package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
	"github.com/robertpelloni/superdaw-mcp/pkg/engine"
	"github.com/hypebeast/go-osc/osc"
	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
	"github.com/robertpelloni/superdaw-mcp/pkg/ux/dashboard"
	"sync"
	"github.com/robertpelloni/superdaw-mcp/pkg/vst"
)

var stdoutMu sync.Mutex

func main() {
	reader := bufio.NewReader(os.Stdin)
	scanner := vst.NewScanner("vst_cache.json")

	dash, dashMux := dashboard.StartDashboard(8081)
	dashboard.RegisterMobileRemote(dashMux)

	router := daw.NewAudioRouter("127.0.0.1", 12000)
	genImporter := engine.NewGenerativeImporter()
	link := engine.NewLinkBridge()

	manager := daw.NewConnectionManager()

	// Register default instances
	manager.Register("ableton", daw.NewAbletonDriver("127.0.0.1", 11000, 11001))
	manager.Register("reaper", daw.NewReaperDriver("127.0.0.1", 8000, 8080))
	manager.Register("ardour", daw.NewArdourDriver("127.0.0.1", 3819))
	manager.Register("bitwig", daw.NewBitwigDriver("127.0.0.1", 8181))
	manager.Register("flstudio", daw.NewFLStudioDriver("127.0.0.1", 9000))
	manager.Register("logic", daw.NewLogicProDriver("127.0.0.1", 12100, 12101))
	manager.Register("cubase", daw.NewCubaseDriver("127.0.0.1", 7001))
	manager.Register("protools", daw.NewProToolsDriver("127.0.0.1", 7002))

	manager.SetDefault("ableton")

	for id, drv := range manager.GetAll() {
		instanceID := id
		drv.SetNotifyHandler(func(method string, params interface{}) {
			notif := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  method,
				"params":  params,
			}
			writeResponseRaw(notif)

			// Collaborative updates default to "default" room for now
			// In production, drivers would be assigned to specific rooms
			roomID := "default"

			if method == "superdaw/arrangement_update" {
				if p, ok := params.(map[string]interface{}); ok {
					if arr, ok := p["arrangement"].(string); ok {
						dash.UpdateArrangement(roomID, instanceID, arr)
					}
				}
			}

			jobs := genImporter.GetJobs()
			jList := make([]interface{}, len(jobs))
			for i, j := range jobs { jList[i] = j }
			dash.UpdateJobs(roomID, jList)
		})
	}

	go func() {
		disp := osc.NewStandardDispatcher()
		handler := func(m *osc.Message) {
			parts := strings.Split(m.Address, "/")
			if len(parts) < 4 { return }
			instanceID := parts[2]
			cmd := parts[3]
			driver, err := manager.Get(instanceID)
			if err != nil { return }
			switch cmd {
			case "transport":
				if len(m.Arguments) > 0 {
					p, _ := m.Arguments[0].(int32)
					driver.SetTransportState(p == 1, 120.0)
				}
			case "volume":
				if len(m.Arguments) > 1 {
					id, _ := m.Arguments[0].(string)
					vol, _ := m.Arguments[1].(float32)
					driver.SetTrackVolume(id, vol)
				}
			case "pan":
				if len(m.Arguments) > 1 {
					id, _ := m.Arguments[0].(string)
					p, _ := m.Arguments[1].(float32)
					driver.SetTrackPan(id, p)
				}
			case "track_create":
				if len(m.Arguments) > 1 {
					name, _ := m.Arguments[0].(string)
					t, _ := m.Arguments[1].(string)
					driver.CreateTrack(name, t)
				}
			}
		}
		disp.AddMsgHandler("/superdaw/*/transport", handler)
		disp.AddMsgHandler("/superdaw/*/volume", handler)
		disp.AddMsgHandler("/superdaw/*/pan", handler)
		disp.AddMsgHandler("/superdaw/*/track_create", handler)
		disp.AddMsgHandler("/superdaw/transport/play", func(m *osc.Message) {
			if len(m.Arguments) > 0 {
				p, _ := m.Arguments[0].(int32)
				if drv, err := manager.Get(""); err == nil {
					drv.SetTransportState(p == 1, 120.0)
				}
			}
		})
		server := &osc.Server{Addr: "0.0.0.0:12001", Dispatcher: disp}
		server.ListenAndServe()
	}()

	vstDirs := []string{}
	if runtime.GOOS == "darwin" {
		vstDirs = append(vstDirs, "/Library/Audio/Plug-Ins/VST3")
	} else if runtime.GOOS == "windows" {
		vstDirs = append(vstDirs, `C:\Program Files\Common Files\VST3`)
	} else {
		vstDirs = append(vstDirs, "/usr/lib/vst3", "/usr/local/lib/vst3")
	}
	scanner.ScanDirectories(vstDirs)

	// TCP MCP Gateway (Remote SDK Access)
	go func() {
		ln, err := net.Listen("tcp", "127.0.0.1:12002")
		if err != nil {
			fmt.Fprintf(os.Stderr, "TCP Gateway failed: %v\n", err)
			return
		}
		for {
			conn, err := ln.Accept()
			if err != nil { continue }
			go func(c net.Conn) {
				defer c.Close()
				tcpScanner := bufio.NewScanner(c)
				for tcpScanner.Scan() {
					line := tcpScanner.Text()
					var req mcp.JSONRPCRequest
					if err := json.Unmarshal([]byte(line), &req); err != nil { continue }

					if req.Method == "tools/call" {
						var params struct {
							Name      string                 `json:"name"`
							Arguments map[string]interface{} `json:"arguments"`
						}
						json.Unmarshal(req.Params, &params)
						res := handleToolCall(params.Name, params.Arguments, manager, scanner, dash, genImporter, link, router)

						resp := mcp.JSONRPCResponse{
							JSONRPC: "2.0",
							ID:      req.ID,
							Result:  map[string]interface{}{"content": []interface{}{map[string]interface{}{"type": "text", "text": fmt.Sprintf("%v", res)}}},
						}
						data, _ := json.Marshal(resp)
						c.Write(append(data, '\n'))
					}
				}
			}(conn)
		}
	}()

	go func() {
		for cmd := range dashboard.CommandBus {
			handleToolCallWithRoom(cmd.RoomID, cmd.Name, cmd.Arguments, manager, scanner, dash, genImporter, link, router)
		}
	}()

	if drv, err := manager.Get("bitwig"); err == nil {
		go func() {
			if err := drv.Connect("127.0.0.1:8181"); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to connect to Bitwig: %v\n", err)
			}
		}()
	}

	version := "3.2.0"
	versionData, err := os.ReadFile("VERSION.md")
	if err == nil {
		version = strings.TrimSpace(string(versionData))
	}

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF { return }
			continue
		}
		var req mcp.JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil { continue }
		if req.Method == "initialize" {
			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities": map[string]interface{}{"tools": map[string]interface{}{"listChanged": true}},
					"serverInfo": map[string]interface{}{"name": "SuperDAW-MCP", "version": version},
				},
			}
			writeResponse(res)
		} else if req.Method == "tools/list" {
			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  mcp.GenerateManifest(),
			}
			writeResponse(res)
		} else if req.Method == "tools/call" {
			var params struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			}
			if err := json.Unmarshal(req.Params, &params); err != nil {
				sendError(req.ID, -32602, "Invalid params")
				continue
			}
			// MCP standard tools default to "default" room
			result := handleToolCallWithRoom("default", params.Name, params.Arguments, manager, scanner, dash, genImporter, link, router)
			var content []interface{}
			if text, ok := result.(string); ok {
				content = append(content, map[string]interface{}{"type": "text", "text": text})
			} else {
				jsonRes, _ := json.Marshal(result)
				content = append(content, map[string]interface{}{"type": "text", "text": string(jsonRes)})
			}
			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  map[string]interface{}{"content": content},
			}
			writeResponse(res)
		}
	}
}

