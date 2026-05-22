package mcp
import "encoding/json"
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{} `json:"id"`
}
type JSONRPCResponse struct {
	JSONRPC string `json:"jsonrpc"`
	Result  interface{} `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
	ID      interface{} `json:"id"`
}
type RPCError struct { Code int `json:"code"`; Message string `json:"message"` }
type Tool struct { Name string `json:"name"`; Description string `json:"description"`; InputSchema interface{} `json:"inputSchema"` }
type ListToolsResult struct { Tools []Tool `json:"tools"` }
func GenerateManifest() ListToolsResult {
	return ListToolsResult{
		Tools: []Tool{
			{
				Name: "superdaw_set_mixer",
				Description: "Set track volume and pan",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"track_id": map[string]interface{}{"type": "string"},
						"volume": map[string]interface{}{"type": "number"},
						"pan": map[string]interface{}{"type": "number"},
						"daw": map[string]interface{}{"type": "string"},
					},
					"required": []string{"track_id", "volume"},
				},
			},
			{
				Name: "superdaw_write_midi",
				Description: "Inject MIDI notes",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"track_id": map[string]interface{}{"type": "string"},
						"daw": map[string]interface{}{"type": "string"},
						"notes": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"pitch": map[string]interface{}{"type": "integer"},
									"velocity": map[string]interface{}{"type": "integer"},
									"start_beat": map[string]interface{}{"type": "number"},
									"duration": map[string]interface{}{"type": "number"},
								},
							},
						},
					},
					"required": []string{"track_id", "notes"},
				},
			},
			{
				Name: "superdaw_generate_euclidean",
				Description: "Generate Euclidean rhythm",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"track_id": map[string]interface{}{"type": "string"},
						"daw": map[string]interface{}{"type": "string"},
						"hits": map[string]interface{}{"type": "integer"},
						"steps": map[string]interface{}{"type": "integer"},
						"pitch": map[string]interface{}{"type": "integer"},
					},
					"required": []string{"track_id", "hits", "steps", "pitch"},
				},
			},
			{
				Name: "superdaw_create_track",
				Description: "Create a new track",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{"type": "string"},
						"type": map[string]interface{}{"type": "string", "enum": []string{"audio", "midi"}},
						"daw": map[string]interface{}{"type": "string"},
					},
					"required": []string{"name", "type"},
				},
			},
			{
				Name: "superdaw_transport_control",
				Description: "Control DAW transport",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"playing": map[string]interface{}{"type": "boolean"},
						"bpm":     map[string]interface{}{"type": "number"},
						"daw":     map[string]interface{}{"type": "string"},
					},
				},
			},
		},
	}
}
