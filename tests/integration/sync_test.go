package integration

import (
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

	time.Sleep(100 * time.Millisecond) // Wait for server to start

	// 2. Simulate OSC message from Ableton Agent
	client := osc.NewClient("127.0.0.1", localPort)
	msg := osc.NewMessage("/superdaw/state/playing")
	msg.Append(int32(1))

	cmd := exec.Command(binPath)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start cmd: %v", err)
	}
	defer cmd.Process.Kill()

	// 3. Wait for async processing
	time.Sleep(200 * time.Millisecond)

	// Send multiple rapid sync states to test race conditions
	for i := 0; i < 100; i++ {
		_ = mockDAW.SendMessage("127.0.0.1", 11001, "/superdaw/state/tempo", float32(120.0+float32(i)))
	}

	time.Sleep(200 * time.Millisecond)

	time.Sleep(200 * time.Millisecond)
	_, tempo, _ := driver.GetTransportState()
	if tempo != 135.0 {
		t.Errorf("Driver did not reflect real-time tempo change. Expected 135, got %f", tempo)
	}
}
