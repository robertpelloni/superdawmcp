package daw

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	var addr string
	if playing {
		addr = "/transport/play"
	} else {
		addr = "/transport/stop"
	}
	msg := osc.NewMessage(addr)
	if err := r.OSCClient.Send(msg); err != nil {
		return err
	}

	// Set Tempo via Web API if possible, or OSC
	tempoMsg := osc.NewMessage("/tempo")
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
	// Using REAPER Web API (wwr) to trigger an action
	// 40001 is the command ID for 'Track: Insert new track'
	url := fmt.Sprintf("http://%s:%d/wwr/_40001", r.Host, r.WebPort)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "reaper_track_new", nil
}

func (r *ReaperDriver) SetTrackVolume(trackID string, volume float32) error {
	// OSC: /track/1/volume
	addr := fmt.Sprintf("/track/%s/volume", trackID)
	msg := osc.NewMessage(addr)
	msg.Append(volume)
	return r.OSCClient.Send(msg)
}

func (r *ReaperDriver) WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error {
	// REAPER requires more complex logic for MIDI instantiation usually via ReaScript.
	// We can use a custom OSC address if the bridge supports it.
	msg := osc.NewMessage("/custom/midi/write")
	msg.Append(trackID)

	payload, _ := json.Marshal(notes)
	msg.Append(string(payload))

	return r.OSCClient.Send(msg)
}

func (r *ReaperDriver) InstantiatePlugin(trackID string, pluginName string) (string, error) {
	// Trigger custom ReaScript via Web API
	url := fmt.Sprintf("http://%s:%d/wwr/instantiate_plugin?track=%s&plugin=%s", r.Host, r.WebPort, trackID, pluginName)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return "reaper_fx_id", nil
}

func (r *ReaperDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	// OSC: /track/1/fx/1/fxparam/1/value
	addr := fmt.Sprintf("/track/%s/fx/%s/fxparam/%d/value", trackID, pluginID, paramIndex)
	msg := osc.NewMessage(addr)
	msg.Append(value)
	return r.OSCClient.Send(msg)
}

// LuaBridgeCall facilitates calling custom ReaScripts via the file-based bridge if needed.
func (r *ReaperDriver) LuaBridgeCall(funcName string, args []interface{}) error {
	// This would implement the file writing logic found in total-reaper-mcp
	return nil
}
