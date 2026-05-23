package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"

	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
	"github.com/robertpelloni/superdaw-mcp/pkg/engine"
	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
	"github.com/robertpelloni/superdaw-mcp/pkg/vst"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	scanner := vst.NewScanner("vst_cache.json")

	dashboard := StartDashboard(8080)
	router := daw.NewAudioRouter("127.0.0.1", 12000)
	genImporter := engine.NewGenerativeImporter()

	// Cross-platform VST scanning paths
	vstDirs := []string{}
	if runtime.GOOS == "darwin" {
		vstDirs = append(vstDirs, "/Library/Audio/Plug-Ins/VST3")
	} else if runtime.GOOS == "windows" {
		vstDirs = append(vstDirs, `C:\Program Files\Common Files\VST3`)
	} else {
		vstDirs = append(vstDirs, "/usr/lib/vst3", "/usr/local/lib/vst3")
	}
	scanner.ScanDirectories(vstDirs)

	drivers := map[string]daw.DAWDriver{
		"ableton":  daw.NewAbletonDriver("127.0.0.1", 11000, 11001),
		"reaper":   daw.NewReaperDriver("127.0.0.1", 8000, 8080),
		"ardour":   daw.NewArdourDriver("127.0.0.1", 3819),
		"bitwig":   daw.NewBitwigDriver("127.0.0.1", 8181),
		"flstudio": daw.NewFLStudioDriver("127.0.0.1", 9000),
	}
	activeDriver := drivers["ableton"]

	// Initialize drivers that require permanent connections
	if drv, ok := drivers["bitwig"]; ok {
		go func() {
			if err := drv.Connect("127.0.0.1:8181"); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to connect to Bitwig: %v\n", err)
			}
		}()
	}

	// Read version from VERSION.md
	version := "1.4.0"
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
					"capabilities": map[string]interface{}{
						"tools": map[string]interface{}{"listChanged": true},
					},
					"serverInfo": map[string]interface{}{
						"name":    "SuperDAW-MCP",
						"version": version,
					},
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

			dawName := "ableton"
			if d, ok := params.Arguments["daw"].(string); ok { dawName = d }
			driver := drivers[dawName]
			if driver == nil { driver = activeDriver }

			var result interface{}
			result = "Success"

			switch params.Name {
			case "superdaw_set_mixer":
				id, _ := params.Arguments["track_id"].(string)
				vol, _ := params.Arguments["volume"].(float64)
				driver.SetTrackVolume(id, float32(vol))
				if p, ok := params.Arguments["pan"].(float64); ok {
					driver.SetTrackPan(id, float32(p))
				}
			case "superdaw_write_midi":
				id, _ := params.Arguments["track_id"].(string)
				notesJSON, _ := json.Marshal(params.Arguments["notes"])
				var notes []daw.MIDINote
				json.Unmarshal(notesJSON, &notes)
				driver.WriteMIDIClip(id, 0, notes)
			case "superdaw_generate_euclidean":
				id, _ := params.Arguments["track_id"].(string)
				hits, _ := params.Arguments["hits"].(float64)
				steps, _ := params.Arguments["steps"].(float64)
				pitch, _ := params.Arguments["pitch"].(float64)
				notes := engine.GenerateEuclidean(int(hits), int(steps), int(pitch), 100, 0, 4.0)
				driver.WriteMIDIClip(id, 0, notes)
			case "superdaw_create_track":
				n, _ := params.Arguments["name"].(string)
				t, _ := params.Arguments["type"].(string)
				driver.CreateTrack(n, t)
			case "superdaw_transport_control":
				p, _ := params.Arguments["playing"].(bool)
				b, ok := params.Arguments["bpm"].(float64)
				if !ok { b = 120.0 }
				driver.SetTransportState(p, b)
				dashboard.UpdateDAW(dawName, p, b)
			case "superdaw_list_clips":
				id, _ := params.Arguments["track_id"].(string)
				result, _ = driver.ListClips(id)
			case "superdaw_delete_clip":
				id, _ := params.Arguments["track_id"].(string)
				idx, _ := params.Arguments["clip_idx"].(float64)
				driver.DeleteClip(id, int(idx))
			case "superdaw_list_plugins":
				result = scanner.ListPlugins()
			case "superdaw_get_plugin_params":
				name, _ := params.Arguments["plugin_name"].(string)
				result, _ = scanner.GetPluginMetadata(name)
			case "superdaw_separate_stems":
				in, _ := params.Arguments["input_path"].(string)
				out, _ := params.Arguments["output_dir"].(string)
				stems, ok := params.Arguments["stems"].(float64)
				if !ok { stems = 4 }
				result, _ = engine.SeparateStems(in, out, int(stems))
			case "superdaw_custom_command":
				cmd, _ := params.Arguments["command"].(string)
				args, _ := params.Arguments["args"].(map[string]interface{})
				result, _ = driver.ExecuteCustomCommand(cmd, args)

			// PHASE 4 TOOLS
			case "superdaw_patch_audio":
				srcDaw, _ := params.Arguments["source_daw"].(string)
				srcTrack, _ := params.Arguments["source_track"].(string)
				dstDaw, _ := params.Arguments["dest_daw"].(string)
				dstTrack, _ := params.Arguments["dest_track"].(string)
				router.Patch(srcDaw, srcTrack, dstDaw, dstTrack)
				dashboard.AddPatch(AudioPatch{SourceDAW: srcDaw, SourceTrack: srcTrack, DestDAW: dstDaw, DestTrack: dstTrack})

			case "superdaw_import_generative":
				prompt, _ := params.Arguments["prompt"].(string)
				target, _ := params.Arguments["target_daw"].(string)
				result, _ = genImporter.ImportStems(prompt, target)
			}

			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"content": []interface{}{
						map[string]interface{}{
							"type": "text",
							"text": fmt.Sprintf("%v", result),
						},
					},
				},
			}
			writeResponse(res)
		}
	}
}

func writeResponse(res mcp.JSONRPCResponse) {
	out, _ := json.Marshal(res)
	os.Stdout.Write(append(out, '\n'))
}

func sendError(id interface{}, code int, message string) {
	res := mcp.JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &mcp.RPCError{
			Code:    code,
			Message: message,
		},
	}
	writeResponse(res)
}
