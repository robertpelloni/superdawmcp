package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

func TestIntegration_StateSync(t *testing.T) {
	binPath := "./superdaw-sync-test"
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

	time.Sleep(200 * time.Millisecond)

	// Send multiple rapid sync states to test race conditions
	for i := 0; i < 100; i++ {
		_ = mockDAW.SendMessage("127.0.0.1", 11001, "/superdaw/state/tempo", float32(120.0+float32(i)))
	}

	time.Sleep(200 * time.Millisecond)

	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params: json.RawMessage(`{"name": "superdaw_get_transport_state", "arguments": {"daw": "ableton"}}`),
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
}
