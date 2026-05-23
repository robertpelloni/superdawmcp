package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type DashboardState struct {
	DAWs    map[string]DAWState `json:"daws"`
	Patches []AudioPatch       `json:"patches"`
	mu      sync.RWMutex
}

type DAWState struct {
	Name      string  `json:"name"`
	IsPlaying bool    `json:"is_playing"`
	BPM       float64 `json:"bpm"`
}

type AudioPatch struct {
	SourceDAW   string `json:"source_daw"`
	SourceTrack string `json:"source_track"`
	DestDAW     string `json:"dest_daw"`
	DestTrack   string `json:"dest_track"`
}

func StartDashboard(port int) *DashboardState {
	state := &DashboardState{
		DAWs:    make(map[string]DAWState),
		Patches: []AudioPatch{},
	}

	http.HandleFunc("/api/state", func(w http.ResponseWriter, r *http.Request) {
		state.mu.RLock()
		defer state.mu.RUnlock()
		json.NewEncoder(w).Encode(state)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
				<head>
					<title>SuperDAW Dashboard</title>
					<style>
						body { font-family: sans-serif; background: #1a1a1a; color: #eee; padding: 20px; }
						.card { background: #2a2a2a; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
						h2 { color: #00ff88; border-bottom: 1px solid #444; padding-bottom: 5px; }
						pre { background: #000; padding: 10px; border-radius: 4px; overflow: auto; }
						.grid { display: flex; gap: 20px; flex-wrap: wrap; }
					</style>
				</head>
				<body>
					<h1>SuperDAW Universal Dashboard</h1>
					<div class="grid">
						<div class="card" style="flex: 1; min-width: 300px;">
							<h2>DAW Status</h2>
							<div id="daws"></div>
						</div>
						<div class="card" style="flex: 1; min-width: 300px;">
							<h2>Audio Routing Matrix</h2>
							<div id="routing"></div>
						</div>
					</div>

					<script>
						async function update() {
							const res = await fetch('/api/state');
							const state = await res.json();

							document.getElementById('daws').innerHTML = '<pre>' + JSON.stringify(state.daws, null, 2) + '</pre>';
							document.getElementById('routing').innerHTML = '<pre>' + JSON.stringify(state.patches, null, 2) + '</pre>';
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

func (s *DashboardState) AddPatch(p AudioPatch) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Patches = append(s.Patches, p)
}
