package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

func main() {
	// Test what Go sends
	client := osc.NewClient("127.0.0.1", 11000)
	
	// Test 1: tempo
	msg1 := osc.NewMessage("/superdaw/transport/tempo", float32(148.0))
	fmt.Println("Tempo message:", msg1.String())
	client.Send(msg1)
	
	// Test 2: track create
	msg2 := osc.NewMessage("/superdaw/track/create", "Psy Bass", "midi")
	fmt.Println("Track create message:", msg2.String())
	client.Send(msg2)
	
	// Test 3: clip write
	msg3 := osc.NewMessage("/superdaw/clip/write", "0", int32(0), `[{"pitch":60,"start_beat":0,"duration":0.25,"velocity":100}]`)
	fmt.Println("Clip write message:", msg3.String())
	client.Send(msg3)
}
