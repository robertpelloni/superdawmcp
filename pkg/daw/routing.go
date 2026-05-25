package daw

import (
	"fmt"
	"os"
	"github.com/hypebeast/go-osc/osc"
)

// AudioRouter handles virtual audio connections between DAWs via JACK or ReRoute.
type AudioRouter struct {
	OSCClient *osc.Client
	Jack      *JackRouter
}

func NewAudioRouter(host string, port int) *AudioRouter {
	return &AudioRouter{
		OSCClient: osc.NewClient(host, port),
		Jack:      NewJackRouter(),
	}
}

// Patch connects a source track from one DAW to a destination track in another.
func (r *AudioRouter) Patch(sourceDAW, sourceTrack, destDAW, destTrack string) error {
	// 1. Send OSC notification for bridge agents
	m := osc.NewMessage("/superdaw/routing/patch")
	m.Append(sourceDAW)
	m.Append(sourceTrack)
	m.Append(destDAW)
	m.Append(destTrack)
	r.OSCClient.Send(m)

	// 2. Execute JACK routing if applicable
	// Heuristic: sourceTrack/destTrack as "out1", "in1" etc.
	srcPort := fmt.Sprintf("%s:%s", sourceDAW, sourceTrack)
	dstPort := fmt.Sprintf("%s:%s", destDAW, destTrack)
	err := r.Jack.Patch(srcPort, dstPort)

	fmt.Fprintf(os.Stderr, "Audio Routing (JACK): %s -> %s\n", srcPort, dstPort)
	return err
}

// Unpatch removes an audio connection.
func (r *AudioRouter) Unpatch(sourceDAW, sourceTrack, destDAW, destTrack string) error {
	// 1. OSC Notification
	m := osc.NewMessage("/superdaw/routing/unpatch")
	m.Append(sourceDAW)
	m.Append(sourceTrack)
	m.Append(destDAW)
	m.Append(destTrack)
	r.OSCClient.Send(m)

	// 2. JACK Unpatch
	srcPort := fmt.Sprintf("%s:%s", sourceDAW, sourceTrack)
	dstPort := fmt.Sprintf("%s:%s", destDAW, destTrack)
	return r.Jack.Unpatch(srcPort, dstPort)
}
