package engine

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type ReasoningSidecar struct {
	Active      bool
	Context     []string
	Suggestions []string
	lock        sync.RWMutex
}

var GlobalReasoningSidecar = &ReasoningSidecar{
	Active:      false,
	Context:     make([]string, 0),
	Suggestions: make([]string, 0),
}

func (r *ReasoningSidecar) StartBackgroundLoop() {
	r.lock.Lock()
	if r.Active {
		r.lock.Unlock()
		return
	}
	r.Active = true
	r.lock.Unlock()

	go func() {
		for {
			r.lock.RLock()
			active := r.Active
			r.lock.RUnlock()
			if !active {
				break
			}

			// Simulate reasoning context tick
			time.Sleep(5 * time.Second)
		}
	}()
}

func (r *ReasoningSidecar) StopBackgroundLoop() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.Active = false
}

func (r *ReasoningSidecar) AnalyzeProjectState(stateDump string) string {
	r.lock.Lock()
	defer r.lock.Unlock()

	// Basic heuristic mock reasoning over a project JSON dump
	var issues []string
	if strings.Contains(strings.ToLower(stateDump), "muddy") {
		issues = append(issues, "Detected potential frequency masking in the lower mids (200-400Hz). Suggest cutting here on the bass track.")
	}
	if strings.Contains(strings.ToLower(stateDump), "clipping") {
		issues = append(issues, "Master bus appears to be peaking over 0dBFS. Suggest applying a limiter or reducing gain staging.")
	}
	if !strings.Contains(strings.ToLower(stateDump), "reverb") {
		issues = append(issues, "Mix feels dry. Suggest adding a subtle room or plate reverb on a return track.")
	}

	if len(issues) == 0 {
		issues = append(issues, "Project structure looks clean. Consider adding automation to build tension.")
	}

	r.Suggestions = issues
	return fmt.Sprintf("Analysis complete. Found %d insights.", len(issues))
}

func (r *ReasoningSidecar) GetSuggestions() []string {
	r.lock.RLock()
	defer r.lock.RUnlock()
	return r.Suggestions
}
