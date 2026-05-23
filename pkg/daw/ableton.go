package daw
import "github.com/hypebeast/go-osc/osc"
type AbletonLiveDriver struct { OSCClient *osc.Client }
func NewAbletonDriver(h string, p, l int) *AbletonLiveDriver { return &AbletonLiveDriver{OSCClient: osc.NewClient(h, p)} }
func (a *AbletonLiveDriver) Connect(e string) error { return nil }
func (a *AbletonLiveDriver) Disconnect() error { return nil }
func (a *AbletonLiveDriver) SetTransportState(p bool, b float64) error {
	m := osc.NewMessage("/superdaw/transport/play"); v := int32(0); if p { v = 1 }; m.Append(v); a.OSCClient.Send(m)
	m2 := osc.NewMessage("/superdaw/transport/tempo"); m2.Append(float32(b)); return a.OSCClient.Send(m2)
}
func (a *AbletonLiveDriver) CreateTrack(n, t string) (string, error) {
	m := osc.NewMessage("/superdaw/track/create"); m.Append(n); m.Append(t); return "id", a.OSCClient.Send(m)
}
func (a *AbletonLiveDriver) SetTrackVolume(id string, v float32) error {
	m := osc.NewMessage("/superdaw/track/volume"); m.Append(id); m.Append(v); return a.OSCClient.Send(m)
}
func (a *AbletonLiveDriver) SetTrackPan(id string, p float32) error {
	m := osc.NewMessage("/superdaw/track/pan"); m.Append(id); m.Append(p); return a.OSCClient.Send(m)
}
func (a *AbletonLiveDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	m := osc.NewMessage("/superdaw/clip/write"); m.Append(id); m.Append(int32(idx))
	for _, n := range notes { m.Append(int32(n.Pitch)); m.Append(int32(n.Velocity)); m.Append(n.StartBeat); m.Append(n.Duration) }
	return a.OSCClient.Send(m)
}
func (a *AbletonLiveDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }
func (a *AbletonLiveDriver) DeleteClip(id string, idx int) error {
	m := osc.NewMessage("/superdaw/clip/delete"); m.Append(id); m.Append(int32(idx)); return a.OSCClient.Send(m)
}
