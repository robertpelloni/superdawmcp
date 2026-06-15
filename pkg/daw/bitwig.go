package daw

import (
	"encoding/json"
	"fmt"
	"net"
)

type BitwigDriver struct {
	conn net.Conn
}

func NewBitwigDriver(host string, port int) *BitwigDriver {
	return &BitwigDriver{}
}

func (b *BitwigDriver) GetType() string { return "bitwig" }
func (b *BitwigDriver) Connect(endpoint string) error {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil { return err }
	b.conn = conn
	return nil
}

func (b *BitwigDriver) Disconnect() error {
	if b.conn != nil { return b.conn.Close() }
	return nil
}

func (b *BitwigDriver) send(method string, params map[string]interface{}) error {
	if b.conn == nil { return fmt.Errorf("not connected") }
	req := map[string]interface{}{"jsonrpc": "2.0", "method": method, "params": params, "id": 1}
	data, _ := json.Marshal(req)
	_, err := b.conn.Write(append(data, '\n'))
	return err
}

func (b *BitwigDriver) SetTransportState(playing bool, bpm float64) error {
	return b.send("transport.set_playing", map[string]interface{}{"playing": playing})
}

func (b *BitwigDriver) GetTransportState() (bool, float64, error) { return false, 120.0, nil }

func (b *BitwigDriver) GetTracks() ([]TrackConfig, error) { return []TrackConfig{}, nil }

func (b *BitwigDriver) CreateTrack(name, trackType string) (string, error) {
	err := b.send("track.create", map[string]interface{}{"name": name, "type": trackType})
	return "id", err
}

func (b *BitwigDriver) SetTrackVolume(id string, volume float32) error {
	return b.send("track.set_volume", map[string]interface{}{"index": id, "volume": volume})
}

func (b *BitwigDriver) SetTrackPan(id string, pan float32) error {
	return b.send("track.set_pan", map[string]interface{}{"index": id, "pan": pan})
}

func (b *BitwigDriver) SetTrackInstrument(id string, instrument string) error {
	return b.send("track.instrument", map[string]interface{}{"index": id, "instrument": instrument})
}

func (b *BitwigDriver) WriteMIDIClip(id string, idx int, notes []MIDINote) error {
	return b.send("clip.write_notes", map[string]interface{}{"track_id": id, "clip_index": idx, "notes": notes})
}

func (b *BitwigDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, nil }

func (b *BitwigDriver) DeleteClip(id string, idx int) error {
	return b.send("clip.delete", map[string]interface{}{"track_id": id, "clip_index": idx})
}

func (b *BitwigDriver) SetNotifyHandler(handler func(method string, params interface{})) {
}

func (b *BitwigDriver) ExecuteCustomCommand(cmd string, args map[string]interface{}) (interface{}, error) {
	err := b.send("custom."+cmd, args)
	return "Sent custom command to Bitwig", err
}

func (b *BitwigDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	return b.send("plugin.set_parameter", map[string]interface{}{"track_id": trackID, "plugin_id": pluginID, "index": paramIndex, "value": value})
}

func (b *BitwigDriver) SendCC(trackID string, controller int, value int) error {
	return b.send("midi.send_cc", map[string]interface{}{"track_id": trackID, "controller": controller, "value": value})
}
