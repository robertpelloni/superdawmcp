package daw

import (
	"github.com/hypebeast/go-osc/osc"
)

// FLStudioDriver uses a MIDI-over-OSC bridge to communicate with FL Studio's MIDI scripting API.
type FLStudioDriver struct {
	OSCClient *osc.Client
}

func NewFLStudioDriver(host string, port int) *FLStudioDriver {
	return &FLStudioDriver{OSCClient: osc.NewClient(host, port)}
}

func (f *FLStudioDriver) GetType() string               { return "flstudio" }
func (f *FLStudioDriver) Connect(endpoint string) error { return nil }
func (f *FLStudioDriver) Disconnect() error { return nil }

func (f *FLStudioDriver) SetTransportState(playing bool, bpm float64) error {
	m := osc.NewMessage("/flstudio/transport/play")
	v := int32(0); if playing { v = 1 }
	m.Append(v)
	f.OSCClient.Send(m)

	m2 := osc.NewMessage("/flstudio/transport/tempo")
	m2.Append(float32(bpm))
	return f.OSCClient.Send(m2)
}

func (f *FLStudioDriver) GetTransportState() (bool, float64, error) { return false, 120.0, nil }

func (f *FLStudioDriver) GetTracks() ([]TrackConfig, error) { return []TrackConfig{}, nil }

func (f *FLStudioDriver) CreateTrack(name, trackType string) (string, error) {
	m := osc.NewMessage("/flstudio/track/create")
	m.Append(name)
	m.Append(trackType)
	return "id", f.OSCClient.Send(m)
}

func (f *FLStudioDriver) SetTrackVolume(id string, volume float32) error {
	m := osc.NewMessage("/flstudio/track/volume")
	m.Append(id)
	m.Append(volume)
	return f.OSCClient.Send(m)
}

func (f *FLStudioDriver) SetTrackPan(id string, pan float32) error {
	m := osc.NewMessage("/flstudio/track/pan")
	m.Append(id)
	m.Append(pan)
	return f.OSCClient.Send(m)
}

func (f *FLStudioDriver) SetTrackInstrument(id string, instrument string) error {
	return nil // Not implemented for FL Studio
}

func (f *FLStudioDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	m := osc.NewMessage("/flstudio/clip/write")
	m.Append(id)
	m.Append(int32(idx))
	for _, n := range notes {
		m.Append(int32(n.Pitch))
		m.Append(int32(n.Velocity))
	}
	return f.OSCClient.Send(m)
}

func (f *FLStudioDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }
func (f *FLStudioDriver) DeleteClip(id string, idx int) error { return nil }

func (f *FLStudioDriver) SetNotifyHandler(handler func(method string, params interface{})) {
}

func (f *FLStudioDriver) ExecuteCustomCommand(cmd string, args map[string]interface{}) (interface{}, error) {
	m := osc.NewMessage("/flstudio/custom/" + cmd)
	return "Sent to FL Studio", f.OSCClient.Send(m)
}

func (f *FLStudioDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	return nil
}

func (f *FLStudioDriver) SendCC(trackID string, controller int, value int) error {
	return nil
}
