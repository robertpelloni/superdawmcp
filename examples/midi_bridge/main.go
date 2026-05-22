package main

import (
	"fmt"
	"log"
	"os"

	"github.com/robertpelloni/superdaw-mcp/pkg/client/go"
)

func main() {
	// A mock MIDI bridge that demonstrates how to map MIDI CC messages to SuperDAW tool calls.
	// In a real implementation, we would use a library like 'gitlab.com/gomidi/midi/v4'.

	if len(os.Args) < 2 {
		log.Fatalf("Usage: midi_bridge <path_to_superdaw_mcp>")
	}

	serverPath := os.Args[1]
	client, err := client.NewClient(serverPath)
	if err != nil {
		log.Fatalf("Failed to connect to SuperDAW: %v", err)
	}
	defer client.Close()

	log.Println("SuperDAW MIDI Bridge Started (Mocking CC 7 -> Volume)")

	// Simulate receiving a MIDI CC 7 (Volume) message
	trackID := "1"
	ccValue := 100 // MIDI range 0-127
	volume := float32(ccValue) / 127.0

	fmt.Printf("MIDI CC 7 received: mapping to volume %.2f for track %s\n", volume, trackID)

	err = client.SetMixer(trackID, volume, 0.0)
	if err != nil {
		log.Printf("Failed to set volume: %v", err)
	}

	log.Println("Bridge shutting down.")
}
