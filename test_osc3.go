package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

func main() {
	m := osc.NewMessage("/superdaw/track/create", "Psy Bass", "midi")
	fmt.Printf("Message: %+v\n", m)
	fmt.Printf("Address: %s\n", m.Address)
	fmt.Printf("Arguments: %+v\n", m.Arguments)
}
