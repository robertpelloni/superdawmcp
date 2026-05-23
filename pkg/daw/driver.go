package daw
type TrackConfig struct { ID, Name string; Volume, Pan float32 }
type MIDINote struct { Pitch, Velocity int; StartBeat, Duration float32 }
type ClipInfo struct { Index int; Name string }
type DAWDriver interface {
	Connect(e string) error; Disconnect() error
	SetTransportState(p bool, b float64) error
	CreateTrack(n, t string) (string, error)
	SetTrackVolume(id string, v float32) error
	SetTrackPan(id string, p float32) error
	WriteMIDIClip(id string, idx int, notes []MIDINote) error
	ListClips(id string) ([]ClipInfo, error)
	DeleteClip(id string, idx int) error
}
