package vst

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// GpuInspectorStream represents a high-frequency WebSocket/TCP channel
// intended to feed a WebGL-based frontend for real-time waveform and spectrum analysis.
type GpuInspectorStream struct {
	Active     bool
	PluginID   string
	SampleRate int
	lock       sync.RWMutex
}

var GlobalGpuStream = &GpuInspectorStream{
	Active:     false,
	SampleRate: 44100,
}

// StartStreaming initiates the real-time high-throughput telemetry required
// for a GPU-based FFT analyzer on the client side.
func (g *GpuInspectorStream) StartStreaming(pluginID string) error {
	g.lock.Lock()
	if g.Active {
		g.lock.Unlock()
		return fmt.Errorf("GPU inspector stream already active for %s", g.PluginID)
	}
	g.PluginID = pluginID
	g.Active = true
	g.lock.Unlock()

	// Background loop emitting dummy spectral data
	go func() {
		ticker := time.NewTicker(33 * time.Millisecond) // ~30 FPS
		defer ticker.Stop()
		for {
			g.lock.RLock()
			active := g.Active
			g.lock.RUnlock()
			if !active {
				break
			}
			<-ticker.C

			// In a real implementation, this would pull from the VST output buffer.
			// Here we generate simulated magnitude bins.
			_ = generateSimulatedFFT()
			// (Broadcast over WebSocket to WebGL clients)
		}
	}()

	return nil
}

// StopStreaming halts the high-frequency telemetry.
func (g *GpuInspectorStream) StopStreaming() {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.Active = false
	g.PluginID = ""
}

func generateSimulatedFFT() []float32 {
	bins := make([]float32, 256)
	for i := 0; i < len(bins); i++ {
		bins[i] = rand.Float32()
	}
	return bins
}
