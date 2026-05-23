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

func NewReaperDriver(h string, p, w int) *ReaperDriver {
	return &ReaperDriver{
		OSCClient: osc.NewClient(h, p),
		WebPort:   w,
		Host:      h,
	}
}

func (r *ReaperDriver) Connect(e string) error { return nil }
func (r *ReaperDriver) Disconnect() error { return nil }

func (r *ReaperDriver) SetTransportState(p bool, bpm float64) error {
	msg := osc.NewMessage("/superdaw/transport/play")
	val := int32(0); if p { val = 1 }
	msg.Append(val)
	if err := r.OSCClient.Send(msg); err != nil { return err }

	// Use Web API for tempo
	url := fmt.Sprintf("http://%s:%d/wwr/_SET/BPM/%.2f", r.Host, r.WebPort, bpm)
	resp, err := http.Get(url)
	if err == nil { resp.Body.Close() }
	return err
}

func (r *ReaperDriver) CreateTrack(n, t string) (string, error) {
	msg := osc.NewMessage("/superdaw/track/create")
	msg.Append(n); msg.Append(t)
	return "id", r.OSCClient.Send(msg)
}

func (r *ReaperDriver) SetTrackVolume(id string, v float32) error {
	msg := osc.NewMessage("/superdaw/track/volume")
	msg.Append(id); msg.Append(v)
	return r.OSCClient.Send(msg)
}

func (r *ReaperDriver) SetTrackPan(id string, p float32) error {
	msg := osc.NewMessage("/superdaw/track/pan")
	msg.Append(id); msg.Append(p)
	return r.OSCClient.Send(msg)
}

func (r *ReaperDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	msg := osc.NewMessage("/superdaw/clip/write")
	msg.Append(id); msg.Append(int32(idx))
	payload, _ := json.Marshal(notes)
	msg.Append(string(payload))
	return r.OSCClient.Send(msg)
}

func (r *ReaperDriver) ListClips(id string) ([]ClipInfo, error) {
	return []ClipInfo{}, nil
}

func (r *ReaperDriver) DeleteClip(id string, idx int) error {
	msg := osc.NewMessage("/superdaw/clip/delete")
	msg.Append(id); msg.Append(int32(idx))
	return r.OSCClient.Send(msg)
}
