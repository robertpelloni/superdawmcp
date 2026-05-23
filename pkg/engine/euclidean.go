package engine
import "github.com/robertpelloni/superdaw-mcp/pkg/daw"
func GenerateEuclidean(h, s, p, v, r int, l float32) []daw.MIDINote {
	pat := make([]int, s); for i := 0; i < s; i++ { if (i*h)%s < h { pat[i] = 1 } }
	sz := l / float32(s); var res []daw.MIDINote
	for i, val := range pat { if val == 1 { idx := (i + r) % s; res = append(res, daw.MIDINote{Pitch: p, Velocity: v, StartBeat: float32(idx) * sz, Duration: sz}) } }
	return res
}
