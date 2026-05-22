package integration

import (
	"bufio"
	"fmt"
	"net"
	"sync"
)

type MockTCPDAW struct {
	messages []string
	mu       sync.Mutex
	port     int
	stop     chan struct{}
}

func NewMockTCPDAW(port int) *MockTCPDAW {
	return &MockTCPDAW{port: port, stop: make(chan struct{})}
}

func (m *MockTCPDAW) Start() error {
	l, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", m.port))
	if err != nil {
		return err
	}
	defer l.Close()

	go func() {
		<-m.stop
		l.Close()
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-m.stop:
				return nil
			default:
				return err
			}
		}
		go m.handleConn(conn)
	}
}

func (m *MockTCPDAW) handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		m.mu.Lock()
		m.messages = append(m.messages, line)
		m.mu.Unlock()

		// Send mock response
		fmt.Fprintf(conn, "{\"jsonrpc\":\"2.0\",\"result\":{\"status\":\"ok\"},\"id\":1}\n")
	}
}

func (m *MockTCPDAW) Stop() {
	close(m.stop)
}

func (m *MockTCPDAW) GetMessages() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string{}, m.messages...)
}
