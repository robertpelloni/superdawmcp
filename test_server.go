package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
	"time"
)

func main() {
	// Test if we can create an OSC server
	dispatcher := osc.NewStandardDispatcher()
	dispatcher.AddMsgHandler("/test", func(msg *osc.Message) {
		fmt.Println("Received:", msg.Address, msg.Arguments)
	})
	
	server := &osc.Server{Addr: "127.0.0.1:11000", Dispatcher: dispatcher}
	go server.ListenAndServe()
	
	fmt.Println("Server started on 11000")
	time.Sleep(2 * time.Second)
	
	// Send test
	client := osc.NewClient("127.0.0.1", 11000)
	client.Send(osc.NewMessage("/test", "hello"))
	
	time.Sleep(1 * time.Second)
}
