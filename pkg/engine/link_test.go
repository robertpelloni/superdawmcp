package engine

import (
	"math"
	"testing"
	"time"
)

func TestLinkBridge_LifecycleAndState(t *testing.T) {
	bridge := NewLinkBridge()
	if bridge == nil {
		t.Fatal("expected NewLinkBridge to return non-nil instance")
	}

	status := bridge.GetStatus()
	if status["bpm"].(float64) != 120.0 {
		t.Errorf("expected initial bpm 120.0, got %v", status["bpm"])
	}
	if status["peers"].(int) != 1 {
		t.Errorf("expected initial peers 1, got %v", status["peers"])
	}

	bridge.SetPeers(3)
	if bridge.Peers != 3 {
		t.Errorf("expected peers to be 3, got %d", bridge.Peers)
	}

	bridge.SetPeers(-1)
	if bridge.Peers != 0 {
		t.Errorf("expected peers clamped to 0, got %d", bridge.Peers)
	}
}

func TestLinkBridge_TempoBoundsAndRamping(t *testing.T) {
	bridge := NewLinkBridge()

	// Test boundary clamps
	bridge.SetTempo(10.0) // below min (20.0)
	if bridge.GetTempo() != 20.0 {
		t.Errorf("expected tempo clamp to 20.0, got %f", bridge.GetTempo())
	}

	bridge.SetTempo(1200.0) // above max (999.0)
	if bridge.GetTempo() != 999.0 {
		t.Errorf("expected tempo clamp to 999.0, got %f", bridge.GetTempo())
	}

	bridge.SetTempo(128.0)
	if bridge.GetTempo() != 128.0 {
		t.Errorf("expected tempo to be 128.0, got %f", bridge.GetTempo())
	}

	// Test smooth tempo ramping
	rampDuration := 50 * time.Millisecond
	bridge.RampTempo(140.0, rampDuration)

	// Give goroutine time to complete
	time.Sleep(75 * time.Millisecond)

	finalTempo := bridge.GetTempo()
	if math.Abs(finalTempo-140.0) > 0.5 {
		t.Errorf("expected ramped tempo to reach ~140.0, got %f", finalTempo)
	}
}

func TestLinkBridge_TransportAndBeatCalculation(t *testing.T) {
	bridge := NewLinkBridge()
	bridge.SetTempo(120.0)

	// At 120 BPM: 2 beats per second = 1 beat every 500ms
	bridge.Sync(true, 120.0)
	if !bridge.IsPlaying {
		t.Errorf("expected transport to be playing")
	}

	time.Sleep(100 * time.Millisecond)
	beat := bridge.GetBeat()
	if beat <= 0.0 {
		t.Errorf("expected beat to advance past 0, got %f", beat)
	}

	phase := bridge.GetPhase()
	if phase < 0.0 || phase >= 4.0 {
		t.Errorf("expected phase to be in [0, 4), got %f", phase)
	}

	quantized := bridge.GetQuantizedBeat(4.0)
	if quantized < beat {
		t.Errorf("expected quantized beat >= current beat, got %f < %f", quantized, beat)
	}

	bridge.Sync(false, 120.0)
	if bridge.IsPlaying {
		t.Errorf("expected transport to be stopped")
	}

	if bridge.GetBeat() != 0.0 {
		t.Errorf("expected beat to be 0 when stopped, got %f", bridge.GetBeat())
	}
}
