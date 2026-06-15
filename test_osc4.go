package main

import (
	"fmt"
	"encoding/json"
	"github.com/hypebeast/go-osc/osc"
)

type MIDINote struct {
	Pitch     int     `json:"pitch"`
	StartBeat float64 `json:"start_beat"`
	Duration  float64 `json:"duration"`
	Velocity  int     `json:"velocity"`
}

func main() {
	notes := []MIDINote{
		{Pitch: 36, StartBeat: 0.0, Duration: 1.0, Velocity: 100},
		{Pitch: 36, StartBeat: 1.0, Duration: 1.0, Velocity: 100},
	}
	data, _ := json.Marshal(notes)
	
	m := osc.NewMessage("/superdaw/clip/write", "0", int32(0), string(data))
	fmt.Printf("Message: %+v\n", m)
	fmt.Printf("Address: %s\n", m.Address)
	fmt.Printf("Arguments: %+v\n", m.Arguments)
}
