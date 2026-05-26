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
	Patches []AudioPatch        `json:"patches"`
	Jobs    []interface{}       `json:"jobs"`
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

type DAWState struct {
	Name        string  `json:"name"`
	IsPlaying   bool    `json:"is_playing"`
	BPM         float64 `json:"bpm"`
	Arrangement string  `json:"arrangement"`
}

type AudioPatch struct {
	SourceDAW   string `json:"source_daw"`
	SourceTrack string `json:"source_track"`
	DestDAW     string `json:"dest_daw"`
	DestTrack   string `json:"dest_track"`
}

type InternalCommand struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

var CommandBus = make(chan InternalCommand, 32)

func StartDashboard(port int) (*DashboardState, *http.ServeMux) {
	state := &DashboardState{
		DAWs:    make(map[string]DAWState),
		Patches: []AudioPatch{},
		clients: make(map[*websocket.Conn]bool),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/call", func(w http.ResponseWriter, r *http.Request) {
		var req InternalCommand
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		CommandBus <- req
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil { return }
		state.mu.Lock()
		state.clients[conn] = true
		state.mu.Unlock()

		// Send initial state
		state.broadcast()
	})

	mux.HandleFunc("/obs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
				<head>
					<style>
						body { font-family: sans-serif; color: #fff; background: transparent; padding: 10px; font-weight: bold; text-shadow: 2px 2px 4px #000; }
						.obs-stat { font-size: 24px; color: #00ff88; }
						.obs-badge { background: rgba(0,0,0,0.5); padding: 5px 10px; border-radius: 8px; margin: 5px 0; border-left: 5px solid #00bcd4; }
					</style>
				</head>
				<body>
					<div id="stats"></div>
					<script>
						const ws = new WebSocket('ws://' + window.location.host + '/ws');
						ws.onmessage = (event) => {
							const state = JSON.parse(event.data);
							let html = '';
							for (const name in state.daws) {
								const d = state.daws[name];
								if (d.is_playing) {
									html += '<div class="obs-badge">DAW: ' + name.toUpperCase() + ' <span class="obs-stat">' + d.bpm.toFixed(1) + ' BPM</span></div>';
								}
							}
							document.getElementById('stats').innerHTML = html;
						};
					</script>
				</body>
			</html>
		`)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
						#timeline { width: 100%; height: 300px; background: #000; margin-top: 20px; border: 1px solid #444; position: relative; overflow-x: auto; }
						#blueprint { width: 100%; height: 200px; background: #1a1a1a; border: 1px dashed #444; margin-top: 10px; display: flex; align-items: center; justify-content: center; font-family: monospace; color: #00ff88; }
						.track-lane { height: 40px; border-bottom: 1px solid #222; display: flex; align-items: center; white-space: nowrap; }
						.clip-block { position: absolute; background: #00bcd4; height: 30px; border-radius: 4px; border: 1px solid #fff; font-size: 10px; color: #000; padding: 2px; overflow: hidden; }
						.keyboard { display: flex; justify-content: center; margin-top: 20px; }
						.key { width: 40px; height: 120px; border: 1px solid #000; background: white; cursor: pointer; }
						.key.black { background: black; height: 80px; width: 30px; margin-left: -15px; margin-right: -15px; z-index: 2; }
						.key:active { background: #00ff88; }
					</style>
				</head>
				<body>
					<div style="display: flex; justify-content: space-between; align-items: center;">
						<h1>SuperDAW Universal Dashboard</h1>
						<div>
							<button onclick="callMcp('superdaw_save_session', {})" class="badge" style="cursor: pointer; background: #00ff88; color: #000; border: none;">SAVE SESSION</button>
							<button onclick="callMcp('superdaw_load_session', {})" class="badge" style="cursor: pointer; background: #00bcd4; color: #000; border: none;">LOAD SESSION</button>
							<div id="version-badge" class="badge">v2.3.0 (Active)</div>
						</div>
					</div>

					<div class="grid">
						<div class="card">
							<h2>DAW Engine Status</h2>
							<div id="daws"></div>
						</div>
						<div class="card">
							<h2>Virtual Audio Patching</h2>
							<div id="routing"></div>
							<div id="blueprint"></div>
						</div>
						<div class="card">
							<h2>Generative AI Activity</h2>
							<div id="jobs"></div>
						</div>
					</div>

					<div class="card">
						<h2>Live Studio Arrangement</h2>
						<div id="timeline"></div>
					</div>

					<div class="card">
						<h2>Virtual MIDI Performance</h2>
						<div class="keyboard">
							<div class="key" onclick="playNote(60)"></div>
							<div class="key black" onclick="playNote(61)"></div>
							<div class="key" onclick="playNote(62)"></div>
							<div class="key black" onclick="playNote(63)"></div>
							<div class="key" onclick="playNote(64)"></div>
							<div class="key" onclick="playNote(65)"></div>
							<div class="key black" onclick="playNote(66)"></div>
							<div class="key" onclick="playNote(67)"></div>
							<div class="key black" onclick="playNote(68)"></div>
							<div class="key" onclick="playNote(69)"></div>
							<div class="key black" onclick="playNote(70)"></div>
							<div class="key" onclick="playNote(71)"></div>
							<div class="key" onclick="playNote(72)"></div>
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

							// Render Jobs
							let jobHtml = '';
							if (state.jobs) {
								state.jobs.forEach(j => {
									jobHtml += `
										<div class="patch-item">
											<span>${j.prompt}</span>
											<span class="badge" style="width: 100px; background: #444; position: relative; overflow: hidden;">
												<div style="background: #00ff88; width: ${j.progress*100}%; height: 10px; border-radius: 5px;"></div>
											</span>
											<span>${j.status}</span>
										</div>
									`;
								});
							}
							document.getElementById('jobs').innerHTML = jobHtml || '<p style="color: #666">No active generation jobs.</p>';

							// Render Blueprint (Mermaid-style text graph)
							let blueprint = 'graph LR\n';
							state.patches.forEach(p => {
								blueprint += `  ${p.source_daw} --> ${p.dest_daw}\n`;
							});
							document.getElementById('blueprint').innerText = blueprint === 'graph LR\n' ? 'No connections.' : blueprint;

							// Render Timeline
							let timelineHtml = '';
							let top = 0;
							for (const name in state.daws) {
								const d = state.daws[name];
								if (d.arrangement) {
									try {
										const arrangement = JSON.parse(d.arrangement);
										arrangement.forEach(track => {
											timelineHtml += `<div class="track-lane" style="top: ${top}px; position: relative;"><span style="width: 100px; display: inline-block;">${track.track}</span>`;
											track.clips.forEach(clip => {
												const left = clip.start * 20; // 20 pixels per second
												const width = (clip.end - clip.start) * 20;
												timelineHtml += `<div class="clip-block" style="left: ${100+left}px; width: ${width}px;">${clip.name}</div>`;
											});
											timelineHtml += `</div>`;
											top += 40;
										});
									} catch(e) {}
								}
							}
							document.getElementById('timeline').innerHTML = timelineHtml || '<p style="color: #666; padding: 20px;">No arrangement data available.</p>';
						};

						function playNote(pitch) {
							callMcp('superdaw_write_midi', {
								track_id: '0',
								notes: [{pitch: pitch, velocity: 100, start_beat: 0, duration: 0.5}]
							});
						}

						async function callMcp(name, args) {
							// For this simplified dashboard, we assume a local API proxy exists or
							// we just log the intent. In a real build, this uses the server's internal RPC.
							console.log("MCP Call:", name, args);
							fetch('/api/call', {
								method: 'POST',
								headers: {'Content-Type': 'application/json'},
								body: JSON.stringify({name, arguments: args})
							});
						}
					</script>
				</body>
			</html>
		`)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}
	go server.ListenAndServe()
	return state, mux
}

func GetMux() *http.ServeMux {
	// This is a bit of a hack to allow remote.go to register on the same mux if needed,
	// but for now we'll just use the global http.DefaultServeMux in remote.go if it's separate,
	// or better, refactor to use a single mux.
	return nil
}

func (s *DashboardState) UpdateDAW(name string, playing bool, bpm float64) {
	s.mu.Lock()
	state := s.DAWs[name]
	state.Name = name
	state.IsPlaying = playing
	state.BPM = bpm
	s.DAWs[name] = state
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) UpdateJobs(jobs []interface{}) {
	s.mu.Lock()
	s.Jobs = jobs
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) GetState() DashboardState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s
}

func (s *DashboardState) SetState(state DashboardState) {
	s.mu.Lock()
	s.DAWs = state.DAWs
	s.Patches = state.Patches
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) UpdateArrangement(name string, arrangement string) {
	s.mu.Lock()
	state := s.DAWs[name]
	state.Arrangement = arrangement
	s.DAWs[name] = state
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
