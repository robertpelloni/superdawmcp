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
	// 1. Build the server binary
	binPath := "./superdaw-mcp-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw/main.go")
	if err := buildCmd.Run(); err != nil {
		t.Fatalf("Failed to build server: %v", err)
	}
	defer os.Remove(binPath)

	// 2. Start Mock DAW (Ableton/Unified default)
	mockDAW := NewMockDAW(11000)
	go mockDAW.Start()
	defer mockDAW.Stop()

	// 3. Start MCP Server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, binPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	reader := json.NewDecoder(stdout)
	writer := json.NewEncoder(stdin)

	// 4. Send Tools/Call Request (Set Mixer)
	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"superdaw_set_mixer","arguments":{"track_id":"1","volume":0.5,"pan":0.0}}`),
		ID:      1,
	}
	if err := writer.Encode(req); err != nil {
		t.Fatal(err)
	}

	// 5. Verify Response
	var res mcp.JSONRPCResponse
	if err := reader.Decode(&res); err != nil {
		t.Fatal(err)
	}
	if res.Error != nil {
		t.Fatalf("RPC Error: %v", res.Error.Message)
	}

	// 6. Verify Mock DAW received OSC
	time.Sleep(200 * time.Millisecond) // Wait for OSC delivery
	msgs := mockDAW.GetMessages()

	found := false
	for _, m := range msgs {
		if m.Address == "/superdaw/track/volume" {
			if len(m.Arguments) >= 2 {
				trackID, _ := m.Arguments[0].(string)
				volume, _ := m.Arguments[1].(float32)
				if trackID == "1" && volume == 0.5 {
					found = true
					break
				}
			}
		}
	}

	if !found {
		t.Errorf("Mock DAW did not receive expected OSC message for track volume. Received: %v", msgs)
	}

	// 7. Test MIDI Write
	reqMIDI := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"superdaw_write_midi","arguments":{"track_id":"1","notes":[{"pitch":60,"velocity":100,"start":0,"duration":1}]}}`),
		ID:      2,
	}
	if err := writer.Encode(reqMIDI); err != nil {
		t.Fatal(err)
	}

	if err := reader.Decode(&res); err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)
	msgs = mockDAW.GetMessages()
	foundMIDI := false
	for _, m := range msgs {
		if m.Address == "/superdaw/clip/write" {
			foundMIDI = true
			break
		}
	}
	if !foundMIDI {
		t.Errorf("Mock DAW did not receive expected OSC message for MIDI write.")
	}

	// 8. Test Euclidean Generation
	reqEuclid := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name":"superdaw_generate_euclidean","arguments":{"track_id":"1","hits":3,"steps":8,"pitch":60}}`),
		ID:      3,
	}
	if err := writer.Encode(reqEuclid); err != nil {
		t.Fatal(err)
	}
	if err := reader.Decode(&res); err != nil {
		t.Fatal(err)
	}

	time.Sleep(200 * time.Millisecond)
	msgs = mockDAW.GetMessages()
	foundEuclid := false
	for _, m := range msgs {
		// Euclidean generation also results in a clip write OSC
		if m.Address == "/superdaw/clip/write" {
			foundEuclid = true
			break
		}
	}
	if !foundEuclid {
		t.Errorf("Mock DAW did not receive expected OSC message for Euclidean write.")
	}

	// 9. Cleanup
	stdin.Close()
	cmd.Wait()
}
