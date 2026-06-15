package daw

import (
	"encoding/json"
	"fmt"
	"sync"
	"github.com/hypebeast/go-osc/osc"
)

type AbletonLiveDriver struct {
	OSCClient     *osc.Client
	notifyHandler func(method string, params interface{})
	state         struct {
		playing     bool
		tempo       float32
		arrangement string
		trackCount  int
		mu          sync.RWMutex
	}
}

func NewAbletonDriver(host string, port int, localPort int) *AbletonLiveDriver {
	d := &AbletonLiveDriver{OSCClient: osc.NewClient(host, port)}
	go d.listen(localPort)
	return d
}

func (a *AbletonLiveDriver) listen(port int) {
	dispatcher := osc.NewStandardDispatcher()
	dispatcher.AddMsgHandler("/superdaw/state/playing", func(msg *osc.Message) {
		a.state.mu.Lock()
		defer a.state.mu.Unlock()
		if len(msg.Arguments) > 0 {
			var b bool
			if val, ok := msg.Arguments[0].(bool); ok {
				b = val
			} else if val, ok := msg.Arguments[0].(int32); ok {
				b = val == 1
			} else {
				return
			}
			a.state.playing = b
			if a.notifyHandler != nil {
				a.notifyHandler("superdaw/transport_update", map[string]interface{}{"daw": "ableton", "playing": b})
			}
		}
	})
	dispatcher.AddMsgHandler("/superdaw/state/arrangement", func(msg *osc.Message) {
		a.state.mu.Lock()
		defer a.state.mu.Unlock()
		if len(msg.Arguments) > 0 {
			if s, ok := msg.Arguments[0].(string); ok {
				a.state.arrangement = s
				if a.notifyHandler != nil {
					a.notifyHandler("superdaw/arrangement_update", map[string]interface{}{"daw": "ableton", "arrangement": s})
				}
			}
		}
	})
	dispatcher.AddMsgHandler("/superdaw/state/tempo", func(msg *osc.Message) {
		a.state.mu.Lock()
		defer a.state.mu.Unlock()
		if len(msg.Arguments) > 0 {
			if f, ok := msg.Arguments[0].(float32); ok {
				a.state.tempo = f
				if a.notifyHandler != nil {
					a.notifyHandler("superdaw/transport_update", map[string]interface{}{"daw": "ableton", "tempo": f})
				}
			}
		}
	})
	dispatcher.AddMsgHandler("/superdaw/state/track_count", func(msg *osc.Message) {
		a.state.mu.Lock()
		defer a.state.mu.Unlock()
		if len(msg.Arguments) > 0 {
			if i, ok := msg.Arguments[0].(int32); ok {
				a.state.trackCount = int(i)
			}
		}
	})
	dispatcher.AddMsgHandler("/superdaw/state/track_created", func(msg *osc.Message) {
		a.state.mu.Lock()
		defer a.state.mu.Unlock()
		if len(msg.Arguments) > 0 {
			if i, ok := msg.Arguments[0].(int32); ok {
				a.state.trackCount = int(i) + 1
			}
		}
	})

	server := &osc.Server{Addr: fmt.Sprintf("0.0.0.0:%d", port), Dispatcher: dispatcher}
	server.ListenAndServe()
}

func (a *AbletonLiveDriver) GetType() string               { return "ableton" }
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
	a.state.mu.RLock()
	count := a.state.trackCount
	a.state.mu.RUnlock()
	tracks := make([]TrackConfig, count)
	for i := 0; i < count; i++ {
		tracks[i] = TrackConfig{ID: fmt.Sprintf("%d", i), Name: fmt.Sprintf("Track %d", i)}
	}
	return tracks, nil
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
	data, _ := json.Marshal(notes)
	m.Append(string(data))
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

func (a *AbletonLiveDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	m := osc.NewMessage("/superdaw/plugin/param")
	m.Append(trackID)
	m.Append(pluginID)
	m.Append(int32(paramIndex))
	m.Append(value)
	return a.OSCClient.Send(m)
}

func (a *AbletonLiveDriver) SendCC(trackID string, controller int, value int) error {
	m := osc.NewMessage("/superdaw/midi/cc")
	m.Append(trackID)
	m.Append(int32(controller))
	m.Append(int32(value))
	return a.OSCClient.Send(m)
}
