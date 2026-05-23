package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type DashboardState struct {
	DAWs map[string]DAWState `json:"daws"`
	mu   sync.RWMutex
}

type DAWState struct {
	Name      string  `json:"name"`
	IsPlaying bool    `json:"is_playing"`
	BPM       float64 `json:"bpm"`
}

func StartDashboard(port int) *DashboardState {
	state := &DashboardState{
		DAWs: make(map[string]DAWState),
	}

	http.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		state.mu.RLock()
		defer state.mu.RUnlock()
		json.NewEncoder(w).Encode(state)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
				<head><title>SuperDAW Dashboard</title></head>
				<body>
					<h1>SuperDAW Dashboard</h1>
					<div id="state"></div>
					<script>
						async function update() {
							const res = await fetch('/api/state');
							const state = await res.json();
							document.getElementById('state').innerText = JSON.stringify(state, null, 2);
						}
						setInterval(update, 1000);
						update();
					</script>
				</body>
			</html>
		`)
	})

	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return state
}

func (s *DashboardState) UpdateDAW(name string, playing bool, bpm float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.DAWs[name] = DAWState{Name: name, IsPlaying: playing, BPM: bpm}
}
