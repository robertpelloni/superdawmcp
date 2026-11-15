package engine

import (
	"encoding/json"
	"sync"
	"time"
)

type UserSession struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Role     string    `json:"role"` // "admin", "viewer", "collaborator"
	JoinedAt time.Time `json:"joined_at"`
	LastSeen time.Time `json:"last_seen"`
}

type StudioSession struct {
	SessionID   string                 `json:"session_id"`
	ProjectName string                 `json:"project_name"`
	Users       map[string]UserSession `json:"users"`
	State       map[string]interface{} `json:"shared_state"`
	lock        sync.RWMutex
}

var GlobalStudioSession *StudioSession

func init() {
	GlobalStudioSession = &StudioSession{
		SessionID:   "studio-001",
		ProjectName: "SuperDAW Multi-User Project",
		Users:       make(map[string]UserSession),
		State:       make(map[string]interface{}),
	}
}

func (s *StudioSession) Join(userID, name, role string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Users[userID] = UserSession{
		ID:       userID,
		Name:     name,
		Role:     role,
		JoinedAt: time.Now(),
		LastSeen: time.Now(),
	}
}

func (s *StudioSession) Ping(userID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if user, exists := s.Users[userID]; exists {
		user.LastSeen = time.Now()
		s.Users[userID] = user
	}
}

func (s *StudioSession) Leave(userID string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.Users, userID)
}

func (s *StudioSession) UpdateState(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.State[key] = value
}

func (s *StudioSession) GetSessionDump() []byte {
	s.lock.RLock()
	defer s.lock.RUnlock()
	data, _ := json.Marshal(s)
	return data
}
