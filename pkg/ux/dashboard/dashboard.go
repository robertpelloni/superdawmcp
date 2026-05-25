package dashboard

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
						body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #121212; color: #e0e0e0; padding: 20px; }
						.card { background: #1e1e1e; padding: 20px; border-radius: 12px; margin-bottom: 20px; border: 1px solid #333; box-shadow: 0 4px 6px rgba(0,0,0,0.3); }
						h1 { color: #00ff88; text-transform: uppercase; letter-spacing: 2px; text-shadow: 0 0 10px rgba(0,255,136,0.3); }
						h2 { color: #00bcd4; border-bottom: 1px solid #333; padding-bottom: 10px; font-weight: 300; }
						pre { background: #080808; padding: 15px; border-radius: 8px; border-left: 4px solid #00ff88; overflow: auto; font-family: 'Consolas', monospace; }
						.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 25px; }
						.patch-item { background: #252525; padding: 10px; margin: 5px 0; border-radius: 4px; display: flex; justify-content: space-between; align-items: center; }
						.patch-arrow { color: #00ff88; font-weight: bold; }
						.badge { background: #333; padding: 2px 8px; border-radius: 10px; font-size: 0.8em; color: #aaa; }
					</style>
				</head>
				<body>
					<div style="display: flex; justify-content: space-between; align-items: center;">
						<h1>SuperDAW Universal Dashboard</h1>
						<div id="version-badge" class="badge">v2.3.0 (Active)</div>
					</div>

					<div class="grid">
						<div class="card">
							<h2>DAW Engine Status</h2>
							<div id="daws"></div>
						</div>
						<div class="card">
							<h2>Virtual Audio Patching</h2>
							<div id="routing"></div>
						</div>
					</div>

					<script>
						const ws = new WebSocket('ws://' + window.location.host + '/ws');
						ws.onmessage = (event) => {
							const state = JSON.parse(event.data);

							// Render DAWs
							let dawHtml = '';
							for (const name in state.daws) {
								const d = state.daws[name];
								dawHtml += `
									<div class="patch-item">
										<span><strong>${name.toUpperCase()}</strong></span>
										<span>${d.is_playing ? '▶️ PLAYING' : '⏹️ STOPPED'}</span>
										<span class="badge">${d.bpm.toFixed(1)} BPM</span>
									</div>
								`;
							}
							document.getElementById('daws').innerHTML = dawHtml || '<p style="color: #666">No active DAWs connected.</p>';

							// Render Patches
							let patchHtml = '';
							state.patches.forEach(p => {
								patchHtml += `
									<div class="patch-item">
										<span>${p.source_daw} (${p.source_track})</span>
										<span class="patch-arrow">➔</span>
										<span>${p.dest_daw} (${p.dest_track})</span>
									</div>
								`;
							});
							document.getElementById('routing').innerHTML = patchHtml || '<p style="color: #666">No active audio patches.</p>';
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

func (s *DashboardState) RemovePatch(p AudioPatch) {
	s.mu.Lock()
	newPatches := []AudioPatch{}
	for _, patch := range s.Patches {
		if patch.SourceDAW == p.SourceDAW && patch.SourceTrack == p.SourceTrack &&
			patch.DestDAW == p.DestDAW && patch.DestTrack == p.DestTrack {
			continue
		}
		newPatches = append(newPatches, patch)
	}
	s.Patches = newPatches
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
