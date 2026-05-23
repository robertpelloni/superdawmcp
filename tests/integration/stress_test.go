package integration

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"
	"time"

	"github.com/robertpelloni/superdaw-mcp/pkg/client/go"
)

func TestStress_HighFrequency(t *testing.T) {
	binPath := "./superdaw-stress-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(binPath)

	mockDAW := NewMockDAW(11000)
	go mockDAW.Start()
	defer mockDAW.Stop()

	client, err := client.NewClient(binPath)
	if err != nil { t.Fatalf("Failed to initialize: %v", err) }
	defer client.Close()

	numMessages := 500
	var wg sync.WaitGroup
	wg.Add(numMessages)

	start := time.Now()
	for i := 0; i < numMessages; i++ {
		go func(idx int) {
			defer wg.Done()
			vol := float32(idx) / float32(numMessages)
			client.SetMixer("1", vol, 0.0)
		}(i)
	}
	wg.Wait()
	duration := time.Since(start)

	t.Logf("Sent %d messages in %v (avg %v per message)", numMessages, duration, duration/time.Duration(numMessages))

	// Verify that the server didn't crash and processed messages
	time.Sleep(500 * time.Millisecond)
	msgs := mockDAW.GetMessages()
	if len(msgs) < numMessages/2 { // Allow some network drop in stress but check for high throughput
		t.Errorf("Too few messages received: %d/%d", len(msgs), numMessages)
	}
}

func TestStress_ConcurrentDAWs(t *testing.T) {
	binPath := "./superdaw-concurrent-test"
	exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw").Run()
	defer os.Remove(binPath)

	client, _ := client.NewClient(binPath)
	defer client.Close()

	daws := []string{"ableton", "reaper", "ardour", "bitwig", "flstudio"}
	var wg sync.WaitGroup
	wg.Add(len(daws))

	for _, d := range daws {
		go func(dawName string) {
			defer wg.Done()
			for i := 0; i < 50; i++ {
				client.TransportControl(true, 120.0, dawName)
				time.Sleep(10 * time.Millisecond)
			}
		}(d)
	}
	wg.Wait()
}
