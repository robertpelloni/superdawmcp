package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
	"time"
)

func main() {
	// Test if we can create an OSC server
	dispatcher := osc.NewStandardDispatcher()
	dispatcher.AddMsgHandler("/superdaw/transport/tempo", func(msg *osc.Message) {
		fmt.Println("Received tempo:", msg.Arguments)
	})
	dispatcher.AddMsgHandler("/superdaw/track/create", func(msg *osc.Message) {
		fmt.Println("Received track create:", msg.Arguments)
	})
	
	server := &osc.Server{Addr: "127.0.0.1:11000", Dispatcher: dispatcher}
	go server.ListenAndServe()
	
	fmt.Println("Server started on 11000")
	
	// Send test
	client := osc.NewClient("127.0.0.1", 11000)
	client.Send(osc.NewMessage("/superdaw/transport/tempo", float32(148.0)))
	client.Send(osc.NewMessage("/superdaw/track/create", "Psy Bass", "midi"))
	
	time.Sleep(2 * time.Second)
	fmt.Println("Done")
}
