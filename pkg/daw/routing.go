package daw

import (
	"fmt"
	"os"
	"github.com/hypebeast/go-osc/osc"
)

// AudioRouter handles virtual audio connections between DAWs via JACK or ReRoute.
type AudioRouter struct {
	OSCClient *osc.Client
}

func NewAudioRouter(host string, port int) *AudioRouter {
	return &AudioRouter{OSCClient: osc.NewClient(host, port)}
}

// Patch connects a source track from one DAW to a destination track in another.
func (r *AudioRouter) Patch(sourceDAW, sourceTrack, destDAW, destTrack string) error {
	m := osc.NewMessage("/superdaw/routing/patch")
	m.Append(sourceDAW)
	m.Append(sourceTrack)
	m.Append(destDAW)
	m.Append(destTrack)

	fmt.Fprintf(os.Stderr, "Audio Routing: %s (%s) -> %s (%s)\n", sourceDAW, sourceTrack, destDAW, destTrack)
	return r.OSCClient.Send(m)
}

// Unpatch removes an audio connection.
func (r *AudioRouter) Unpatch(sourceDAW, sourceTrack, destDAW, destTrack string) error {
	m := osc.NewMessage("/superdaw/routing/unpatch")
	m.Append(sourceDAW)
	m.Append(sourceTrack)
	m.Append(destDAW)
	m.Append(destTrack)
	return r.OSCClient.Send(m)
}
