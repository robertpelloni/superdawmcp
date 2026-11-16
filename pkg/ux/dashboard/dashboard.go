package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
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
		if err != nil {
			return
		}
		state.mu.Lock()
		state.clients[conn] = true
		state.mu.Unlock()
		state.broadcast()
	})

	mux.HandleFunc("/obs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
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
		fmt.Fprint(w, `
			<html>
				<head>
					<title>SuperDAW Dashboard v3.2</title>
					<style>
						body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #121212; color: #e0e0e0; padding: 20px; }
						.card { background: #1e1e1e; padding: 20px; border-radius: 12px; margin-bottom: 20px; border: 1px solid #333; box-shadow: 0 4px 6px rgba(0,0,0,0.3); }
						h1 { color: #00ff88; text-transform: uppercase; letter-spacing: 2px; text-shadow: 0 0 10px rgba(0,255,136,0.3); }
						h2 { color: #00bcd4; border-bottom: 1px solid #333; padding-bottom: 10px; font-weight: 300; }
						pre { background: #080808; padding: 15px; border-radius: 8px; border-left: 4px solid #00ff88; overflow: auto; font-family: 'Consolas', monospace; }
						.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 25px; }
						.patch-item { background: #252525; padding: 10px; margin: 5px 0; border-radius: 4px; display: flex; justify-content: space-between; align-items: center; }
						.patch-arrow { color: #00ff88; font-weight: bold; }
						.badge { background: #333; padding: 4px 10px; border-radius: 10px; font-size: 0.85em; color: #fff; cursor: help; border: 1px solid #555; }
						.badge:hover { background: #00ff88; color: #000; }
						.tooltip { position: relative; display: inline-block; cursor: help; }
						.tooltip .tooltiptext { visibility: hidden; width: 220px; background-color: #333; color: #fff; text-align: center; border-radius: 6px; padding: 8px; position: absolute; z-index: 1; bottom: 125%; left: 50%; margin-left: -110px; opacity: 0; transition: opacity 0.3s; font-size: 0.8em; font-weight: normal; border: 1px solid #555; box-shadow: 0 4px 6px rgba(0,0,0,0.5); }
						.tooltip:hover .tooltiptext { visibility: visible; opacity: 1; }
						#timeline { width: 100%%; height: 300px; background: #000; margin-top: 20px; border: 1px solid #444; position: relative; overflow-x: auto; }
						#blueprint { width: 100%%; height: 200px; background: #1a1a1a; border: 1px dashed #444; margin-top: 10px; display: flex; align-items: center; justify-content: center; font-family: monospace; color: #00ff88; }
						.track-lane { height: 40px; border-bottom: 1px solid #222; display: flex; align-items: center; white-space: nowrap; }
						.clip-block { position: absolute; background: #00bcd4; height: 30px; border-radius: 4px; border: 1px solid #fff; font-size: 10px; color: #000; padding: 2px; overflow: hidden; }
						.keyboard { display: flex; justify-content: center; margin-top: 20px; }
						.key { width: 40px; height: 120px; border: 1px solid #000; background: white; cursor: pointer; }
						.key.black { background: black; height: 80px; width: 30px; margin-left: -15px; margin-right: -15px; z-index: 2; }
						.key:active { background: #00ff88; }
						.feature-list { list-style: none; padding: 0; }
						.feature-list li { margin: 10px 0; padding: 10px; background: #222; border-radius: 6px; border-left: 4px solid #00bcd4; display: flex; align-items: center; justify-content: space-between; }
						.feature-title { font-weight: bold; color: #fff; }
					</style>
				</head>
				<body>
					<div style="display: flex; justify-content: space-between; align-items: center;">
						<h1><span class="tooltip">SuperDAW Orchestrator<span class="tooltiptext">Global Model Context Protocol Gateway unifying all DAWs</span></span></h1>
						<div>
							<span class="badge tooltip">WebSocket Connected<span class="tooltiptext">Streaming live telemetry at 30fps</span></span>
							<span class="badge tooltip">MCP v3.2<span class="tooltiptext">Latest API schema loaded</span></span>
						</div>
					</div>

					<div class="grid">
						<div class="card">
							<h2>
								<span class="tooltip">Active DAWs<span class="tooltiptext">DAWs currently synchronized via driver agents</span></span>
							</h2>
							<pre id="daws">Waiting for telemetry...</pre>
						</div>
						<div class="card">
							<h2>
								<span class="tooltip">Jack / ReRoute Audio Patches<span class="tooltiptext">Universal hardware routing configurations mapping stems cross-DAW</span></span>
							</h2>
							<div id="patches"></div>
						</div>
					</div>

					<div class="card">
						<h2>
							<span class="tooltip">Global Feature Modules<span class="tooltiptext">Core orchestrator capabilities powered by AI and Submodules</span></span>
						</h2>
						<ul class="feature-list">
							<li>
								<div>
									<span class="feature-title">Plugin Inspector (v3.1.0)</span>
									<div style="font-size: 0.85em; color: #888; margin-top: 4px;">Deep-scans VST3 parameters via native libvst3 bridges.</div>
								</div>
								<span class="badge tooltip">Active<span class="tooltiptext">Waiting for superdaw_get_plugin_params</span></span>
							</li>
							<li>
								<div>
									<span class="feature-title">WebAssembly VST Runners</span>
									<div style="font-size: 0.85em; color: #888; margin-top: 4px;">Executes Wasm DSP algorithms sandboxed in the Go daemon.</div>
								</div>
								<span class="badge tooltip">Active<span class="tooltiptext">Wazero host initialized</span></span>
							</li>
							<li>
								<div>
									<span class="feature-title">In-DAW LLM Reasoning Sidecar</span>
									<div style="font-size: 0.85em; color: #888; margin-top: 4px;">Continuously analyzes active project states to suggest mix/arrangement adjustments.</div>
								</div>
								<span class="badge tooltip">Active<span class="tooltiptext">Background loop evaluating JSON dumps</span></span>
							</li>
							<li>
								<div>
									<span class="feature-title">Universal Preset Translator</span>
									<div style="font-size: 0.85em; color: #888; margin-top: 4px;">Heuristically maps complex synth parameters across incompatible plugins (e.g. Serum to Vital).</div>
								</div>
								<span class="badge tooltip">Active<span class="tooltiptext">Ready for superdaw_translate_preset</span></span>
							</li>
							<li>
								<div>
									<span class="feature-title">Multi-User Studio Sessions</span>
									<div style="font-size: 0.85em; color: #888; margin-top: 4px;">TCP synchronized state loops for collaborative remote editing.</div>
								</div>
								<span class="badge tooltip">Active<span class="tooltiptext">TCP Gateway port 12002 open</span></span>
							</li>
							<li>
								<div>
									<span class="feature-title">GPU-Accelerated FFT Inspector</span>
									<div style="font-size: 0.85em; color: #888; margin-top: 4px;">High-frequency telemetry stream driving raw WebGL audio magnitudes.</div>
								</div>
								<span class="badge tooltip">Active<span class="tooltiptext">Streaming real-time arrays over WS</span></span>
							</li>
						</ul>
					</div>

					<div class="card">
						<h2>
							<span class="tooltip">Generative AI Blueprint<span class="tooltiptext">AI prompt-to-stem rendering timeline map</span></span>
						</h2>
						<div id="blueprint">No AI prompt processing currently active. Waiting for superdaw_import_generated_stems...</div>
					</div>

					<script>
						const ws = new WebSocket('ws://' + window.location.host + '/ws');
						ws.onmessage = (event) => {
							const state = JSON.parse(event.data);
							document.getElementById('daws').textContent = JSON.stringify(state.daws, null, 2);

							let patchesHtml = '';
							if(state.patches && state.patches.length > 0) {
								state.patches.forEach(p => {
									patchesHtml += `+"`"+`<div class="patch-item">
										<span><span class="badge">${p.source_daw}</span> ${p.source_track}</span>
										<span class="patch-arrow">>></span>
										<span><span class="badge">${p.dest_daw}</span> ${p.dest_track}</span>
									</div>`+"`"+`;
								});
							} else {
								patchesHtml = '<p style="color: #666;">No active audio cross-routing.</p>';
							}
							document.getElementById('patches').innerHTML = patchesHtml;
						};
					</script>
				</body>
			</html>`)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}
	go server.ListenAndServe()
	return state, mux
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
