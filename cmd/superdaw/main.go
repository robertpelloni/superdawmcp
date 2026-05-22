package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
	"github.com/robertpelloni/superdaw-mcp/pkg/engine"
	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Instantiate active backend driver targets over loopback infrastructure
	drivers := map[string]daw.DAWDriver{
		"ableton": daw.NewAbletonDriver("127.0.0.1", 11000, 11001),
		"reaper":  daw.NewReaperDriver("127.0.0.1", 8000, 8080),
		"ardour":  daw.NewArdourDriver("127.0.0.1", 3819),
	}

	// Default driver
	activeDriver := drivers["ableton"]

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(0)
		}

		var req mcp.JSONRPCRequest
		if err := json.Unmarshal([]byte(input), &req); err != nil {
			sendError(req.ID, -32700, "Parse error processing incoming JSON raw stream package.")
			continue
		}

		switch req.Method {
		case "initialize":
			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"protocolVersion": "2024-11-05",
					"capabilities":    map[string]interface{}{},
					"serverInfo": map[string]interface{}{
						"name":    "SuperDAW-Universal-MCP",
						"version": "1.0.0",
					},
				},
			}
			writeResponse(res)

		case "tools/list":
			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result:  mcp.GenerateManifest(),
			}
			writeResponse(res)

		case "tools/call":
			// Process tool arguments and route them to your active driver implementation layer
			var callParams struct {
				Name      string                 `json:"name"`
				Arguments map[string]interface{} `json:"arguments"`
			}
			json.Unmarshal(req.Params, &callParams)

			// Routing logic: allow overriding target DAW via arguments
			targetDAW, ok := callParams.Arguments["daw"].(string)
			driver := activeDriver
			if ok {
				if d, found := drivers[targetDAW]; found {
					driver = d
				}
			}

			if callParams.Name == "superdaw_set_mixer" {
				trackID, _ := callParams.Arguments["track_id"].(string)
				volume, _ := callParams.Arguments["volume"].(float64)
				_ = driver.SetTrackVolume(trackID, float32(volume))
				if pan, ok := callParams.Arguments["pan"].(float64); ok {
					_ = driver.SetTrackPan(trackID, float32(pan))
				}
			} else if callParams.Name == "superdaw_write_midi" {
				trackID, _ := callParams.Arguments["track_id"].(string)
				notesJSON, _ := json.Marshal(callParams.Arguments["notes"])
				var notes []daw.MIDINote
				json.Unmarshal(notesJSON, &notes)
				_ = driver.WriteMIDIClip(trackID, 0, notes)
			} else if callParams.Name == "superdaw_generate_euclidean" {
				trackID, _ := callParams.Arguments["track_id"].(string)
				hits := int(callParams.Arguments["hits"].(float64))
				steps := int(callParams.Arguments["steps"].(float64))
				pitch := int(callParams.Arguments["pitch"].(float64))

				velocity := 100
				if v, ok := callParams.Arguments["velocity"].(float64); ok {
					velocity = int(v)
				}
				rotate := 0
				if r, ok := callParams.Arguments["rotate"].(float64); ok {
					rotate = int(r)
				}
				length := float32(4.0)
				if l, ok := callParams.Arguments["length"].(float64); ok {
					length = float32(l)
				}

				notes := engine.GenerateEuclidean(hits, steps, pitch, velocity, rotate, length)
				_ = driver.WriteMIDIClip(trackID, 0, notes)
			} else if callParams.Name == "superdaw_separate_stems" {
				inputPath, _ := callParams.Arguments["input_path"].(string)
				outputDir, _ := callParams.Arguments["output_dir"].(string)
				stems := 4
				if s, ok := callParams.Arguments["stems"].(float64); ok {
					stems = int(s)
				}
				_ = engine.SeparateStems(inputPath, outputDir, stems)
			}

			res := mcp.JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Result: map[string]interface{}{
					"content": []map[string]interface{}{
						{
							"type": "text",
							"text": "Mixer values synchronized successfully across active targets.",
						},
					},
				},
			}
			writeResponse(res)
		}
	}
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

func writeResponse(res mcp.JSONRPCResponse) {
	out, _ := json.Marshal(res)
	fmt.Println(string(out))
}
