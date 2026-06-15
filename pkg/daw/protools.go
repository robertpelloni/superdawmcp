package daw

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

type ProToolsDriver struct {
	OSCClient *osc.Client
	Host      string
	Port      int
	Notify    func(method string, params interface{})
}

func NewProToolsDriver(host string, port int) *ProToolsDriver {
	return &ProToolsDriver{
		OSCClient: osc.NewClient(host, port),
		Host:      host,
		Port:      port,
	}
}

func (p *ProToolsDriver) GetType() string { return "protools" }
func (p *ProToolsDriver) Connect(endpoint string) error {
	return nil
}

func (p *ProToolsDriver) Disconnect() error {
	return nil
}

func (p *ProToolsDriver) SetTransportState(playing bool, bpm float64) error {
	var cmd int32 = 0
	if playing { cmd = 1 }
	msg := osc.NewMessage("/protools/transport/play")
	msg.Append(cmd)
	return p.OSCClient.Send(msg)
}

func (p *ProToolsDriver) GetTransportState() (bool, float64, error) {
	return false, 120.0, nil
}

func (p *ProToolsDriver) SetTrackVolume(trackID string, volume float32) error {
	msg := osc.NewMessage("/protools/track/volume")
	msg.Append(trackID)
	msg.Append(volume)
	return p.OSCClient.Send(msg)
}

func (p *ProToolsDriver) SetTrackPan(trackID string, pan float32) error {
	msg := osc.NewMessage("/protools/track/pan")
	msg.Append(trackID)
	msg.Append(pan)
	return p.OSCClient.Send(msg)
}

func (p *ProToolsDriver) SetTrackInstrument(trackID string, instrument string) error {
	return nil // Not implemented for Pro Tools
}

func (p *ProToolsDriver) CreateTrack(name string, trackType string) (string, error) {
	msg := osc.NewMessage("/protools/track/create")
	msg.Append(name)
	msg.Append(trackType)
	err := p.OSCClient.Send(msg)
	return "pt_track_" + name, err
}

func (p *ProToolsDriver) WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error {
	return fmt.Errorf("MIDI clip injection not natively supported in Pro Tools via OSC bridge")
}

func (p *ProToolsDriver) ListClips(trackID string) ([]ClipInfo, error) {
	return []ClipInfo{}, nil
}

func (p *ProToolsDriver) DeleteClip(trackID string, clipIndex int) error {
	return nil
}

func (p *ProToolsDriver) GetTracks() ([]TrackConfig, error) {
	return []TrackConfig{}, nil
}

func (p *ProToolsDriver) InstantiatePlugin(trackID string, pluginName string) (string, error) {
	return "", nil
}

func (p *ProToolsDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	return nil
}

func (p *ProToolsDriver) SendCC(trackID string, controller int, value int) error {
	return nil
}

func (p *ProToolsDriver) ExecuteCustomCommand(command string, args map[string]interface{}) (interface{}, error) {
	return nil, nil
}

func (p *ProToolsDriver) SetNotifyHandler(h func(method string, params interface{})) {
	p.Notify = h
}

// SendCC stub – Pro Tools does not support CC via OSC.
func (p *ProToolsDriver) SendCC(trackID string, controller int, value int) error { return nil }

