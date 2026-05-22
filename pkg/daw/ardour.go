package daw

import (
	"fmt"

	"github.com/hypebeast/go-osc/osc"
)

type ArdourDriver struct {
	OSCClient *osc.Client
	Host      string
	Port      int
}

func NewArdourDriver(host string, port int) *ArdourDriver {
	return &ArdourDriver{
		OSCClient: osc.NewClient(host, port),
		Host:      host,
		Port:      port,
	}
}

func (a *ArdourDriver) Connect(endpoint string) error {
	return nil
}

func (a *ArdourDriver) Disconnect() error {
	return nil
}

func (a *ArdourDriver) SetTransportState(playing bool, bpm float64) error {
	var addr string
	if playing {
		addr = "/transport_play"
	} else {
		addr = "/transport_stop"
	}
	msg := osc.NewMessage(addr)
	if err := a.OSCClient.Send(msg); err != nil {
		return err
	}

	tempoMsg := osc.NewMessage("/set_transport_speed")
	tempoMsg.Append(float32(1.0)) // speed
	return a.OSCClient.Send(tempoMsg)
}

func (a *ArdourDriver) GetTransportState() (bool, float64, error) {
	return false, 120.0, nil
}

func (a *ArdourDriver) GetTracks() ([]TrackConfig, error) {
	return []TrackConfig{}, nil
}

func (a *ArdourDriver) CreateTrack(name string, trackType string) (string, error) {
	msg := osc.NewMessage("/access_action")
	msg.Append("Track/add-audio-track")
	err := a.OSCClient.Send(msg)
	return "ardour_track_new", err
}

func (a *ArdourDriver) SetTrackVolume(trackID string, volume float32) error {
	// Ardour 5+ uses /strip/fader
	addr := fmt.Sprintf("/strip/fader")
	msg := osc.NewMessage(addr)
	// trackID in Ardour is often an integer index
	msg.Append(trackID)
	msg.Append(volume)
	return a.OSCClient.Send(msg)
}

func (a *ArdourDriver) WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error {
	return fmt.Errorf("MIDI writing not natively supported via Ardour OSC")
}

func (a *ArdourDriver) InstantiatePlugin(trackID string, pluginName string) (string, error) {
	return "", fmt.Errorf("plugin instantiation not supported via Ardour OSC")
}

func (a *ArdourDriver) SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error {
	addr := "/strip/plugin/parameter"
	msg := osc.NewMessage(addr)
	msg.Append(trackID)
	msg.Append(pluginID)
	msg.Append(int32(paramIndex))
	msg.Append(value)
	return a.OSCClient.Send(msg)
}
