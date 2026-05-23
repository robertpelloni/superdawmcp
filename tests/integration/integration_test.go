package integration

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/robertpelloni/superdaw-mcp/pkg/client/go"
)

func TestIntegration_EndToEnd(t *testing.T) {
	binPath := "./superdaw-mcp-test"
	// Build the main binary
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw/main.go")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(binPath)

	mockDAW := NewMockDAW(11000)
	go mockDAW.Start()
	defer mockDAW.Stop()

	client, err := client.NewClient(binPath)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	if client == nil {
		t.Fatal("Client is nil")
	}
	defer client.Close()

	err = client.SetMixer("1", 0.7, -0.5)
	if err != nil { t.Errorf("SetMixer failed: %v", err) }

	time.Sleep(200 * time.Millisecond)
	msgs := mockDAW.GetMessages()
	found := false
	for _, m := range msgs {
		if m.Address == "/superdaw/track/volume" {
			found = true
			break
		}
	}
	if !found {
		t.Error("No volume message received")
	}

	client.CreateTrack("Synth", "midi")
	time.Sleep(200 * time.Millisecond)
	msgs = mockDAW.GetMessages()
	found = false
	for _, m := range msgs {
		if m.Address == "/superdaw/track/create" {
			found = true
			break
		}
	}
	if !found {
		t.Error("No track creation message received")
	}

	client.TransportControl(true, 130.0)
	time.Sleep(200 * time.Millisecond)
	msgs = mockDAW.GetMessages()
	found = false
	for _, m := range msgs {
		if m.Address == "/superdaw/transport/play" {
			found = true
			break
		}
	}
	if !found {
		t.Error("No transport message received")
	}
}

func TestIntegration_Bitwig(t *testing.T) {
	binPath := "./superdaw-mcp-bitwig-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw/main.go")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(binPath)

	mockBitwig := NewMockTCPDAW(8181)
	go mockBitwig.Start()
	defer mockBitwig.Stop()

	client, err := client.NewClient(binPath)
	if err != nil {
		t.Fatalf("Failed to initialize client: %v", err)
	}
	if client == nil {
		t.Fatal("Client is nil")
	}
	defer client.Close()

	err = client.SetMixer("1", 0.8, 0.0, "bitwig")
	if err != nil { t.Errorf("SetMixer failed: %v", err) }

	time.Sleep(200 * time.Millisecond)
	msgs := mockBitwig.GetMessages()
	if len(msgs) == 0 {
		t.Error("No TCP messages received")
	}
}
