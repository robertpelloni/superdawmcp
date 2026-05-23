package connector

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

// Connector enables bidirectional communication between DAW plugins and the MCP server.
type Connector struct {
	listenPort int
	clients    map[string]net.Conn
	mu         sync.RWMutex
	onUpdate   func(pluginID string, state map[string]interface{})
}

func NewConnector(port int) *Connector {
	return &Connector{
		listenPort: port,
		clients:    make(map[string]net.Conn),
	}
}

func (c *Connector) Start() error {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", c.listenPort))
	if err != nil { return err }

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil { continue }
			go c.handleClient(conn)
		}
	}()

	return nil
}

func (c *Connector) SetOnUpdate(callback func(string, map[string]interface{})) {
	c.onUpdate = callback
}

func (c *Connector) handleClient(conn net.Conn) {
	decoder := json.NewDecoder(conn)
	for {
		var msg struct {
			PluginID string                 `json:"plugin_id"`
			State    map[string]interface{} `json:"state"`
		}
		if err := decoder.Decode(&msg); err != nil {
			break
		}

		c.mu.Lock()
		c.clients[msg.PluginID] = conn
		c.mu.Unlock()

		if c.onUpdate != nil {
			c.onUpdate(msg.PluginID, msg.State)
		}
	}
	conn.Close()
}

// SendCommand sends a state update or command back to a specific plugin instance.
func (c *Connector) SendCommand(pluginID string, cmd map[string]interface{}) error {
	c.mu.RLock()
	conn, ok := c.clients[pluginID]
	c.mu.RUnlock()

	if !ok { return fmt.Errorf("plugin %s not connected", pluginID) }

	data, _ := json.Marshal(cmd)
	_, err := conn.Write(append(data, '\n'))
	return err
}
