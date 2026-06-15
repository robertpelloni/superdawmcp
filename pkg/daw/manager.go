package daw

import (
	"fmt"
	"sync"
)

// TransportState holds the last-known playback state for a DAW instance.
type TransportState struct {
	Playing bool
	BPM     float64
}

// ConnectionManager handles multiple concurrent DAW instances.
type ConnectionManager struct {
	instances         map[string]DAWDriver
	defaultID         string
	transportStates   map[string]TransportState
	mu                sync.RWMutex
}

// NewConnectionManager creates a new instance manager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		instances:       make(map[string]DAWDriver),
		transportStates: make(map[string]TransportState),
	}
}

// Register adds a new DAW driver instance to the manager.
func (m *ConnectionManager) Register(id string, driver DAWDriver) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.instances[id] = driver
	if m.defaultID == "" {
		m.defaultID = id
	}
}

// Unregister removes a DAW instance.
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

// SetDefault sets the default DAW instance to use when none is specified.
func (m *ConnectionManager) SetDefault(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.instances[id]; !ok {
		return fmt.Errorf("instance %s not found", id)
	}
	m.defaultID = id
	return nil
}

// Get retrieves a DAW driver by its unique ID.
func (m *ConnectionManager) Get(id string) (DAWDriver, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if id == "" {
		id = m.defaultID
	}
	driver, ok := m.instances[id]
	if !ok {
		return nil, fmt.Errorf("DAW instance '%s' not found", id)
	}
	return driver, nil
}

// ResolveID returns the effective instance ID (resolves empty to default).
func (m *ConnectionManager) ResolveID(id string) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if id == "" {
		return m.defaultID
	}
	return id
}

// CacheTransportState stores the last known playback state for a DAW instance.
func (m *ConnectionManager) CacheTransportState(id string, playing bool, bpm float64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.transportStates[id] = TransportState{Playing: playing, BPM: bpm}
}

// GetCachedTransportState retrieves a previously cached transport state.
func (m *ConnectionManager) GetCachedTransportState(id string) (TransportState, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s, ok := m.transportStates[id]
	return s, ok
}

// GetAll returns a map of all registered DAW instances.
func (m *ConnectionManager) GetAll() map[string]DAWDriver {
	m.mu.RLock()
	defer m.mu.RUnlock()
	res := make(map[string]DAWDriver)
	for k, v := range m.instances {
		res[k] = v
	}
	return res
}
