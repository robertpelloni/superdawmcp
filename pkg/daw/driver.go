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

type DAWDriver interface {
	Connect(endpoint string) error
	Disconnect() error
	SetTransportState(playing bool, bpm float64) error
	CreateTrack(name string, trackType string) (string, error)
	SetTrackVolume(trackID string, volume float32) error
	SetTrackPan(trackID string, pan float32) error
	WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error
	ListClips(trackID string) ([]ClipInfo, error)
	DeleteClip(trackID string, clipIndex int) error
}
