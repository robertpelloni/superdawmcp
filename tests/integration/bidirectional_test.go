package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
	"strings"

	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

func TestIntegration_BidirectionalSync(t *testing.T) {
	binPath := "./superdaw-bidirectional-test"
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
	cmd.Start()
	defer cmd.Process.Kill()

	err := mockDAW.SendMessage("127.0.0.1", 11001, "/superdaw/state/tempo", float32(145.0))
	if err != nil {
		t.Fatalf("Failed to send mock state update: %v", err)
	}
	err = mockDAW.SendMessage("127.0.0.1", 11001, "/superdaw/state/playing", true)
	if err != nil {
		t.Fatalf("Failed to send mock state update: %v", err)
	}

	time.Sleep(500 * time.Millisecond)

	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: json.RawMessage(`{"name": "superdaw_get_transport_state", "arguments": {"daw": "ableton"}}`),
		ID:      1,
	}
	reqBytes, _ := json.Marshal(req)
	fmt.Fprintf(stdin, "%s\n", string(reqBytes))

	var res mcp.JSONRPCResponse
	dec := json.NewDecoder(stdout)
	if err := dec.Decode(&res); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	resultMap, ok := res.Result.(map[string]interface{})
	if !ok {
		t.Fatalf("Unexpected result format: %v", res.Result)
	}

	content, _ := resultMap["content"].([]interface{})
	textMap, _ := content[0].(map[string]interface{})
	text, _ := textMap["text"].(string)

	if !strings.Contains(text, "145") || !strings.Contains(text, "true") {
		t.Errorf("State sync failed. Got: %s", text)
	}
}
