package daw

import (
	"encoding/json"

	"github.com/hypebeast/go-osc/osc"
)

type ReaperDriver struct {
	OSCClient *osc.Client
	WebPort   int
	Host      string
}

func NewReaperDriver(host string, oscPort int, webPort int) *ReaperDriver {
	return &ReaperDriver{
		OSCClient: osc.NewClient(host, oscPort),
		WebPort:   webPort,
		Host:      host,
	}
}

func (r *ReaperDriver) Connect(endpoint string) error {
	return nil
}

func (r *ReaperDriver) Disconnect() error {
	return nil
}

func (r *ReaperDriver) SetTransportState(playing bool, bpm float64) error {
	var cmd int32 = 0
	if playing {
		cmd = 1
	}

	// Unified Protocol: /superdaw/transport/play
	msg1 := osc.NewMessage("/superdaw/transport/play")
	msg1.Append(cmd)
	if err := r.OSCClient.Send(msg1); err != nil {
		return err
	}

	// Unified Protocol: /superdaw/transport/tempo
	tempoMsg := osc.NewMessage("/superdaw/transport/tempo")
	tempoMsg.Append(float32(bpm))
	return r.OSCClient.Send(tempoMsg)
}

func (r *ReaperDriver) GetTransportState() (bool, float64, error) {
	return false, 120.0, nil
}

func (r *ReaperDriver) GetTracks() ([]TrackConfig, error) {
	return []TrackConfig{}, nil
}

func (r *ReaperDriver) CreateTrack(name string, trackType string) (string, error) {
	// Unified Protocol: /superdaw/track/create
	msg := osc.NewMessage("/superdaw/track/create")
	msg.Append(name)
	msg.Append(trackType)
	err := r.OSCClient.Send(msg)
	return "reaper_track_new", err
}

func (r *ReaperDriver) SetTrackVolume(trackID string, volume float32) error {
	// Unified Protocol: /superdaw/track/volume
	msg := osc.NewMessage("/superdaw/track/volume")
	msg.Append(trackID)
	msg.Append(volume)
	return r.OSCClient.Send(msg)
}

func (r *ReaperDriver) WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error {
	// Unified Protocol: /superdaw/clip/write
	msg := osc.NewMessage("/superdaw/clip/write")
	msg.Append(trackID)
	msg.Append(int32(clipIndex))

	payload, _ := json.Marshal(notes)
	msg.Append(string(payload))

	return r.OSCClient.Send(msg)
}

func (r *ReaperDriver) InstantiatePlugin(trackID string, pluginName string) (string, error) {
	// Unified Protocol: /superdaw/device/instantiate
	msg := osc.NewMessage("/superdaw/device/instantiate")
	msg.Append(trackID)
	msg.Append(pluginName)
	err := r.OSCClient.Send(msg)
	return "reaper_fx_id", err
}

func (r *ReaperDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	// Unified Protocol: /superdaw/device/param
	msg := osc.NewMessage("/superdaw/device/param")
	msg.Append(trackID)
	msg.Append(pluginID)
	msg.Append(int32(paramIndex))
	msg.Append(value)
	return r.OSCClient.Send(msg)
}

func (r *ReaperDriver) LuaBridgeCall(funcName string, args []interface{}) error {
	return nil
}
