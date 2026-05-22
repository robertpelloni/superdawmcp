package daw
import (
	"github.com/hypebeast/go-osc/osc"
)
type AbletonLiveDriver struct {
	OSCClient *osc.Client
	OSCHost   string
	OSCPort   int
}
func NewAbletonDriver(host string, port int, localPort int) *AbletonLiveDriver {
	return &AbletonLiveDriver{OSCClient: osc.NewClient(host, port), OSCHost: host, OSCPort: port}
}
func (a *AbletonLiveDriver) Connect(e string) error { return nil }
func (a *AbletonLiveDriver) Disconnect() error { return nil }
func (a *AbletonLiveDriver) SetTransportState(p bool, bpm float64) error {
	msg := osc.NewMessage("/superdaw/transport/play")
	val := int32(0); if p { val = 1 }
	msg.Append(val)
	a.OSCClient.Send(msg)
	msg2 := osc.NewMessage("/superdaw/transport/tempo")
	msg2.Append(float32(bpm))
	return a.OSCClient.Send(msg2)
}
func (a *AbletonLiveDriver) CreateTrack(n, t string) (string, error) {
	msg := osc.NewMessage("/superdaw/track/create")
	msg.Append(n); msg.Append(t)
	return "id", a.OSCClient.Send(msg)
}
func (a *AbletonLiveDriver) SetTrackVolume(id string, v float32) error {
	msg := osc.NewMessage("/superdaw/track/volume")
	msg.Append(id); msg.Append(v)
	return a.OSCClient.Send(msg)
}
func (a *AbletonLiveDriver) SetTrackPan(id string, p float32) error {
	msg := osc.NewMessage("/superdaw/track/pan")
	msg.Append(id); msg.Append(p)
	return a.OSCClient.Send(msg)
}
func (a *AbletonLiveDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	msg := osc.NewMessage("/superdaw/clip/write")
	msg.Append(id); msg.Append(int32(idx))
	for _, n := range notes {
		msg.Append(int32(n.Pitch)); msg.Append(int32(n.Velocity)); msg.Append(n.StartBeat); msg.Append(n.Duration)
	}
	return a.OSCClient.Send(msg)
}
