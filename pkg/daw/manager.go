package daw

import (
	"fmt"
	"sync"
)

// ConnectionManager handles multiple concurrent DAW instances.
type ConnectionManager struct {
	instances map[string]DAWDriver
	defaultID string
	mu        sync.RWMutex
}

// NewConnectionManager creates a new instance manager.
func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		instances: make(map[string]DAWDriver),
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
