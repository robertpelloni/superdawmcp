package daw

type TrackConfig struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Volume float32 `json:"volume"`
	Pan    float32 `json:"pan"`
}

type MIDINote struct {
	Pitch     int     `json:"pitch"`
	Velocity  int     `json:"velocity"`
	StartBeat float32 `json:"start_beat"`
	Duration  float32 `json:"duration"`
}

type ClipInfo struct {
	Index int    `json:"index"`
	Name  string `json:"name"`
}

// DAWDriver defines the capability interface for all integrated audio engines.
type DAWDriver interface {
	GetType() string
	Connect(endpoint string) error
	Disconnect() error

	// Transport Actions
	SetTransportState(playing bool, bpm float64) error
	GetTransportState() (bool, float64, error)

	// Mixer Actions
	GetTracks() ([]TrackConfig, error)
	CreateTrack(name, trackType string) (string, error)
	SetTrackVolume(trackID string, volume float32) error
	SetTrackPan(trackID string, pan float32) error
	SetTrackInstrument(trackID string, instrument string) error

	// Clip & MIDI
	WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error
	ListClips(trackID string) ([]ClipInfo, error)
	DeleteClip(trackID string, clipIndex int) error

	// Plugins
	SetPluginParameter(trackID string, pluginID string, paramIndex int, value float32) error

	// MIDI
	SendCC(trackID string, controller int, value int) error

	// Custom DAW-specific extensions
	ExecuteCustomCommand(command string, args map[string]interface{}) (interface{}, error)

	// Notifications
	SetNotifyHandler(handler func(method string, params interface{}))
}
