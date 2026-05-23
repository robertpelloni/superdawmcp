package daw

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"
	"sync"
	"github.com/hypebeast/go-osc/osc"
)

type ReaperDriver struct {
	OSCClient *osc.Client
	bridgeDir string
	requestID int
	state     struct {
		playing bool
		tempo   float64
		mu      sync.RWMutex
	}
}

func NewReaperDriver(host string, port, webPort int) *ReaperDriver {
	home, _ := os.UserHomeDir()
	bridgeDir := filepath.Join(home, "Library/Application Support/REAPER/Scripts/mcp_bridge_data")
	if runtime.GOOS == "windows" {
		bridgeDir = filepath.Join(os.Getenv("APPDATA"), "REAPER/Scripts/mcp_bridge_data")
	}

	return &ReaperDriver{
		OSCClient: osc.NewClient(host, port),
		bridgeDir: bridgeDir,
		requestID: 1,
	}
}

func (r *ReaperDriver) Connect(endpoint string) error {
	os.MkdirAll(r.bridgeDir, 0755)
	return nil
}

func (r *ReaperDriver) Disconnect() error { return nil }

func (r *ReaperDriver) callBridge(funcName string, args []interface{}) (map[string]interface{}, error) {
	id := r.requestID
	r.requestID++

	reqFile := filepath.Join(r.bridgeDir, fmt.Sprintf("request_%d.json", id))
	resFile := filepath.Join(r.bridgeDir, fmt.Sprintf("response_%d.json", id))

	reqData := map[string]interface{}{
		"func": funcName,
		"args": args,
	}
	data, _ := json.Marshal(reqData)
	err := os.WriteFile(reqFile, data, 0644)
	if err != nil { return nil, err }

	// Poll for response
	for i := 0; i < 50; i++ {
		time.Sleep(100 * time.Millisecond)
		if _, err := os.Stat(resFile); err == nil {
			resData, _ := os.ReadFile(resFile)
			var result map[string]interface{}
			json.Unmarshal(resData, &result)
			os.Remove(resFile)
			return result, nil
		}
	}
	return nil, fmt.Errorf("bridge timeout")
}

func (r *ReaperDriver) SetTransportState(playing bool, bpm float64) error {
	m := osc.NewMessage("/superdaw/transport/play")
	v := int32(0); if playing { v = 1 }
	m.Append(v)
	r.OSCClient.Send(m)

	_, err := r.callBridge("SetTempo", []interface{}{bpm})
	if err == nil {
		r.state.mu.Lock()
		r.state.playing = playing
		r.state.tempo = bpm
		r.state.mu.Unlock()
	}
	return err
}

func (r *ReaperDriver) GetTransportState() (bool, float64, error) {
	r.state.mu.RLock()
	defer r.state.mu.RUnlock()
	return r.state.playing, r.state.tempo, nil
}

func (r *ReaperDriver) CreateTrack(name, trackType string) (string, error) {
	_, err := r.callBridge("InsertTrackAtIndex", []interface{}{-1, true})
	return "id", err
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
	r.callBridge("CreateMIDIItem", []interface{}{0, 0, 4.0})
	for _, n := range notes {
		r.callBridge("InsertMIDINote", []interface{}{0, 0, n.Pitch, n.StartBeat, 0.5, n.Velocity, 1})
	}
	return nil
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
