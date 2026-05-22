package integration
import (
	"fmt"
	"net"
	"sync"
	"github.com/hypebeast/go-osc/osc"
)
type MockDAW struct {
	messages []osc.Packet
	mu       sync.Mutex
	port     int
	stop     chan struct{}
}
func NewMockDAW(port int) *MockDAW {
	return &MockDAW{port: port, stop: make(chan struct{})}
}
func (m *MockDAW) Start() error {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", m.port))
	if err != nil { return err }
	conn, err := net.ListenUDP("udp", addr)
	if err != nil { return err }
	defer conn.Close()
	go func() { <-m.stop; conn.Close() }()
	for {
		buf := make([]byte, 65536)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			select {
			case <-m.stop: return nil
			default: return err
			}
		}
		packet, err := osc.ParsePacket(string(buf[:n]))
		if err != nil { continue }
		m.mu.Lock()
		m.messages = append(m.messages, packet)
		m.mu.Unlock()
	}
}
func (m *MockDAW) Stop() { close(m.stop) }
func (m *MockDAW) GetMessages() []osc.Message {
	m.mu.Lock(); defer m.mu.Unlock()
	var res []osc.Message
	for _, p := range m.messages {
		switch packet := p.(type) {
		case *osc.Message: res = append(res, *packet)
		case *osc.Bundle:
			for _, item := range packet.Messages { res = append(res, *item) }
		}
	}
	return res
}
