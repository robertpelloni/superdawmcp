package integration

import (
	"testing"
	"time"

	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
	"github.com/hypebeast/go-osc/osc"
)

func TestRealtimeSync_Ableton(t *testing.T) {
	// 1. Start Ableton Driver with a specific local port
	localPort := 11005
	driver := daw.NewAbletonDriver("127.0.0.1", 11000, localPort)

	// 2. Simulate OSC message from Ableton Agent
	client := osc.NewClient("127.0.0.1", localPort)
	msg := osc.NewMessage("/superdaw/state/playing")
	msg.Append(true)

	err := client.Send(msg)
	if err != nil { t.Fatalf("Failed to send OSC: %v", err) }

	// 3. Wait for async processing
	time.Sleep(100 * time.Millisecond)

	// 4. Verify Driver Cache
	playing, _, _ := driver.GetTransportState()
	if !playing {
		t.Error("Driver did not reflect real-time playing state change")
	}

	// 5. Test Tempo
	msg2 := osc.NewMessage("/superdaw/state/tempo")
	msg2.Append(float32(135.0))
	client.Send(msg2)

	time.Sleep(100 * time.Millisecond)
	_, tempo, _ := driver.GetTransportState()
	if tempo != 135.0 {
		t.Errorf("Driver did not reflect real-time tempo change. Expected 135, got %f", tempo)
	}
}
