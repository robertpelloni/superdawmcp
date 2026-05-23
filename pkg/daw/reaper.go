package daw
import ( "encoding/json"; "github.com/hypebeast/go-osc/osc" )
type ReaperDriver struct { OSCClient *osc.Client }
func NewReaperDriver(h string, p, w int) *ReaperDriver { return &ReaperDriver{OSCClient: osc.NewClient(h, p)} }
func (r *ReaperDriver) Connect(e string) error { return nil }
func (r *ReaperDriver) Disconnect() error { return nil }
func (r *ReaperDriver) SetTransportState(p bool, b float64) error {
	m := osc.NewMessage("/superdaw/transport/play"); v := int32(0); if p { v = 1 }; m.Append(v); return r.OSCClient.Send(m)
}
func (r *ReaperDriver) CreateTrack(n, t string) (string, error) {
	m := osc.NewMessage("/superdaw/track/create"); m.Append(n); m.Append(t); return "id", r.OSCClient.Send(m)
}
func (r *ReaperDriver) SetTrackVolume(id string, v float32) error {
	m := osc.NewMessage("/superdaw/track/volume"); m.Append(id); m.Append(v); return r.OSCClient.Send(m)
}
func (r *ReaperDriver) SetTrackPan(id string, p float32) error {
	m := osc.NewMessage("/superdaw/track/pan"); m.Append(id); m.Append(p); return r.OSCClient.Send(m)
}
func (r *ReaperDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	m := osc.NewMessage("/superdaw/clip/write"); m.Append(id); m.Append(int32(idx))
	p, _ := json.Marshal(notes); m.Append(string(p)); return r.OSCClient.Send(m)
}
func (r *ReaperDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }
func (r *ReaperDriver) DeleteClip(id string, idx int) error {
	m := osc.NewMessage("/superdaw/clip/delete"); m.Append(id); m.Append(int32(idx)); return r.OSCClient.Send(m)
}
