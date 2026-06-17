package daw

import (
	"fmt"
	"sync"
)

type TransportState struct {
	Playing bool
	BPM     float64
}

type ConnectionInfo struct {
	Driver DAWDriver
	RoomID string
}

type ConnectionManager struct {
	instances       map[string]ConnectionInfo
	defaultID       string
	transportStates map[string]TransportState
	mu              sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		instances:       make(map[string]ConnectionInfo),
		transportStates: make(map[string]TransportState),
	}
}

func (m *ConnectionManager) Register(id string, driver DAWDriver) {
	m.RegisterWithRoom(id, driver, "default")
}

func (m *ConnectionManager) RegisterWithRoom(id string, driver DAWDriver, roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if roomID == "" { roomID = "default" }
	m.instances[id] = ConnectionInfo{Driver: driver, RoomID: roomID}
	if m.defaultID == "" {
		m.defaultID = id
	}
}

func (m *ConnectionManager) Unregister(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.instances, id)
	if m.defaultID == id {
		m.defaultID = ""
		for k := range m.instances {
			m.defaultID = k
			break
		}
	}
}

func (m *ConnectionManager) SetDefault(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.instances[id]; !ok {
		return fmt.Errorf("instance %s not found", id)
	}
	m.defaultID = id
	return nil
}

func (m *ConnectionManager) Get(id string) (DAWDriver, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if id == "" {
		id = m.defaultID
	}
	info, ok := m.instances[id]
	if !ok {
		return nil, fmt.Errorf("DAW instance '%s' not found", id)
	}
	return info.Driver, nil
}

func (m *ConnectionManager) GetInfo(id string) (ConnectionInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if id == "" {
		id = m.defaultID
	}
	info, ok := m.instances[id]
	if !ok {
		return ConnectionInfo{}, fmt.Errorf("DAW instance '%s' not found", id)
	}
	return info, nil
}

func (m *ConnectionManager) ResolveID(id string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if id == "" {
		return m.defaultID
	}
	return id
}

func (m *ConnectionManager) CacheTransportState(id string, playing bool, bpm float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transportStates[id] = TransportState{Playing: playing, BPM: bpm}
}

func (m *ConnectionManager) GetCachedTransportState(id string) (TransportState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.transportStates[id]
	return s, ok
}

func (m *ConnectionManager) GetAll() map[string]DAWDriver {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make(map[string]DAWDriver)
	for k, v := range m.instances {
		res[k] = v.Driver
	}
	return res
}
