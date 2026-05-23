package engine

import (
	"fmt"
	"time"
)

// LinkBridge provides a unified network clock for all connected DAWs.
// In a real implementation, this would use a CGO bridge to the Ableton Link C++ SDK.
type LinkBridge struct {
	BPM       float64
	IsPlaying bool
	StartTime time.Time
}

func NewLinkBridge() *LinkBridge {
	return &LinkBridge{
		BPM:       120.0,
		IsPlaying: false,
		StartTime: time.Now(),
	}
}

// GetBeat returns the current beat count since StartTime.
func (l *LinkBridge) GetBeat() float64 {
	if !l.IsPlaying { return 0 }
	elapsed := time.Since(l.StartTime).Seconds()
	return elapsed * (l.BPM / 60.0)
}

// Sync updates the Link session state.
func (l *LinkBridge) Sync(playing bool, bpm float64) {
	if playing && !l.IsPlaying {
		l.StartTime = time.Now()
	}
	l.IsPlaying = playing
	l.BPM = bpm
	fmt.Printf("Link Bridge: Sync state -> playing=%v, bpm=%.2f\n", playing, bpm)
}
