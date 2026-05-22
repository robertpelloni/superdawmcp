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
func TestIntegration_EndToEnd(t *testing.T) {
	binPath := "./superdaw-mcp-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw/main.go")
	if err := buildCmd.Run(); err != nil { t.Fatalf("Failed to build: %v", err) }
	defer os.Remove(binPath)
	mockDAW := NewMockDAW(11000)
	go mockDAW.Start()
	defer mockDAW.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binPath)
	stdin, _ := cmd.StdinPipe(); stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	writer := json.NewEncoder(stdin); reader := json.NewDecoder(stdout)
	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0", Method: "tools/call",
		Params: json.RawMessage(`{"name":"superdaw_set_mixer","arguments":{"track_id":"1","volume":0.5}}`),
		ID: 1,
	}
	writer.Encode(req)
	var res mcp.JSONRPCResponse
	reader.Decode(&res)
	time.Sleep(100 * time.Millisecond)
	msgs := mockDAW.GetMessages()
	found := false
	for _, m := range msgs {
		if m.Address == "/superdaw/track/volume" { found = true; break }
	}
	if !found { t.Error("No OSC volume message received") }
	stdin.Close(); cmd.Wait()
}

func TestIntegration_Bitwig(t *testing.T) {
	binPath := "./superdaw-mcp-bitwig-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw/main.go")
	if err := buildCmd.Run(); err != nil { t.Fatalf("Failed to build: %v", err) }
	defer os.Remove(binPath)

	mockBitwig := NewMockTCPDAW(8181)
	go mockBitwig.Start()
	defer mockBitwig.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath)
	stdin, _ := cmd.StdinPipe(); stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	writer := json.NewEncoder(stdin); reader := json.NewDecoder(stdout)

	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0", Method: "tools/call",
		Params: json.RawMessage(`{"name":"superdaw_set_mixer","arguments":{"track_id":"1","volume":0.8,"daw":"bitwig"}}`),
		ID: 1,
	}
	writer.Encode(req)
	var res mcp.JSONRPCResponse
	reader.Decode(&res)

	time.Sleep(100 * time.Millisecond)
	msgs := mockBitwig.GetMessages()
	if len(msgs) == 0 {
		t.Error("No TCP messages received by Mock Bitwig")
	}

	stdin.Close(); cmd.Wait()
}
