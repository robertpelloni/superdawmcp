package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type DashboardState struct {
	DAWs    map[string]DAWState `json:"daws"`
	Patches []AudioPatch       `json:"patches"`
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
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
		clients: make(map[*websocket.Conn]bool),
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil { return }
		state.mu.Lock()
		state.clients[conn] = true
		state.mu.Unlock()

		// Send initial state
		state.broadcast()
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
				<head>
					<title>SuperDAW Dashboard v1.6</title>
					<style>
						body { font-family: sans-serif; background: #1a1a1a; color: #eee; padding: 20px; }
						.card { background: #2a2a2a; padding: 15px; border-radius: 8px; margin-bottom: 20px; }
						h2 { color: #00ff88; border-bottom: 1px solid #444; padding-bottom: 5px; }
						pre { background: #000; padding: 10px; border-radius: 4px; overflow: auto; }
						.grid { display: flex; gap: 20px; flex-wrap: wrap; }
					</style>
				</head>
				<body>
					<h1>SuperDAW Universal Dashboard (WebSocket)</h1>
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
						const ws = new WebSocket('ws://' + window.location.host + '/ws');
						ws.onmessage = (event) => {
							const state = JSON.parse(event.data);
							document.getElementById('daws').innerHTML = '<pre>' + JSON.stringify(state.daws, null, 2) + '</pre>';
							document.getElementById('routing').innerHTML = '<pre>' + JSON.stringify(state.patches, null, 2) + '</pre>';
						};
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
	s.DAWs[name] = DAWState{Name: name, IsPlaying: playing, BPM: bpm}
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) AddPatch(p AudioPatch) {
	s.mu.Lock()
	s.Patches = append(s.Patches, p)
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) broadcast() {
	s.mu.RLock()
	data, _ := json.Marshal(s)
	for conn := range s.clients {
		conn.WriteMessage(websocket.TextMessage, data)
	}
	s.mu.RUnlock()
}
