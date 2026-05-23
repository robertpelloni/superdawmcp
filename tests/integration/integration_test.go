package integration
import ( "context"; "encoding/json"; "os"; "os/exec"; "testing"; "time"; "github.com/robertpelloni/superdaw-mcp/pkg/mcp" )
func TestIntegration_EndToEnd(t *testing.T) {
	bin := "./superdaw-mcp-test"; exec.Command("go", "build", "-o", bin, "../../cmd/superdaw/main.go").Run(); defer os.Remove(bin)
	mock := NewMockDAW(11000); go mock.Start(); defer mock.Stop()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second); defer cancel()
	cmd := exec.CommandContext(ctx, bin); stdin, _ := cmd.StdinPipe(); stdout, _ := cmd.StdoutPipe(); cmd.Start()
	writer := json.NewEncoder(stdin); reader := json.NewDecoder(stdout)
	req := mcp.JSONRPCRequest{JSONRPC: "2.0", Method: "tools/call", Params: json.RawMessage(`{"name":"superdaw_set_mixer","arguments":{"track_id":"1","volume":0.5}}`), ID: 1}
	writer.Encode(req); var res mcp.JSONRPCResponse; reader.Decode(&res)
	time.Sleep(100 * time.Millisecond); msgs := mock.GetMessages()
	found := false; for _, m := range msgs { if m.Address == "/superdaw/track/volume" { found = true; break } }; if !found { t.Error("No OSC volume message received") }
	stdin.Close(); cmd.Wait()
}
