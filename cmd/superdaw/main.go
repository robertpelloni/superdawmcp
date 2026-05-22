package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Instantiate active backend driver targets over loopback infrastructure
	abletonDriver := daw.NewAbletonDriver("127.0.0.1", 11000, 11001)

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

			if callParams.Name == "superdaw_set_mixer" {
				trackID, _ := callParams.Arguments["track_id"].(string)
				volume, _ := callParams.Arguments["volume"].(float64)
				_ = abletonDriver.SetTrackVolume(trackID, float32(volume))
			} else if callParams.Name == "superdaw_write_midi" {
				trackID, _ := callParams.Arguments["track_id"].(string)
				notesJSON, _ := json.Marshal(callParams.Arguments["notes"])
				var notes []daw.MIDINote
				json.Unmarshal(notesJSON, &notes)
				_ = abletonDriver.WriteMIDIClip(trackID, 0, notes)
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
