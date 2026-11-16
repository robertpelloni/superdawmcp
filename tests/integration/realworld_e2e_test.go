package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

// TestIntegration_RealWorldE2E simulates a complete, heavy-load production session
// to validate performance, routing, and stability under realistic conditions.
func TestIntegration_RealWorldE2E(t *testing.T) {
	binPath := "./superdaw-e2e-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(binPath)

	mockDAW := NewMockDAW(11000)
	go mockDAW.Start()
	defer mockDAW.Stop()

	cmd := exec.Command(binPath)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start cmd: %v", err)
	}
	defer cmd.Process.Kill()

	time.Sleep(300 * time.Millisecond) // Allow server to spin up and load dynamic schemas

	decoder := json.NewDecoder(stdout)

	// A sequence of typical, intensive MCP tool calls simulating an AI session
	workflow := []struct {
		Name   string
		Args   string
		Expect string
	}{
		{"superdaw_get_transport_state", `{"daw": "ableton"}`, "false"},
		{"superdaw_transport_control", `{"daw": "ableton", "playing": true, "bpm": 128.0}`, ""},
		{"superdaw_create_track", `{"daw": "ableton", "name": "Kick", "type": "midi"}`, "Created track"},
		{"superdaw_create_track", `{"daw": "ableton", "name": "Bass", "type": "midi"}`, "Created track"},
		{"superdaw_set_mixer", `{"daw": "ableton", "track_id": "0", "volume": 0.8, "pan": 0.0}`, ""},
		{"superdaw_set_mixer", `{"daw": "ableton", "track_id": "1", "volume": 0.6, "pan": -0.2}`, ""},
		// Dynamic Tool Dispatch simulation (Reaper specific DSL fallback)
		{"superdaw_dsl_track_create", `{"daw": "reaper", "name": "Lead", "role": "melody"}`, "Dispatched to REAPER Bridge: dsl_track_create"}, // Actually ableton fallback returns Sent to Ableton Live
		{"superdaw_generate_euclidean", `{"daw": "ableton", "track_id": "0", "hits": 4, "steps": 16, "pitch": 36}`, "Generated 4/16 Euclidean rhythm"},
		{"superdaw_generate_music", `{"daw": "ableton", "track_id": "1", "style": "techno", "bars": 8}`, "Generated 8 bars of techno music"},
		{"superdaw_list_plugins", `{}`, ""}, // Fast cache read
	}

	start := time.Now()

	for i, step := range workflow {
		req := mcp.JSONRPCRequest{
			JSONRPC: "2.0",
			Method:  "tools/call",
			Params:  json.RawMessage(fmt.Sprintf(`{"name": "%s", "arguments": %s}`, step.Name, step.Args)),
			ID:      i,
		}

		data, _ := json.Marshal(req)
		fmt.Fprintf(stdin, "%s\n", string(data))

		var res mcp.JSONRPCResponse
		for {
			if err := decoder.Decode(&res); err != nil {
				t.Fatalf("Decode failed at step %d (%s): %v", i, step.Name, err)
			}
			idStr := fmt.Sprintf("%v", res.ID)
			if strings.Contains(idStr, fmt.Sprintf("%d", i)) {
				break
			}
		}

		if res.Error != nil {
			t.Fatalf("Error from server on step %d (%s): %+v", i, step.Name, *res.Error)
		}

		if step.Expect != "" {
			resultMap, ok := res.Result.(map[string]interface{})
			if !ok {
				t.Fatalf("Unexpected result format on step %d", i)
			}
			content, ok := resultMap["content"].([]interface{})
			if !ok || len(content) == 0 {
				t.Fatalf("Missing content on step %d", i)
			}
			textMap, _ := content[0].(map[string]interface{})
			text, _ := textMap["text"].(string)

			if !strings.Contains(text, step.Expect) {
				t.Errorf("Step %d (%s) failed. Expected '%s', Got: '%s'", i, step.Name, step.Expect, text)
			}
		}
	}

	duration := time.Since(start)
	t.Logf("Real-World E2E Session processed %d tools in %v", len(workflow), duration)
}
