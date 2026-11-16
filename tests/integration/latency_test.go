package integration

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"testing"
	"time"
	"strings"

	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

// Benchmark plugin host tool dispatch latency.
func BenchmarkToolDispatch(b *testing.B) {
	binPath := "./superdaw-bench-test"
	buildCmd := exec.Command("go", "build", "-o", binPath, "../../cmd/superdaw")
	if out, err := buildCmd.CombinedOutput(); err != nil {
		b.Fatalf("Failed to build: %v\nOutput: %s", err, string(out))
	}
	defer os.Remove(binPath)

	cmd := exec.Command(binPath)
	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	if err := cmd.Start(); err != nil {
		b.Fatalf("Failed to start cmd: %v", err)
	}
	defer cmd.Process.Kill()

	time.Sleep(200 * time.Millisecond)
	decoder := json.NewDecoder(stdout)

	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  "tools/call",
		Params:  json.RawMessage(`{"name": "superdaw_list_plugins", "arguments": {}}`),
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req.ID = i
		data, _ := json.Marshal(req)
		fmt.Fprintf(stdin, "%s\n", string(data))

		var res mcp.JSONRPCResponse
		// We have to wait for the correct ID, sometimes other background logs occur
		for {
			if err := decoder.Decode(&res); err != nil {
				b.Fatalf("Decode failed at %d: %v", i, err)
			}

			// We handle integer IDs now
			idStr := fmt.Sprintf("%v", res.ID)
			if strings.Contains(idStr, fmt.Sprintf("%d", i)) {
				break
			}
		}
	}
}
