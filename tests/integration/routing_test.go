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

func TestIntegration_Routing(t *testing.T) {
	binPath := "./superdaw-routing-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(binPath)

	cmd := exec.Command(binPath)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start cmd: %v", err)
	}
	defer cmd.Process.Kill()

	time.Sleep(200 * time.Millisecond)

	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: json.RawMessage(`{"name": "superdaw_dsl_track_create", "arguments": {"daw": "ableton", "name": "bass"}}`),
		ID:      "1",
	}
	reqBytes, _ := json.Marshal(req)
	fmt.Fprintf(stdin, "%s\n", string(reqBytes))

	var res mcp.JSONRPCResponse
	dec := json.NewDecoder(stdout)
	for {
		if err := dec.Decode(&res); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if res.ID == "1" {
			break
		}
	}

	if res.Error != nil {
		t.Fatalf("Error from server: %+v", *res.Error)
	}

	resultMap, ok := res.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Unexpected result format")
	}

	content, _ := resultMap["content"].([]interface{})
	textMap, _ := content[0].(map[string]interface{})
	text, _ := textMap["text"].(string)

	if !strings.Contains(text, "Sent to Ableton Live") {
		t.Errorf("Routing failed. Got: %s", text)
	}
}
