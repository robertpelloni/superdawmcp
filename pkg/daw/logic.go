package daw

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

// LogicProDriver maps unified commands to Logic Pro's standard OSC schema.
type LogicProDriver struct {
	OSCClient     *osc.Client
	notifyHandler func(method string, params interface{})
}

func NewLogicProDriver(host string, port int, localPort int) *LogicProDriver {
	d := &LogicProDriver{OSCClient: osc.NewClient(host, port)}
	if localPort > 0 {
		go d.listen(localPort)
	}
	return d
}

func (l *LogicProDriver) listen(port int) {
	dispatcher := osc.NewStandardDispatcher()
	handler := func(msg *osc.Message) {
		if l.notifyHandler == nil {
			return
		}
		if msg.Address == "/superdaw/state/playing" {
			if len(msg.Arguments) > 0 {
				if b, ok := msg.Arguments[0].(int32); ok {
					l.notifyHandler("superdaw/transport_update", map[string]interface{}{"daw": "logic", "playing": b == 1})
				}
			}
		}
	}
	dispatcher.AddMsgHandler("/superdaw/state/playing", handler)
	server := &osc.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Dispatcher: dispatcher}
	server.ListenAndServe()
}

func (l *LogicProDriver) GetType() string               { return "logic" }
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

func (l *LogicProDriver) GetTracks() ([]TrackConfig, error) { return []TrackConfig{}, nil }

func (l *LogicProDriver) CreateTrack(name, trackType string) (string, error) {
	// Trigger shortcut mapping in Logic
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

func (l *LogicProDriver) SetTrackInstrument(id string, instrument string) error {
	m := osc.NewMessage("/logic/track/instrument")
	m.Append(id)
	m.Append(instrument)
	return l.OSCClient.Send(m)
}

func (l *LogicProDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error { return nil }
func (l *LogicProDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }
func (l *LogicProDriver) DeleteClip(id string, idx int) error { return nil }

func (l *LogicProDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	return nil
}

func (l *LogicProDriver) SendCC(trackID string, controller int, value int) error {
	return nil
}

func (l *LogicProDriver) SetNotifyHandler(handler func(method string, params interface{})) {
	l.notifyHandler = handler
}

func (l *LogicProDriver) ExecuteCustomCommand(cmd string, args map[string]interface{}) (interface{}, error) {
	m := osc.NewMessage("/logic/custom/" + cmd)
	return "Sent to Logic Pro", l.OSCClient.Send(m)
}
