package daw

type TrackConfig struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Volume float32 `json:"volume"` // Range 0.0 - 1.0
	Pan    float32 `json:"pan"`    // Range -1.0 to 1.0
	Solo   bool    `json:"solo"`
	Mute   bool    `json:"mute"`
}

type MIDINote struct {
	Pitch     int     `json:"pitch"`    // MIDI Note 0-127
	Velocity  int     `json:"velocity"` // 0-127
	StartBeat float32 `json:"start_beat"`
	Duration  float32 `json:"duration"`
}

// DAWDriver defines the capability interface for all integrated audio engines.
type DAWDriver interface {
	Connect(endpoint string) error
	Disconnect() error

	// Transport Actions
	SetTransportState(playing bool, bpm float64) error
	GetTransportState() (bool, float64, error)

	// Mixer Actions
	GetTracks() ([]TrackConfig, error)
	CreateTrack(name string, trackType string) (string, error)
	SetTrackVolume(trackID string, volume float32) error

	// Clip & Generation Engine
	WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error

	// Plugin / VST Interface
	InstantiatePlugin(trackID string, pluginName string) (string, error)
	SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error
}
