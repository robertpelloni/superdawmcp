package compatibility

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/robertpelloni/superdaw-mcp/pkg/client/go"
)

func TestDAWCompatibility(t *testing.T) {
	binPath := "./superdaw-compatibility-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build server: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(binPath)

	daws := []string{"ableton", "reaper", "ardour", "bitwig", "flstudio", "logic", "cubase", "protools"}

	for _, dawName := range daws {
		t.Run(dawName, func(t *testing.T) {
			client, err := client.NewClient(binPath)
			if err != nil { t.Fatalf("Failed to start client: %v", err) }
			defer client.Close()

			// Test Transport
			if err := client.TransportControl(true, 128.0, dawName); err != nil {
				t.Errorf("%s: Transport Play failed: %v", dawName, err)
			}
			time.Sleep(100 * time.Millisecond)

			// Test Track Creation
			if err := client.CreateTrack("TestTrack", "midi", dawName); err != nil {
				t.Errorf("%s: CreateTrack failed: %v", dawName, err)
			}

			// Test Mixer
			if err := client.SetMixer("1", 0.75, 0.0, dawName); err != nil {
				t.Errorf("%s: SetMixer failed: %v", dawName, err)
			}
		})
	}
}
