package mcp

import (
	"encoding/json"
)

type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id"`
}

type JSONRPCResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError   `json:"error,omitempty"`
	ID      interface{}     `json:"id"`
}

type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type ListToolsResult struct {
	Tools []Tool `json:"tools"`
}

// GenerateManifest produces the schema parameters expected by the LLM client engine.
func GenerateManifest() ListToolsResult {
	return ListToolsResult{
		Tools: []Tool{
			{
				Name:        "superdaw_set_mixer",
				Description: "Modify track configurations (volume, panning, mute) inside the targeted active DAW pipeline.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"track_id": map[string]interface{}{"type": "string"},
						"volume":   map[string]interface{}{"type": "number", "description": "Volume float range 0.0 to 1.0"},
						"pan":      map[string]interface{}{"type": "number", "description": "Panning float range -1.0 to 1.0"},
					},
					"required": []string{"track_id", "volume"},
				},
			},
			{
				Name:        "superdaw_write_midi",
				Description: "Inject sequences of musical pitches directly onto an active structural instrument lane track.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"track_id": map[string]interface{}{"type": "string"},
						"notes": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"pitch":    map[string]interface{}{"type": "integer"},
									"velocity": map[string]interface{}{"type": "integer"},
									"start":    map[string]interface{}{"type": "number"},
									"duration": map[string]interface{}{"type": "number"},
								},
							},
						},
					},
					"required": []string{"track_id", "notes"},
				},
			},
		},
	}
}
