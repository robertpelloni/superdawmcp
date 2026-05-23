package engine

import (
	"math/rand"
	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
)

// GenerateMusic generates a generative MIDI sequence based on a style prompt.
func GenerateMusic(style string, bars int) []daw.MIDINote {
	notes := []daw.MIDINote{}

	// Example Logic: C Major Pentatonic (0, 2, 4, 7, 9)
	scale := []int{60, 62, 64, 67, 69}

	if style == "techno" {
		// 16th note pulses
		for i := 0; i < bars*16; i++ {
			if rand.Float32() > 0.7 {
				notes = append(notes, daw.MIDINote{
					Pitch:     scale[rand.Intn(len(scale))],
					Velocity:  90 + rand.Intn(30),
					StartBeat: float32(i) * 0.25,
					Duration:  0.2,
				})
			}
		}
	} else if style == "ambient" {
		// Longer pads
		for i := 0; i < bars; i++ {
			notes = append(notes, daw.MIDINote{
				Pitch:     scale[rand.Intn(len(scale))],
				Velocity:  60,
				StartBeat: float32(i) * 4.0,
				Duration:  3.8,
			})
		}
	} else {
		// Basic 4/4 floor
		for i := 0; i < bars*4; i++ {
			notes = append(notes, daw.MIDINote{
				Pitch:     60,
				Velocity:  100,
				StartBeat: float32(i),
				Duration:  0.1,
			})
		}
	}

	return notes
}
