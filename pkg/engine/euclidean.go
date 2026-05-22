package engine

import (
	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
)

// GenerateEuclidean creates a MIDI note pattern using the Euclidean algorithm.
func GenerateEuclidean(hits, steps, pitch, velocity int, rotate int, lengthBeats float32) []daw.MIDINote {
	if steps <= 0 || hits <= 0 {
		return nil
	}
	if hits > steps {
		hits = steps
	}

	// Bresenham-like Euclidean generation
	pattern := make([]int, steps)
	for i := 0; i < steps; i++ {
		if (i*hits)%steps < hits {
			pattern[i] = 1
		}
	}

	// Apply rotation
	finalPattern := make([]int, steps)
	for i := 0; i < steps; i++ {
		idx := (i + rotate) % steps
		if idx < 0 {
			idx += steps
		}
		finalPattern[idx] = pattern[i]
	}

	stepSize := lengthBeats / float32(steps)
	var notes []daw.MIDINote

	for i, val := range finalPattern {
		if val == 1 {
			notes = append(notes, daw.MIDINote{
				Pitch:     pitch,
				Velocity:  velocity,
				StartBeat: float32(i) * stepSize,
				Duration:  stepSize,
			})
		}
	}

	return notes
}
