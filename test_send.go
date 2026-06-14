package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

func main() {
	client := osc.NewClient("127.0.0.1", 11000)
	
	// Send tempo
	client.Send(osc.NewMessage("/superdaw/transport/tempo", float32(148.0)))
	fmt.Println("Sent tempo")
	
	// Send track create
	client.Send(osc.NewMessage("/superdaw/track/create", "Psy Bass", "midi"))
	fmt.Println("Sent track create")
}
