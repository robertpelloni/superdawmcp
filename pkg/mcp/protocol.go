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
						"daw":      map[string]interface{}{"type": "string", "description": "Target DAW: 'ableton', 'reaper', or 'ardour'"},
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
						"daw":      map[string]interface{}{"type": "string", "description": "Target DAW: 'ableton', 'reaper', or 'ardour'"},
						"notes": map[string]interface{}{
							"type": "array",
							"items": map[string]interface{}{
								"type": "object",
								"properties": map[string]interface{}{
									"pitch":      map[string]interface{}{"type": "integer"},
									"velocity":   map[string]interface{}{"type": "integer"},
									"start_beat": map[string]interface{}{"type": "number"},
									"duration":   map[string]interface{}{"type": "number"},
								},
							},
						},
					},
					"required": []string{"track_id", "notes"},
				},
			},
			{
				Name:        "superdaw_generate_euclidean",
				Description: "Generate a Euclidean rhythm pattern and inject it into a track.",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"track_id": map[string]interface{}{"type": "string"},
						"daw":      map[string]interface{}{"type": "string", "description": "Target DAW: 'ableton', 'reaper', or 'ardour'"},
						"hits":     map[string]interface{}{"type": "integer", "description": "Number of hits"},
						"steps":    map[string]interface{}{"type": "integer", "description": "Total steps (e.g., 16 for 1 bar at 1/16th)"},
						"pitch":    map[string]interface{}{"type": "integer", "description": "MIDI note number"},
						"velocity": map[string]interface{}{"type": "integer", "default": 100},
						"rotate":   map[string]interface{}{"type": "integer", "default": 0},
						"length":   map[string]interface{}{"type": "number", "default": 4.0, "description": "Clip length in beats"},
					},
					"required": []string{"track_id", "hits", "steps", "pitch"},
				},
			},
			{
				Name:        "superdaw_separate_stems",
				Description: "Separate an audio file into stems using AI (Spleeter).",
				InputSchema: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"input_path": map[string]interface{}{"type": "string", "description": "Path to the source audio file"},
						"output_dir": map[string]interface{}{"type": "string", "description": "Directory to save the separated stems"},
						"stems":      map[string]interface{}{"type": "integer", "enum": []int{2, 4, 5}, "default": 4, "description": "Number of stems (2, 4, or 5)"},
					},
					"required": []string{"input_path", "output_dir"},
				},
			},
		},
	}
}
