package main

import (
	"fmt"
	"github.com/hypebeast/go-osc/osc"
	"time"
)

func main() {
	dispatcher := osc.NewStandardDispatcher()
	dispatcher.AddMsgHandler("/test", func(msg *osc.Message) {
		fmt.Println("Received:", msg.Address, msg.Arguments)
	})
	
	server := &osc.Server{Addr: "127.0.0.1:11001", Dispatcher: dispatcher}
	go server.ListenAndServe()
	
	fmt.Println("Server started on 11001")
	time.Sleep(10 * time.Second)
}
