package engine

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"
)

// LinkBridge provides a unified network clock and tempo arbiter for all connected DAWs.
// Supports beat time calculation, quantization quantum alignment (e.g. 4 beats = 1 bar),
// peer session count monitoring, and continuous tempo ramping.
type LinkBridge struct {
	mu          sync.RWMutex
	BPM         float64
	IsPlaying   bool
	StartTime   time.Time
	Peers       int
	Quantum     float64 // Quantum in beats (default 4.0 for standard 4/4 bar)
}

// NewLinkBridge instantiates a new Link session bridge with default 120.0 BPM and 4/4 quantum.
func NewLinkBridge() *LinkBridge {
	return &LinkBridge{
		BPM:       120.0,
		IsPlaying: false,
		StartTime: time.Now(),
		Peers:     1,
		Quantum:   4.0,
	}
}

// GetBeat returns the fractional beat count since transport start.
func (l *LinkBridge) GetBeat() float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if !l.IsPlaying {
		return 0.0
	}
	elapsed := time.Since(l.StartTime).Seconds()
	return elapsed * (l.BPM / 60.0)
}

// GetPhase returns the current phase within the quantum (0.0 <= phase < Quantum).
func (l *LinkBridge) GetPhase() float64 {
	beat := l.GetBeat()
	l.mu.RLock()
	quantum := l.Quantum
	l.mu.RUnlock()

	if quantum <= 0.0 {
		quantum = 4.0
	}
	return math.Mod(beat, quantum)
}

// GetQuantizedBeat returns the nearest upcoming beat quantized to quantum boundaries.
func (l *LinkBridge) GetQuantizedBeat(quantum float64) float64 {
	if quantum <= 0.0 {
		quantum = 4.0
	}
	beat := l.GetBeat()
	return math.Ceil(beat/quantum) * quantum
}

// Sync updates the Link session state thread-safely.
func (l *LinkBridge) Sync(playing bool, bpm float64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if playing && !l.IsPlaying {
		l.StartTime = time.Now()
	}
	l.IsPlaying = playing
	if bpm >= 20.0 && bpm <= 999.0 {
		l.BPM = bpm
	}
	fmt.Fprintf(os.Stderr, "Link Bridge: Sync state -> playing=%v, bpm=%.2f, peers=%d\n", playing, l.BPM, l.Peers)
}

// SetPeers sets the number of connected Link-compatible peers on the local network.
func (l *LinkBridge) SetPeers(peers int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if peers < 0 {
		peers = 0
	}
	l.Peers = peers
}

// SetTempo clamps and assigns the global Link BPM.
func (l *LinkBridge) SetTempo(bpm float64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if bpm < 20.0 {
		bpm = 20.0
	} else if bpm > 999.0 {
		bpm = 999.0
	}
	l.BPM = bpm
}

// GetTempo returns current BPM.
func (l *LinkBridge) GetTempo() float64 {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.BPM
}

// RampTempo smoothly transitions BPM from current to target over the given duration.
func (l *LinkBridge) RampTempo(targetBPM float64, duration time.Duration) {
	if duration <= 0 {
		l.SetTempo(targetBPM)
		return
	}
	startBPM := l.GetTempo()
	steps := 10
	stepDuration := duration / time.Duration(steps)

	go func() {
		for i := 1; i <= steps; i++ {
			time.Sleep(stepDuration)
			alpha := float64(i) / float64(steps)
			current := startBPM + (targetBPM-startBPM)*alpha
			l.SetTempo(current)
		}
	}()
}

// GetStatus returns the complete Link session status report.
func (l *LinkBridge) GetStatus() map[string]interface{} {
	l.mu.RLock()
	defer l.mu.RUnlock()

	beat := 0.0
	phase := 0.0
	if l.IsPlaying {
		elapsed := time.Since(l.StartTime).Seconds()
		beat = elapsed * (l.BPM / 60.0)
		phase = math.Mod(beat, l.Quantum)
	}

	return map[string]interface{}{
		"playing": l.IsPlaying,
		"bpm":     l.BPM,
		"peers":   l.Peers,
		"quantum": l.Quantum,
		"beat":    math.Round(beat*100) / 100,
		"phase":   math.Round(phase*100) / 100,
	}
}
