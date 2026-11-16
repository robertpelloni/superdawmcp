package integration

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"testing"
	"os"
	"time"
	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

func TestIntegration_StressTest(t *testing.T) {
	binPath := "./superdaw-stress-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("Failed to build: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(binPath)

	cmd := exec.Command(binPath)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start cmd: %v", err)
	}
	defer func() {
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
	}()

	decoder := json.NewDecoder(stdout)

	start := time.Now()
	iterations := 500

	for i := 0; i < iterations; i++ {
		req := mcp.JSONRPCRequest{
			JSONRPC: "2.0",
			Method:  "tools/call",
			Params:  json.RawMessage(fmt.Sprintf(`{"name": "superdaw_set_mixer", "arguments": {"track_id": "%d", "volume": 0.5}}`, i)),
			ID:      i,
		}
		data, _ := json.Marshal(req)
		fmt.Fprintf(stdin, "%s\n", string(data))

		var res mcp.JSONRPCResponse
		if err := decoder.Decode(&res); err != nil {
			t.Fatalf("Decode failed at %d: %v", i, err)
		}
	}

	duration := time.Since(start)
	t.Logf("Processed %d tool calls in %v (%f req/sec)", iterations, duration, float64(iterations)/duration.Seconds())

	if duration.Seconds() > 5.0 {
		t.Errorf("Performance below threshold: %v for %d calls", duration, iterations)
	}
}
