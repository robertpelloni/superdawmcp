package daw
import ( "fmt"; "github.com/hypebeast/go-osc/osc" )
type ArdourDriver struct {
	OSCClient *osc.Client
}

func NewArdourDriver(h string, p int) *ArdourDriver { return &ArdourDriver{OSCClient: osc.NewClient(h, p)} }
func (a *ArdourDriver) Connect(e string) error { return nil }
func (a *ArdourDriver) Disconnect() error { return nil }
func (a *ArdourDriver) SetTransportState(p bool, b float64) error {
	addr := "/transport_stop"; if p { addr = "/transport_play" }; return a.OSCClient.Send(osc.NewMessage(addr))
}
func (a *ArdourDriver) GetTransportState() (bool, float64, error) { return false, 120.0, nil }
func (a *ArdourDriver) GetTracks() ([]TrackConfig, error)           { return []TrackConfig{}, nil }
func (a *ArdourDriver) CreateTrack(n, t string) (string, error) {
	m := osc.NewMessage("/access_action"); m.Append("Track/add-audio-track"); return "id", a.OSCClient.Send(m)
}
func (a *ArdourDriver) SetTrackVolume(id string, v float32) error {
	m := osc.NewMessage("/strip/fader"); m.Append(id); m.Append(v); return a.OSCClient.Send(m)
}
func (a *ArdourDriver) SetTrackPan(id string, p float32) error {
	m := osc.NewMessage("/strip/pan_stereo_pan"); m.Append(id); m.Append(p); return a.OSCClient.Send(m)
}
func (a *ArdourDriver) WriteMIDIClip(id string, idx int, n []MIDINote) error { return nil }
func (a *ArdourDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }
func (a *ArdourDriver) DeleteClip(id string, idx int) error { return nil }

func (a *ArdourDriver) SetNotifyHandler(handler func(method string, params interface{})) {
}

func (a *ArdourDriver) ExecuteCustomCommand(cmd string, args map[string]interface{}) (interface{}, error) {
	return nil, fmt.Errorf("custom commands not implemented for Ardour")
}
