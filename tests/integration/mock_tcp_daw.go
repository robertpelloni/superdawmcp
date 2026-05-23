package integration
import ( "bufio"; "fmt"; "net"; "sync" )
type MockTCPDAW struct { messages []string; mu sync.Mutex; port int; stop chan struct{} }
func NewMockTCPDAW(p int) *MockTCPDAW { return &MockTCPDAW{port: p, stop: make(chan struct{})} }
func (m *MockTCPDAW) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", m.port)); if err != nil { return err }; defer l.Close()
	go func() { <-m.stop; l.Close() }()
	for {
		c, err := l.Accept(); if err != nil { return nil }
		go func(conn net.Conn) {
			defer conn.Close(); r := bufio.NewReader(conn)
			for {
				line, err := r.ReadString('\n'); if err != nil { return }
				m.mu.Lock(); m.messages = append(m.messages, line); m.mu.Unlock()
				fmt.Fprintf(conn, "{\"jsonrpc\":\"2.0\",\"result\":\"Success\",\"id\":1}\n")
			}
		}(c)
	}
}
func (m *MockTCPDAW) Stop() { close(m.stop) }
func (m *MockTCPDAW) GetMessages() []string {
	m.mu.Lock(); defer m.mu.Unlock(); return append([]string{}, m.messages...)
}
