package daw

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type BitwigDriver struct {
	addr string
	conn net.Conn
	mu   sync.Mutex
	id   int
}

func NewBitwigDriver(host string, port int) *BitwigDriver {
	return &BitwigDriver{
		addr: fmt.Sprintf("%s:%d", host, port),
		id:   1,
	}
}

func (b *BitwigDriver) Connect(endpoint string) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	conn, err := net.Dial("tcp", b.addr)
	if err != nil {
		return err
	}
	b.conn = conn
	return nil
}

func (b *BitwigDriver) Disconnect() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.conn != nil {
		err := b.conn.Close()
		b.conn = nil
		return err
	}
	return nil
}

func (b *BitwigDriver) call(method string, params map[string]interface{}) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.conn == nil {
		conn, err := net.Dial("tcp", b.addr)
		if err != nil {
			return err
		}
		b.conn = conn
	}

	req := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      b.id,
	}
	b.id++

	data, _ := json.Marshal(req)
	_, err := fmt.Fprintf(b.conn, "%s\n", string(data))
	if err != nil {
		b.conn.Close()
		b.conn = nil
		return err
	}

	// Bitwig responses are single-line JSON
	reader := bufio.NewReader(b.conn)
	_, err = reader.ReadString('\n')
	return err
}

func (b *BitwigDriver) SetTransportState(playing bool, bpm float64) error {
	method := "transport.stop"
	if playing {
		method = "transport.play"
	}
	return b.call(method, map[string]interface{}{"bpm": bpm})
}

func (b *BitwigDriver) CreateTrack(name string, trackType string) (string, error) {
	err := b.call("track.create", map[string]interface{}{"name": name, "type": trackType})
	return "bitwig_track", err
}

func (b *BitwigDriver) SetTrackVolume(trackID string, volume float32) error {
	return b.call("track.volume", map[string]interface{}{"index": trackID, "volume": volume})
}

func (b *BitwigDriver) SetTrackPan(trackID string, pan float32) error {
	return b.call("track.pan", map[string]interface{}{"index": trackID, "pan": pan})
}

func (b *BitwigDriver) WriteMIDIClip(trackID string, clipIndex int, notes []MIDINote) error {
	return b.call("clip.set_notes", map[string]interface{}{
		"trackIndex": trackID,
		"slotIndex":  clipIndex,
		"notes":      notes,
	})
}
