package engine

import (
	"fmt"
	"sync"
	"time"
)

// PeerSession represents an individual collaborator in a studio jam session.
type PeerSession struct {
	ID       string `json:"id"`
	Handle   string `json:"handle"`
	JoinedAt int64  `json:"joined_at"`
	IsHost   bool   `json:"is_host"`
}

// CollabParamState represents the synchronized state of a parameter with timestamp for LWW resolution.
type CollabParamState struct {
	PluginName string    `json:"plugin_name"`
	ParamIndex int       `json:"param_index"`
	Value      float32   `json:"value"`
	UpdatedBy  string    `json:"updated_by"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ParamUpdateEvent is emitted to all connected peers when a parameter changes.
type ParamUpdateEvent struct {
	SenderID   string    `json:"sender_id"`
	PluginName string    `json:"plugin_name"`
	ParamIndex int       `json:"param_index"`
	Value      float32   `json:"value"`
	Timestamp  time.Time `json:"timestamp"`
}

// CollabHub manages multi-user studio jam sessions and parameter synchronization across peers.
type CollabHub struct {
	mu         sync.RWMutex
	roomName   string
	peers      map[string]*PeerSession
	paramState map[string]CollabParamState // key: "pluginName:paramIndex"
	eventLog   []ParamUpdateEvent
	subscribers map[string]chan ParamUpdateEvent
}

// NewCollabHub creates a new collaborative studio session manager.
func NewCollabHub(roomName string) *CollabHub {
	if roomName == "" {
		roomName = "Studio-Jam-Main"
	}
	return &CollabHub{
		roomName:    roomName,
		peers:       make(map[string]*PeerSession),
		paramState:  make(map[string]CollabParamState),
		eventLog:    make([]ParamUpdateEvent, 0),
		subscribers: make(map[string]chan ParamUpdateEvent),
	}
}

// RegisterPeer adds a collaborator to the jam session and returns their event channel.
func (c *CollabHub) RegisterPeer(peerID, handle string, isHost bool) (chan ParamUpdateEvent, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if peerID == "" {
		return nil, fmt.Errorf("peerID cannot be empty")
	}

	session := &PeerSession{
		ID:       peerID,
		Handle:   handle,
		JoinedAt: time.Now().Unix(),
		IsHost:   isHost,
	}
	c.peers[peerID] = session

	ch := make(chan ParamUpdateEvent, 64)
	c.subscribers[peerID] = ch

	return ch, nil
}

// UnregisterPeer removes a collaborator from the session and cleans up their resources.
func (c *CollabHub) UnregisterPeer(peerID string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.peers, peerID)
	if ch, ok := c.subscribers[peerID]; ok {
		close(ch)
		delete(c.subscribers, peerID)
	}
}

// GetPeerCount returns the number of active collaborators currently in the session.
func (c *CollabHub) GetPeerCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.peers)
}

// GetPeers returns a slice of all active peer sessions.
func (c *CollabHub) GetPeers() []*PeerSession {
	c.mu.RLock()
	defer c.mu.RUnlock()

	list := make([]*PeerSession, 0, len(c.peers))
	for _, p := range c.peers {
		list = append(list, p)
	}
	return list
}

// BroadcastParamChange updates the parameter state with Last-Write-Wins (LWW) conflict resolution
// and fans out the update to all active subscribers.
func (c *CollabHub) BroadcastParamChange(senderPeerID, pluginName string, paramIndex int, value float32, ts time.Time) (bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ts.IsZero() {
		ts = time.Now()
	}

	key := fmt.Sprintf("%s:%d", pluginName, paramIndex)
	existing, exists := c.paramState[key]
	if exists && ts.Before(existing.UpdatedAt) {
		// Reject stale change: newer update already applied (LWW)
		return false, nil
	}

	// Update consensus state
	c.paramState[key] = CollabParamState{
		PluginName: pluginName,
		ParamIndex: paramIndex,
		Value:      value,
		UpdatedBy:  senderPeerID,
		UpdatedAt:  ts,
	}

	event := ParamUpdateEvent{
		SenderID:   senderPeerID,
		PluginName: pluginName,
		ParamIndex: paramIndex,
		Value:      value,
		Timestamp:  ts,
	}
	c.eventLog = append(c.eventLog, event)

	// Fan out to all other peers non-blockingly
	for peerID, ch := range c.subscribers {
		if peerID == senderPeerID {
			continue // Don't echo back to sender
		}
		select {
		case ch <- event:
		default:
			// Buffer full, drop or skip to avoid blocking other peers
		}
	}

	return true, nil
}

// GetParamState retrieves the current consensus state for a specific plugin parameter.
func (c *CollabHub) GetParamState(pluginName string, paramIndex int) (CollabParamState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	key := fmt.Sprintf("%s:%d", pluginName, paramIndex)
	state, ok := c.paramState[key]
	return state, ok
}
