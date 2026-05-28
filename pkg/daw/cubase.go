package daw

import (
	"github.com/hypebeast/go-osc/osc"
)

// CubaseDriver maps unified commands to Cubase's MIDI Remote / OSC bridge.
type CubaseDriver struct {
	OSCClient *osc.Client
}

func NewCubaseDriver(host string, port int) *CubaseDriver {
	return &CubaseDriver{OSCClient: osc.NewClient(host, port)}
}

func (c *CubaseDriver) GetType() string               { return "cubase" }
func (c *CubaseDriver) Connect(endpoint string) error { return nil }
func (c *CubaseDriver) Disconnect() error { return nil }

func (c *CubaseDriver) SetTransportState(playing bool, bpm float64) error {
	addr := "/cubase/transport/stop"; if playing { addr = "/cubase/transport/start" }
	return c.OSCClient.Send(osc.NewMessage(addr))
}

func (c *CubaseDriver) GetTransportState() (bool, float64, error) { return false, 120.0, nil }

func (c *CubaseDriver) GetTracks() ([]TrackConfig, error) { return []TrackConfig{}, nil }

func (c *CubaseDriver) CreateTrack(name, trackType string) (string, error) {
	m := osc.NewMessage("/cubase/track/add")
	m.Append(name)
	m.Append(trackType)
	return "cubase_track", c.OSCClient.Send(m)
}

func (c *CubaseDriver) SetTrackVolume(id string, volume float32) error {
	m := osc.NewMessage("/cubase/track/volume")
	m.Append(id)
	m.Append(volume)
	return c.OSCClient.Send(m)
}

func (c *CubaseDriver) SetTrackPan(id string, pan float32) error {
	m := osc.NewMessage("/cubase/track/pan")
	m.Append(id)
	m.Append(pan)
	return c.OSCClient.Send(m)
}

func (c *CubaseDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error { return nil }
func (c *CubaseDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }
func (c *CubaseDriver) DeleteClip(id string, idx int) error { return nil }

func (c *CubaseDriver) SetNotifyHandler(handler func(method string, params interface{})) {
}

func (c *CubaseDriver) ExecuteCustomCommand(cmd string, args map[string]interface{}) (interface{}, error) {
	m := osc.NewMessage("/cubase/custom/" + cmd)
	return "Sent to Cubase", c.OSCClient.Send(m)
}

func (c *CubaseDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	return nil
}

func (c *CubaseDriver) SendCC(trackID string, controller int, value int) error {
	return nil
}
