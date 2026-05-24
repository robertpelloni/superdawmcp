package daw

import (
	"fmt"
	"sync"
	"github.com/hypebeast/go-osc/osc"
)

type AbletonLiveDriver struct {
	OSCClient     *osc.Client
	notifyHandler func(method string, params interface{})
	state         struct {
		playing bool
		tempo   float32
		mu      sync.RWMutex
	}
}

func NewAbletonDriver(host string, port int, localPort int) *AbletonLiveDriver {
	d := &AbletonLiveDriver{OSCClient: osc.NewClient(host, port)}

	// Start OSC listener for state updates
	go d.listen(localPort)

	return d
}

func (a *AbletonLiveDriver) listen(port int) {
	dispatcher := osc.NewStandardDispatcher()
	dispatcher.AddMsgHandler("/superdaw/state/playing", func(msg *osc.Message) {
		a.state.mu.Lock()
		if len(msg.Arguments) > 0 {
			if b, ok := msg.Arguments[0].(bool); ok {
				a.state.playing = b
				if a.notifyHandler != nil {
					a.notifyHandler("superdaw/transport_update", map[string]interface{}{
						"daw":     "ableton",
						"playing": b,
					})
				}
			}
		}
		a.state.mu.Unlock()
	})
	dispatcher.AddMsgHandler("/superdaw/state/tempo", func(msg *osc.Message) {
		a.state.mu.Lock()
		if len(msg.Arguments) > 0 {
			if f, ok := msg.Arguments[0].(float32); ok {
				a.state.tempo = f
				if a.notifyHandler != nil {
					a.notifyHandler("superdaw/transport_update", map[string]interface{}{
						"daw":   "ableton",
						"tempo": f,
					})
				}
			}
		}
		a.state.mu.Unlock()
	})

	server := &osc.Server{Addr: fmt.Sprintf("127.0.0.1:%d", port), Dispatcher: dispatcher}
	server.ListenAndServe()
}

func (a *AbletonLiveDriver) Connect(endpoint string) error { return nil }
func (a *AbletonLiveDriver) Disconnect() error { return nil }

func (a *AbletonLiveDriver) SetTransportState(playing bool, bpm float64) error {
	m := osc.NewMessage("/superdaw/transport/play")
	v := int32(0); if playing { v = 1 }
	m.Append(v)
	a.OSCClient.Send(m)
	m2 := osc.NewMessage("/superdaw/transport/tempo")
	m2.Append(float32(bpm))
	return a.OSCClient.Send(m2)
}

func (a *AbletonLiveDriver) GetTransportState() (bool, float64, error) {
	a.state.mu.RLock()
	defer a.state.mu.RUnlock()
	return a.state.playing, float64(a.state.tempo), nil
}

func (a *AbletonLiveDriver) GetTracks() ([]TrackConfig, error) {
	return []TrackConfig{}, nil
}

func (a *AbletonLiveDriver) CreateTrack(name, trackType string) (string, error) {
	m := osc.NewMessage("/superdaw/track/create")
	m.Append(name)
	m.Append(trackType)
	return "id", a.OSCClient.Send(m)
}

func (a *AbletonLiveDriver) SetTrackVolume(id string, volume float32) error {
	m := osc.NewMessage("/superdaw/track/volume")
	m.Append(id)
	m.Append(volume)
	return a.OSCClient.Send(m)
}

func (a *AbletonLiveDriver) SetTrackPan(id string, pan float32) error {
	m := osc.NewMessage("/superdaw/track/pan")
	m.Append(id)
	m.Append(pan)
	return a.OSCClient.Send(m)
}

func (a *AbletonLiveDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	m := osc.NewMessage("/superdaw/clip/write")
	m.Append(id)
	m.Append(int32(idx))
	for _, n := range notes {
		m.Append(int32(n.Pitch))
		m.Append(int32(n.Velocity))
		m.Append(n.StartBeat)
		m.Append(n.Duration)
	}
	return a.OSCClient.Send(m)
}

func (a *AbletonLiveDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }

func (a *AbletonLiveDriver) DeleteClip(id string, idx int) error {
	m := osc.NewMessage("/superdaw/clip/delete")
	m.Append(id)
	m.Append(int32(idx))
	return a.OSCClient.Send(m)
}

func (a *AbletonLiveDriver) SetNotifyHandler(handler func(method string, params interface{})) {
	a.notifyHandler = handler
}

func (a *AbletonLiveDriver) ExecuteCustomCommand(cmd string, args map[string]interface{}) (interface{}, error) {
	switch cmd {
	case "fire_scene":
		idx, ok := args["scene_index"].(float64)
		if !ok { return nil, fmt.Errorf("missing scene_index") }
		m := osc.NewMessage("/live/scene/fire")
		m.Append(int32(idx))
		return "Fired scene", a.OSCClient.Send(m)
	}
	return nil, fmt.Errorf("unknown command: %s", cmd)
}
