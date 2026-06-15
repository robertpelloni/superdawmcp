package integration
import ( "fmt"; "net"; "sync"; "github.com/hypebeast/go-osc/osc" )
type MockDAW struct { messages []osc.Packet; mu sync.Mutex; port int; stop chan struct{} }
func NewMockDAW(p int) *MockDAW { return &MockDAW{port: p, stop: make(chan struct{})} }
func (m *MockDAW) Start() error {
	a, err := net.ResolveUDPAddr("udp", fmt.Sprintf("127.0.0.1:%d", m.port)); if err != nil { return err }
	c, err := net.ListenUDP("udp", a); if err != nil { return err }; defer c.Close(); go func() { <-m.stop; c.Close() }()
	for {
		b := make([]byte, 65536); n, _, err := c.ReadFromUDP(b); if err != nil { return nil }
		p, _ := osc.ParsePacket(string(b[:n])); m.mu.Lock(); m.messages = append(m.messages, p); m.mu.Unlock()
	}
}
func (m *MockDAW) Stop() { close(m.stop) }
func (m *MockDAW) GetMessages() []osc.Message {
	m.mu.Lock(); defer m.mu.Unlock(); var res []osc.Message
	for _, p := range m.messages { switch packet := p.(type) { case *osc.Message: res = append(res, *packet); case *osc.Bundle: for _, item := range packet.Messages { res = append(res, *item) } } }
	return res
}
func (m *MockDAW) SendMessage(host string, port int, addr string, args ...interface{}) error {
	client := osc.NewClient(host, port)
	msg := osc.NewMessage(addr)
	for _, a := range args { msg.Append(a) }
	return client.Send(msg)
}
