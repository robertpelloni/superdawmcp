package daw
import (
	"github.com/hypebeast/go-osc/osc"
)
type ArdourDriver struct { OSCClient *osc.Client }
func NewArdourDriver(h string, p int) *ArdourDriver { return &ArdourDriver{OSCClient: osc.NewClient(h, p)} }
func (a *ArdourDriver) Connect(e string) error { return nil }
func (a *ArdourDriver) Disconnect() error { return nil }
func (a *ArdourDriver) SetTransportState(p bool, bpm float64) error {
	addr := "/transport_stop"; if p { addr = "/transport_play" }
	return a.OSCClient.Send(osc.NewMessage(addr))
}
func (a *ArdourDriver) CreateTrack(n, t string) (string, error) {
	msg := osc.NewMessage("/access_action"); msg.Append("Track/add-audio-track")
	return "id", a.OSCClient.Send(msg)
}
func (a *ArdourDriver) SetTrackVolume(id string, v float32) error {
	msg := osc.NewMessage("/strip/fader"); msg.Append(id); msg.Append(v)
	return a.OSCClient.Send(msg)
}
func (a *ArdourDriver) SetTrackPan(id string, p float32) error {
	msg := osc.NewMessage("/strip/pan_stereo_pan"); msg.Append(id); msg.Append(p)
	return a.OSCClient.Send(msg)
}
func (a *ArdourDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error { return nil }
