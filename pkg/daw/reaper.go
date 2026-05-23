package daw

import (
	"encoding/json"
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

type ReaperDriver struct {
	OSCClient *osc.Client
}

func NewReaperDriver(host string, port, webPort int) *ReaperDriver {
	return &ReaperDriver{OSCClient: osc.NewClient(host, port)}
}

func (r *ReaperDriver) Connect(endpoint string) error { return nil }
func (r *ReaperDriver) Disconnect() error { return nil }

func (r *ReaperDriver) SetTransportState(playing bool, bpm float64) error {
	m := osc.NewMessage("/superdaw/transport/play")
	v := int32(0); if playing { v = 1 }
	m.Append(v)
	return r.OSCClient.Send(m)
}

func (r *ReaperDriver) CreateTrack(name, trackType string) (string, error) {
	m := osc.NewMessage("/superdaw/track/create")
	m.Append(name)
	m.Append(trackType)
	return "id", r.OSCClient.Send(m)
}

func (r *ReaperDriver) SetTrackVolume(id string, volume float32) error {
	m := osc.NewMessage("/superdaw/track/volume")
	m.Append(id)
	m.Append(volume)
	return r.OSCClient.Send(m)
}

func (r *ReaperDriver) SetTrackPan(id string, pan float32) error {
	m := osc.NewMessage("/superdaw/track/pan")
	m.Append(id)
	m.Append(pan)
	return r.OSCClient.Send(m)
}

func (r *ReaperDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	m := osc.NewMessage("/superdaw/clip/write")
	m.Append(id)
	m.Append(int32(idx))
	p, _ := json.Marshal(notes)
	m.Append(string(p))
	return r.OSCClient.Send(m)
}

func (r *ReaperDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }

func (r *ReaperDriver) DeleteClip(id string, idx int) error {
	m := osc.NewMessage("/superdaw/clip/delete")
	m.Append(id)
	m.Append(int32(idx))
	return r.OSCClient.Send(m)
}

func (r *ReaperDriver) ExecuteCustomCommand(cmd string, args map[string]interface{}) (interface{}, error) {
	switch cmd {
	case "run_action":
		actionID, ok := args["action_id"].(string)
		if !ok { return nil, fmt.Errorf("missing action_id") }
		m := osc.NewMessage("/superdaw/action")
		m.Append(actionID)
		return "Triggered REAPER action", r.OSCClient.Send(m)
	}
	return nil, fmt.Errorf("unknown command: %s", cmd)
}
