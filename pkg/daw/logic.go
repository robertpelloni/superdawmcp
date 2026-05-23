package daw

import (
	"github.com/hypebeast/go-osc/osc"
)

// LogicProDriver maps unified commands to Logic Pro's standard OSC schema.
type LogicProDriver struct {
	OSCClient *osc.Client
}

func NewLogicProDriver(host string, port int) *LogicProDriver {
	return &LogicProDriver{OSCClient: osc.NewClient(host, port)}
}

func (l *LogicProDriver) Connect(endpoint string) error { return nil }
func (l *LogicProDriver) Disconnect() error { return nil }

func (l *LogicProDriver) SetTransportState(playing bool, bpm float64) error {
	addr := "/stop"; if playing { addr = "/play" }
	l.OSCClient.Send(osc.NewMessage(addr))

	m := osc.NewMessage("/tempo")
	m.Append(float32(bpm))
	return l.OSCClient.Send(m)
}

func (l *LogicProDriver) GetTransportState() (bool, float64, error) { return false, 120.0, nil }

func (l *LogicProDriver) CreateTrack(name, trackType string) (string, error) {
	m := osc.NewMessage("/shortcut/create_track")
	return "logic_track", l.OSCClient.Send(m)
}

func (l *LogicProDriver) SetTrackVolume(id string, volume float32) error {
	m := osc.NewMessage("/fader" + id)
	m.Append(volume)
	return l.OSCClient.Send(m)
}

func (l *LogicProDriver) SetTrackPan(id string, pan float32) error {
	m := osc.NewMessage("/pan" + id)
	m.Append(pan)
	return l.OSCClient.Send(m)
}

func (l *LogicProDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	return nil
}

func (l *LogicProDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }
func (l *LogicProDriver) DeleteClip(id string, idx int) error { return nil }

func (l *LogicProDriver) ExecuteCustomCommand(cmd string, args map[string]interface{}) (interface{}, error) {
	m := osc.NewMessage("/logic/custom/" + cmd)
	return "Sent to Logic Pro", l.OSCClient.Send(m)
}
