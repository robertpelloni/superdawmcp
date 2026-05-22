package daw

import (
	"net"
	"testing"
	"time"

	"github.com/hypebeast/go-osc/osc"
)

func TestLatency(t *testing.T) {
	// Setup mock OSC server to measure loopback latency
	addr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	port := conn.LocalAddr().(*net.UDPAddr).Port
	client := osc.NewClient("127.0.0.1", port)

	start := time.Now()
	msg := osc.NewMessage("/ping")
	if err := client.Send(msg); err != nil {
		t.Fatal(err)
	}

	// Read from conn
	buf := make([]byte, 1024)
	conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
	_, _, err = conn.ReadFromUDP(buf)
	if err != nil {
		t.Fatal(err)
	}

	latency := time.Since(start)
	t.Logf("Loopback latency: %v", latency)

	if latency > 2*time.Millisecond {
		t.Errorf("Latency too high: %v (expected < 2ms)", latency)
	}
}