func writeResponse(res mcp.JSONRPCResponse) { writeResponseRaw(res) }
func writeResponseRaw(res interface{}) {
	stdoutMu.Lock()
	defer stdoutMu.Unlock()
	out, _ := json.Marshal(res)
	os.Stdout.Write(append(out, '\n'))
}
func sendError(id interface{}, code int, message string) {
	res := mcp.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error:   &mcp.RPCError{Code: code, Message: message},
	}
	writeResponse(res)
}

func handleToolCall(name string, args map[string]interface{}, manager *daw.ConnectionManager, scanner *vst.Scanner, dash *dashboard.DashboardState, genImporter *engine.GenerativeImporter, link *engine.LinkBridge, router *daw.AudioRouter) interface{} {
	return handleToolCallWithRoom("default", name, args, manager, scanner, dash, genImporter, link, router)
}

func handleToolCallWithRoom(roomID string, name string, args map[string]interface{}, manager *daw.ConnectionManager, scanner *vst.Scanner, dash *dashboard.DashboardState, genImporter *engine.GenerativeImporter, link *engine.LinkBridge, router *daw.AudioRouter) interface{} {
	instanceID := ""
	if id, ok := args["daw"].(string); ok { instanceID = id }
	driver, err := manager.Get(instanceID)
	if err != nil { return err.Error() }

	var result interface{}
	result = "Success"
	switch name {
	case "superdaw_set_mixer":
		id, _ := args["track_id"].(string); vol, _ := args["volume"].(float64)
		driver.SetTrackVolume(id, float32(vol))
		if p, ok := args["pan"].(float64); ok { driver.SetTrackPan(id, float32(p)) }
	case "superdaw_set_instrument":
		id, _ := args["track_id"].(string); inst, _ := args["instrument"].(string)
		driver.SetTrackInstrument(id, inst)
		result = fmt.Sprintf("Loaded %s on track %s", inst, id)
	case "superdaw_write_midi":
		id, _ := args["track_id"].(string); clipIdx := 0
		if idx, ok := args["clip_index"].(float64); ok { clipIdx = int(idx) }
		notesJSON, _ := json.Marshal(args["notes"])
		var notes []daw.MIDINote
		json.Unmarshal(notesJSON, &notes)
		driver.WriteMIDIClip(id, clipIdx, notes)
	case "superdaw_create_track":
		n, _ := args["name"].(string); t, _ := args["type"].(string)
		driver.CreateTrack(n, t)
		result = fmt.Sprintf("Created track: %s (%s)", n, t)
	case "superdaw_generate_euclidean":
		id, _ := args["track_id"].(string); hits, _ := args["hits"].(float64); steps, _ := args["steps"].(float64); pitch, _ := args["pitch"].(float64)
		if pitch == 0 { pitch = 60 }
		velocity := 100; if v, ok := args["velocity"].(float64); ok { velocity = int(v) }
		rotation := 0; if r, ok := args["rotation"].(float64); ok { rotation = int(r) }
		length := float32(4.0); if l, ok := args["length"].(float64); ok { length = float32(l) }
		notes := engine.GenerateEuclidean(int(hits), int(steps), int(pitch), velocity, rotation, length)
		driver.WriteMIDIClip(id, 0, notes)
		result = fmt.Sprintf("Generated %d/%d Euclidean rhythm on track %s", int(hits), int(steps), id)
	case "superdaw_transport_control":
		p, _ := args["playing"].(bool); b, ok := args["bpm"].(float64)
		if !ok { b = 120.0 }
		driver.SetTransportState(p, b); dash.UpdateDAW(roomID, instanceID, p, b); link.Sync(p, b)
		manager.CacheTransportState(manager.ResolveID(instanceID), p, b)
	case "superdaw_get_tracks":
		tracks, _ := driver.GetTracks(); result = tracks
	case "superdaw_get_transport_state":
		if s, ok := manager.GetCachedTransportState(manager.ResolveID(instanceID)); ok {
			result = map[string]interface{}{"playing": s.Playing, "bpm": s.BPM}
		} else {
			p, b, _ := driver.GetTransportState()
			result = map[string]interface{}{"playing": p, "bpm": b}
		}
	case "superdaw_list_clips":
		id, _ := args["track_id"].(string); result, _ = driver.ListClips(id)
	case "superdaw_delete_clip":
		id, _ := args["track_id"].(string); idx, _ := args["clip_idx"].(float64); driver.DeleteClip(id, int(idx))
	case "superdaw_list_plugins":
		result = scanner.ListPlugins()
	case "superdaw_get_plugin_params":
		name, _ := args["plugin_name"].(string); result, _ = scanner.GetPluginMetadata(name)
	case "superdaw_set_plugin_parameter":
		pName, _ := args["plugin_name"].(string); paramName, _ := args["parameter_name"].(string); val, _ := args["value"].(float64); trackID, _ := args["track_id"].(string)
		meta, ok := scanner.GetPluginMetadata(pName)
		if !ok { result = "Plugin not found." } else {
			idx := -1
			for _, pm := range meta.Parameters {
				if strings.EqualFold(pm.Name, paramName) { idx = pm.Index; break }
			}
			if idx == -1 { result = "Parameter not found." } else {
				driver.SetPluginParameter(trackID, pName, idx, float32(val))
				result = fmt.Sprintf("Set %s:%s to %f", pName, paramName, val)
			}
		}
	case "superdaw_send_cc":
		id, _ := args["track_id"].(string); ctrl, _ := args["controller"].(float64); val, _ := args["value"].(float64)
		driver.SendCC(id, int(ctrl), int(val))
		result = fmt.Sprintf("Sent CC %d:%d to track %s", int(ctrl), int(val), id)
	case "superdaw_fire_scene":
		result, _ = driver.ExecuteCustomCommand("fire_scene", args)
	case "superdaw_separate_stems":
		in, _ := args["input_path"].(string); out, _ := args["output_dir"].(string); stems, ok := args["stems"].(float64); if !ok { stems = 4 }
		if r, err := engine.SeparateStems(in, out, int(stems)); err != nil {
			result = err.Error()
		} else {
			result = r
		}
	case "superdaw_custom_command":
		cmd, _ := args["command"].(string); cargs, _ := args["args"].(map[string]interface{}); result, _ = driver.ExecuteCustomCommand(cmd, cargs)
	case "superdaw_patch_audio":
		srcDaw, _ := args["source_daw"].(string); srcTrack, _ := args["source_track"].(string); dstDaw, _ := args["dest_daw"].(string); dstTrack, _ := args["dest_track"].(string)
		router.Patch(srcDaw, srcTrack, dstDaw, dstTrack); dash.AddPatch(roomID, dashboard.AudioPatch{SourceDAW: srcDaw, SourceTrack: srcTrack, DestDAW: dstDaw, DestTrack: dstTrack})
	case "superdaw_unpatch_audio":
		srcDaw, _ := args["source_daw"].(string); srcTrack, _ := args["source_track"].(string); dstDaw, _ := args["dest_daw"].(string); dstTrack, _ := args["dest_track"].(string)
		router.Unpatch(srcDaw, srcTrack, dstDaw, dstTrack); dash.RemovePatch(roomID, dashboard.AudioPatch{SourceDAW: srcDaw, SourceTrack: srcTrack, DestDAW: dstDaw, DestTrack: dstTrack})
	case "superdaw_import_generative":
		prompt, _ := args["prompt"].(string); target, _ := args["target_daw"].(string); result, _ = genImporter.ImportStems(prompt, target)
	case "superdaw_list_generative_jobs":
		result = genImporter.GetJobs()
	case "superdaw_save_session":
		fname, _ := args["filename"].(string); if fname == "" { fname = "studio_session.json" }
		data, _ := json.Marshal(dash.GetState(roomID)); os.WriteFile(fname, data, 0644); result = "Session saved."
	case "superdaw_load_session":
		fname, _ := args["filename"].(string); if fname == "" { fname = "studio_session.json" }
		data, _ := os.ReadFile(fname); var state dashboard.RoomState; json.Unmarshal(data, &state); dash.SetState(roomID, state); result = "Session loaded."
	case "superdaw_generate_music":
		style, _ := args["style"].(string); bars, _ := args["bars"].(float64); trackID, _ := args["track_id"].(string)
		notes := engine.GenerateMusic(style, int(bars)); driver.WriteMIDIClip(trackID, 0, notes); result = fmt.Sprintf("Generated %d bars of %s music.", int(bars), style)
	}
	return result
}
