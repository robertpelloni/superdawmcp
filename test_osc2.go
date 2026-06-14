package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
)

func main() {
	m := osc.NewMessage("/superdaw/transport/tempo", float32(148.0))
	fmt.Printf("Message: %+v\n", m)
	fmt.Printf("Address: %s\n", m.Address)
	fmt.Printf("Arguments: %+v\n", m.Arguments)
}
