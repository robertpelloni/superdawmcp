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

	// Test MIDI Write
	reqMIDI := mcp.JSONRPCRequest{
		JSONRPC: "2.0", Method: "tools/call",
		Params: json.RawMessage(`{"name":"superdaw_write_midi","arguments":{"track_id":"1","notes":[{"pitch":60,"velocity":100,"start_beat":0,"duration":1}]}}`),
		ID: 2,
	}
	writer.Encode(reqMIDI)
	reader.Decode(&res)
	time.Sleep(100 * time.Millisecond)
	msgs = mockDAW.GetMessages()
	foundMIDI := false
	for _, m := range msgs {
		if m.Address == "/superdaw/clip/write" { foundMIDI = true; break }
	}
	if !foundMIDI { t.Error("No OSC MIDI write message received") }

	// Test Euclidean Generation
	reqEuclid := mcp.JSONRPCRequest{
		JSONRPC: "2.0", Method: "tools/call",
		Params: json.RawMessage(`{"name":"superdaw_generate_euclidean","arguments":{"track_id":"1","hits":3,"steps":8,"pitch":60}}`),
		ID: 3,
	}
	writer.Encode(reqEuclid)
	reader.Decode(&res)
	time.Sleep(100 * time.Millisecond)
	msgs = mockDAW.GetMessages()
	foundEuclid := false
	for _, m := range msgs {
		if m.Address == "/superdaw/clip/write" { foundEuclid = true; break }
	}
	if !foundEuclid { t.Error("No OSC Euclidean write message received") }

	// Test Create Track
	reqTrack := mcp.JSONRPCRequest{
		JSONRPC: "2.0", Method: "tools/call",
		Params: json.RawMessage(`{"name":"superdaw_create_track","arguments":{"name":"Synth","type":"midi"}}`),
		ID: 4,
	}
	writer.Encode(reqTrack)
	reader.Decode(&res)
	time.Sleep(100 * time.Millisecond)
	msgs = mockDAW.GetMessages()
	foundTrack := false
	for _, m := range msgs {
		if m.Address == "/superdaw/track/create" { foundTrack = true; break }
	}
	if !foundTrack { t.Error("No OSC create track message received") }

	// Test Transport Control
	reqTransport := mcp.JSONRPCRequest{
		JSONRPC: "2.0", Method: "tools/call",
		Params: json.RawMessage(`{"name":"superdaw_transport_control","arguments":{"playing":true,"bpm":128.0}}`),
		ID: 5,
	}
	writer.Encode(reqTransport)
	reader.Decode(&res)
	time.Sleep(100 * time.Millisecond)
	msgs = mockDAW.GetMessages()
	foundTransport := false
	for _, m := range msgs {
		if m.Address == "/superdaw/transport/play" { foundTransport = true; break }
	}
	if !foundTransport { t.Error("No OSC transport play message received") }

	stdin.Close(); cmd.Wait()
}

func TestIntegration_Plugins(t *testing.T) {
	binPath := "./superdaw-mcp-plugin-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw/main.go")
	if err := buildCmd.Run(); err != nil { t.Fatalf("Failed to build: %v", err) }
	defer os.Remove(binPath)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath)
	stdin, _ := cmd.StdinPipe(); stdout, _ := cmd.StdoutPipe()
	cmd.Start()
	writer := json.NewEncoder(stdin); reader := json.NewDecoder(stdout)

	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0", Method: "tools/call",
		Params: json.RawMessage(`{"name":"superdaw_list_plugins","arguments":{}}`),
		ID: 1,
	}
	writer.Encode(req)
	var res mcp.JSONRPCResponse
	reader.Decode(&res)

	if res.Error != nil {
		t.Errorf("Unexpected error: %v", res.Error.Message)
	}

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
