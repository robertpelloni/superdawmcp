package client

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"github.com/robertpelloni/superdaw-mcp/pkg/mcp"
)

type Client struct {
	cmd     *exec.Cmd
	stdin   io.WriteCloser
	stdout  io.ReadCloser
	decoder *json.Decoder
	id      int
}

func NewClient(binPath string) (*Client, error) {
	cmd := exec.Command(binPath)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Client{
		cmd:     cmd,
		stdin:   stdin,
		stdout:  stdout,
		decoder: json.NewDecoder(stdout),
		id:      1,
	}, nil
}

func (c *Client) Close() error {
	c.stdin.Close()
	return c.cmd.Wait()
}

func (c *Client) call(method string, params interface{}) (interface{}, error) {
	paramsJSON, _ := json.Marshal(params)
	req := mcp.JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  paramsJSON,
		ID:      c.id,
	}
	c.id++

	reqJSON, _ := json.Marshal(req)
	fmt.Fprintf(c.stdin, "%s\n", string(reqJSON))

	var res mcp.JSONRPCResponse
	if err := c.decoder.Decode(&res); err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, fmt.Errorf("RPC error: %s (code %d)", res.Error.Message, res.Error.Code)
	}

	return res.Result, nil
}

func (c *Client) SetMixer(trackID string, volume float32, pan float32, targetDAW ...string) error {
	args := map[string]interface{}{"track_id": trackID, "volume": volume, "pan": pan}
	if len(targetDAW) > 0 {
		args["daw"] = targetDAW[0]
	}
	params := map[string]interface{}{
		"name":      "superdaw_set_mixer",
		"arguments": args,
	}
	_, err := c.call("tools/call", params)
	return err
}

func (c *Client) CreateTrack(name string, trackType string, targetDAW ...string) error {
	args := map[string]interface{}{"name": name, "type": trackType}
	if len(targetDAW) > 0 {
		args["daw"] = targetDAW[0]
	}
	params := map[string]interface{}{
		"name":      "superdaw_create_track",
		"arguments": args,
	}
	_, err := c.call("tools/call", params)
	return err
}

func (c *Client) TransportControl(playing bool, bpm float64, targetDAW ...string) error {
	args := map[string]interface{}{"playing": playing, "bpm": bpm}
	if len(targetDAW) > 0 {
		args["daw"] = targetDAW[0]
	}
	params := map[string]interface{}{
		"name":      "superdaw_transport_control",
		"arguments": args,
	}
	_, err := c.call("tools/call", params)
	return err
}
