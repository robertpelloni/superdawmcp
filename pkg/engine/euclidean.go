package engine
import "github.com/robertpelloni/superdaw-mcp/pkg/daw"
func GenerateEuclidean(hits, steps, pitch, velocity, rotate int, length float32) []daw.MIDINote {
	pattern := make([]int, steps)
	for i := 0; i < steps; i++ { if (i*hits)%steps < hits { pattern[i] = 1 } }
	stepSize := length / float32(steps)
	var notes []daw.MIDINote
	for i, val := range pattern {
		if val == 1 {
			idx := (i + rotate) % steps
			notes = append(notes, daw.MIDINote{Pitch: pitch, Velocity: velocity, StartBeat: float32(idx) * stepSize, Duration: stepSize})
		}
	}
	return notes
}
