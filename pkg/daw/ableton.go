package daw

import (
	"fmt"
	"net"

	"github.com/hypebeast/go-osc/osc"
)

type AbletonLiveDriver struct {
	OSCClient *osc.Client
	OSCHost   string
	OSCPort   int
	LocalPort int
}

func NewAbletonDriver(host string, port int, localPort int) *AbletonLiveDriver {
	return &AbletonLiveDriver{
		OSCClient: osc.NewClient(host, port),
		OSCHost:   host,
		OSCPort:   port,
		LocalPort: localPort,
	}
}

func (a *AbletonLiveDriver) Connect(endpoint string) error {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", a.OSCHost, a.OSCPort))
	if err != nil {
		return fmt.Errorf("failed to establish connection boundary to Ableton local socket listener: %w", err)
	}
	conn.Close()
	return nil
}

func (a *AbletonLiveDriver) Disconnect() error {
	return nil
}

func (a *AbletonLiveDriver) SetTransportState(playing bool, bpm float64) error {
	var cmd int32 = 0
	if playing {
		cmd = 1
	}

	// Unified Protocol: /superdaw/transport/play
	msg1 := osc.NewMessage("/superdaw/transport/play")
	msg1.Append(cmd)
	if err := a.OSCClient.Send(msg1); err != nil {
		return err
	}

	// Unified Protocol: /superdaw/transport/tempo
	msg2 := osc.NewMessage("/superdaw/transport/tempo")
	msg2.Append(float32(bpm))
	return a.OSCClient.Send(msg2)
}

func (a *AbletonLiveDriver) GetTransportState() (bool, float64, error) {
	return false, 120.0, nil
}

func (a *AbletonLiveDriver) GetTracks() ([]TrackConfig, error) {
	return []TrackConfig{}, nil
}

func (a *AbletonLiveDriver) CreateTrack(name string, trackType string) (string, error) {
	// Unified Protocol: /superdaw/track/create
	msg := osc.NewMessage("/superdaw/track/create")
	msg.Append(name)
	msg.Append(trackType)
	err := a.OSCClient.Send(msg)
	return "dynamic_generated_id", err
}

func (a *AbletonLiveDriver) SetTrackVolume(trackID string, volume float32) error {
	// Unified Protocol: /superdaw/track/volume
	msg := osc.NewMessage("/superdaw/track/volume")
	msg.Append(trackID)
	msg.Append(volume)
	return a.OSCClient.Send(msg)
}

func (a *AbletonLiveDriver) WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error {
	// Unified Protocol: /superdaw/clip/write
	msg := osc.NewMessage("/superdaw/clip/write")
	msg.Append(trackID)
	msg.Append(int32(clipIndex))
	for _, n := range notes {
		msg.Append(int32(n.Pitch))
		msg.Append(int32(n.Velocity))
		msg.Append(float32(n.StartBeat))
		msg.Append(float32(n.Duration))
	}
	return a.OSCClient.Send(msg)
}

func (a *AbletonLiveDriver) InstantiatePlugin(trackID string, pluginName string) (string, error) {
	msg := osc.NewMessage("/superdaw/device/instantiate")
	msg.Append(trackID)
	msg.Append(pluginName)
	err := a.OSCClient.Send(msg)
	return "vst_device_node", err
}

func (a *AbletonLiveDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	msg := osc.NewMessage("/superdaw/device/param")
	msg.Append(trackID)
	msg.Append(pluginID)
	msg.Append(int32(paramIndex))
	msg.Append(value)
	return a.OSCClient.Send(msg)
}
