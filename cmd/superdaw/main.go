package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
	"github.com/robertpelloni/superdaw-mcp/pkg/engine"
	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	drivers := map[string]daw.DAWDriver{
		"ableton": daw.NewAbletonDriver("127.0.0.1", 11000, 11001),
		"reaper":  daw.NewReaperDriver("127.0.0.1", 8000, 8080),
		"ardour":  daw.NewArdourDriver("127.0.0.1", 3819),
		"bitwig":  daw.NewBitwigDriver("127.0.0.1", 8181),
	}
	activeDriver := drivers["ableton"]

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return
			}
			continue
		}

		var req mcp.JSONRPCRequest
		if err := json.Unmarshal([]byte(line), &req); err != nil {
			continue
		}

		if req.Method == "initialize" {
			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"serverInfo": map[string]interface{}{
						"name":    "SuperDAW-MCP",
						"version": "1.5.0",
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

			driver := activeDriver
			if d, ok := params.Arguments["daw"].(string); ok {
				if drv, found := drivers[d]; found {
					driver = drv
				}
			}

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
				if !ok {
					b = 120.0
				}
				driver.SetTransportState(p, b)
			}

			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"content": []interface{}{
						map[string]string{
							"type": "text",
							"text": "Success",
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
	fmt.Println(string(out))
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
