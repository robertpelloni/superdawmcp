package engine

import (
	"math/rand"
	"github.com/robertpelloni/superdaw-mcp/pkg/daw"
)

// GenerateMusic generates a generative MIDI sequence based on a style prompt.
func GenerateMusic(style string, bars int) []daw.MIDINote {
	notes := []daw.MIDINote{}

	// Root C (60), Scale: Minor for techno/ambient, Major for basic
	root := 60
	scaleName := "minor"
	if style == "basic" || style == "pop" { scaleName = "major" }

	if style == "techno" || style == "ambient" {
		prog := GenerateProgression(style, root, scaleName)
		for b := 0; b < bars; b++ {
			chord := prog[b%len(prog)]
			for i, pitch := range chord {
				if style == "techno" {
					// 16th note pulses
					for pulse := 0; i == 0 && pulse < 16; pulse++ {
						if rand.Float32() > 0.7 {
							notes = append(notes, daw.MIDINote{
								Pitch:     pitch,
								Velocity:  90 + rand.Intn(30),
								StartBeat: float32(b)*4.0 + float32(pulse)*0.25,
								Duration:  0.2,
							})
						}
					}
				} else {
					// Ambient pads
					notes = append(notes, daw.MIDINote{
						Pitch:     pitch + 12, // Octave up
						Velocity:  60,
						StartBeat: float32(b) * 4.0,
						Duration:  3.8,
					})
				}
			}
		}
	} else if style == "ambient" {
		// Longer pads
		scale := GetScaleNotes(root, scaleName)
		for i := 0; i < bars; i++ {
			notes = append(notes, daw.MIDINote{
				Pitch:     root + scale[rand.Intn(len(scale))] + 12,
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
