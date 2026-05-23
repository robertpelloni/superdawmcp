package integration

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

func TestIntegration_Ableton_E2E(t *testing.T) {
	binPath := "./superdaw-mcp-ableton-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw/main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build server: %v", err)
	}
	defer os.Remove(binPath)

	// Ableton OSC port
	mockAbleton := NewMockDAW(11000)
	go mockAbleton.Start()
	defer mockAbleton.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	reader := json.NewDecoder(stdout)
	writer := json.NewEncoder(stdin)

	// Send Transport Play command
	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"superdaw_transport_control","arguments":{"playing":true,"bpm":130.0,"daw":"ableton"}}`),
		ID:      1,
	}
	writer.Encode(req)

	var res mcp.JSONRPCResponse
	if err := reader.Decode(&res); err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)
	msgs := mockAbleton.GetMessages()

	found := false
	for _, m := range msgs {
		if m.Address == "/superdaw/transport/play" {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("Ableton did not receive expected transport play message")
	}

	stdin.Close()
	cmd.Wait()
}
