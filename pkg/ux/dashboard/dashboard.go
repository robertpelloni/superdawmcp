package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"github.com/gorilla/websocket"
	"html"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// RoomState holds the state and connected clients for a specific collaborative session
type RoomState struct {
	ID      string              `json:"id"`
	DAWs    map[string]DAWState `json:"daws"`
	Patches []AudioPatch        `json:"patches"`
	Jobs    []interface{}       `json:"jobs"`
	clients map[*websocket.Conn]bool
	mu      sync.RWMutex
}

type DashboardState struct {
	rooms map[string]*RoomState
	mu    sync.RWMutex
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
	RoomID    string                 `json:"room_id"`
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

var CommandBus = make(chan InternalCommand, 32)

func StartDashboard(port int) (*DashboardState, *http.ServeMux) {
	orchestrator := &DashboardState{
		rooms: make(map[string]*RoomState),
	}
	// Create default room
	orchestrator.getOrCreateRoom("default")

	mux := http.NewServeMux()

	mux.HandleFunc("/api/call", func(w http.ResponseWriter, r *http.Request) {
		var req InternalCommand
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.RoomID == "" { req.RoomID = "default" }
		CommandBus <- req
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("room")
		if roomID == "" { roomID = "default" }

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil { return }

		room := orchestrator.getOrCreateRoom(roomID)
		room.mu.Lock()
		room.clients[conn] = true
		room.mu.Unlock()

		// Clean up on disconnect
		defer func() {
			room.mu.Lock()
			delete(room.clients, conn)
			room.mu.Unlock()
			conn.Close()
		}()

		room.broadcast()

		// Keep connection alive and read messages (though we mostly push state)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				break
			}
		}
	})

	mux.HandleFunc("/obs", func(w http.ResponseWriter, r *http.Request) {
		roomID := r.URL.Query().Get("room")
		if roomID == "" { roomID = "default" }
		// Escape the roomID to prevent XSS
		safeRoomID := html.EscapeString(roomID)
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
						const roomID = '%s';
						const ws = new WebSocket('ws://' + window.location.host + '/ws?room=' + roomID);
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
		`, safeRoomID)
	})

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
				<head>
					<title>SuperDAW Dashboard v3.2</title>
					<style>
						body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; background: #121212; color: #e0e0e0; padding: 20px; }
						.card { background: #1e1e1e; padding: 20px; border-radius: 12px; margin-bottom: 20px; border: 1px solid #333; box-shadow: 0 4px 6px rgba(0,0,0,0.3); }
						h1 { color: #00ff88; text-transform: uppercase; letter-spacing: 2px; text-shadow: 0 0 10px rgba(0,255,136,0.3); }
						h2 { color: #00bcd4; border-bottom: 1px solid #333; padding-bottom: 10px; font-weight: 300; }
						h3 { font-size: 14px; color: #888; text-transform: uppercase; margin-bottom: 10px; }
						pre { background: #080808; padding: 15px; border-radius: 8px; border-left: 4px solid #00ff88; overflow: auto; font-family: 'Consolas', monospace; }
						.grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(400px, 1fr)); gap: 25px; }
						.patch-item { background: #252525; padding: 10px; margin: 5px 0; border-radius: 4px; display: flex; justify-content: space-between; align-items: center; }
						.patch-arrow { color: #00ff88; font-weight: bold; }
						.badge { background: #333; padding: 2px 8px; border-radius: 10px; font-size: 0.8em; color: #aaa; }
						#timeline { width: 100%%; height: 300px; background: #000; margin-top: 20px; border: 1px solid #444; position: relative; overflow-x: auto; }
						#blueprint { width: 100%%; height: 200px; background: #1a1a1a; border: 1px dashed #444; margin-top: 10px; display: flex; align-items: center; justify-content: center; font-family: monospace; color: #00ff88; }
						.track-lane { height: 40px; border-bottom: 1px solid #222; display: flex; align-items: center; white-space: nowrap; }
						.clip-block { position: absolute; background: #00bcd4; height: 30px; border-radius: 4px; border: 1px solid #fff; font-size: 10px; color: #000; padding: 2px; overflow: hidden; }
						.keyboard { display: flex; justify-content: center; margin-top: 20px; }
						.key { width: 40px; height: 120px; border: 1px solid #000; background: white; cursor: pointer; }
						.key.black { background: black; height: 80px; width: 30px; margin-left: -15px; margin-right: -15px; z-index: 2; }
						.key:active { background: #00ff88; }
						input, select { padding: 8px; background: #222; border: 1px solid #444; color: #fff; border-radius: 4px; }
						button.action { cursor: pointer; background: #00ff88; color: #000; border: none; padding: 8px 15px; border-radius: 4px; font-weight: bold; }
						.room-info { background: #333; padding: 5px 15px; border-radius: 20px; font-size: 14px; color: #00ff88; display: flex; align-items: center; gap: 10px; }
						.presence-indicator { width: 10px; height: 10px; background: #00ff88; border-radius: 50%%; display: inline-block; }
					</style>
				</head>
				<body>
					<div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px;">
						<h1>SuperDAW Universal Dashboard</h1>
						<div style="display: flex; gap: 10px; align-items: center;">
							<div class="room-info">
								Room: <input id="room-input" style="background: none; border: none; color: #fff; width: 80px; padding: 0;" value="default" onchange="switchRoom(this.value)">
								<span id="presence-count" class="badge">1 users</span>
								<div class="presence-indicator"></div>
							</div>
							<button onclick="callMcp('superdaw_save_session', {})" class="badge" style="cursor: pointer; background: #00ff88; color: #000; border: none;">SAVE SESSION</button>
							<button onclick="callMcp('superdaw_load_session', {})" class="badge" style="cursor: pointer; background: #00bcd4; color: #000; border: none;">LOAD SESSION</button>
							<div id="version-badge" class="badge">v3.2.0 (Alpha)</div>
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
							<div style="margin-top: 15px; border-top: 1px solid #333; padding-top: 10px;">
								<h3>New Patch</h3>
								<div style="display: flex; gap: 5px; flex-wrap: wrap;">
									<input id="src-daw" placeholder="Source DAW" style="flex: 1;">
									<input id="src-track" placeholder="Source Track" style="flex: 1;">
									<div class="patch-arrow" style="align-self: center;">➔</div>
									<input id="dst-daw" placeholder="Dest DAW" style="flex: 1;">
									<input id="dst-track" placeholder="Dest Track" style="flex: 1;">
									<button onclick="addPatch()" class="action">ADD</button>
								</div>
							</div>
							<div id="blueprint"></div>
						</div>
						<div class="card">
							<h2>Plugin Inspector</h2>
							<div id="plugin-selector">
								<select id="plugin-list" onchange="loadPluginParams(this.value)" style="width: 100%%;">
									<option>Select a plugin...</option>
								</select>
							</div>
							<div id="plugin-params" style="margin-top: 15px;"></div>
						</div>
						<div class="card">
							<h2>MIDI CC & Custom Commands</h2>
							<div style="margin-bottom: 20px;">
								<h3>Send MIDI CC</h3>
								<div style="display: flex; gap: 5px;">
									<input id="cc-track" placeholder="Trk ID" style="width: 60px;">
									<input id="cc-num" placeholder="CC #" style="width: 60px;">
									<input id="cc-val" placeholder="Value" style="width: 60px;">
									<button onclick="sendCC()" class="action">SEND</button>
								</div>
							</div>
							<div style="border-top: 1px solid #333; padding-top: 10px;">
								<h3>Custom Command</h3>
								<div style="display: flex; gap: 5px; flex-direction: column;">
									<input id="custom-cmd" placeholder="Command Name">
									<textarea id="custom-args" placeholder='{"arg1": "val"}' style="background: #222; border: 1px solid #444; color: #fff; border-radius: 4px; padding: 8px;"></textarea>
									<button onclick="sendCustom()" class="action">EXECUTE</button>
								</div>
							</div>
						</div>
						<div class="card">
							<h2>Generative AI Activities</h2>
							<div id="jobs"></div>
							<div style="margin-top: 15px; border-top: 1px solid #333; padding-top: 10px;">
								<h3>Prompt-to-Stem</h3>
								<div style="display: flex; gap: 5px; flex-direction: column;">
									<input id="gen-prompt" placeholder="Symphonic psytrance with tribal drums...">
									<button onclick="importGenerative()" class="action">GENERATE & IMPORT</button>
								</div>
							</div>
						</div>
						<div class="card">
							<h2>Music Theory & Algorithms</h2>
							<div style="margin-bottom: 20px;">
								<h3>Euclidean Rhythm</h3>
								<div style="display: flex; gap: 5px; flex-wrap: wrap;">
									<input id="euc-track" placeholder="Trk ID" style="width: 60px;">
									<input id="euc-hits" placeholder="Hits" style="width: 60px;">
									<input id="euc-steps" placeholder="Steps" style="width: 60px;">
									<input id="euc-pitch" placeholder="Pitch" style="width: 60px;">
									<button onclick="generateEuclidean()" class="action">GENERATE</button>
								</div>
							</div>
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
						let currentRoom = new URLSearchParams(window.location.search).get('room') || 'default';
						document.getElementById('room-input').value = currentRoom;

						let ws;
						function connectWS() {
							if (ws) ws.close();
							ws = new WebSocket('ws://' + window.location.host + '/ws?room=' + encodeURIComponent(currentRoom));
							ws.onopen = () => { loadPlugins(); };
							ws.onmessage = (event) => {
								const msg = JSON.parse(event.data);
								if (msg.method === 'superdaw/plugin_params_update') {
									const params = JSON.parse(msg.params.parameters);
									params.forEach(p => {
										const slider = document.getElementById('param-' + p.n);
										if (slider) slider.value = p.v;
									});
									return;
								}
								const state = msg;
								document.getElementById('presence-count').innerText = state.user_count + ' users';

								// Render DAWs
								let dawHtml = '';
								for (const name in state.daws) {
									const d = state.daws[name];
									dawHtml += ' 									<div class="patch-item"> 										<span><strong>' + name.toUpperCase() + '</strong></span> 										<span>' + (d.is_playing ? '▶️ PLAYING' : '⏹️ STOPPED') + '</span> 										<span class="badge">' + d.bpm.toFixed(1) + ' BPM</span> 									</div> 								';
								}
								document.getElementById('daws').innerHTML = dawHtml || '<p style="color: #666">No active DAWs connected.</p>';

								// Render Patches
								let patchHtml = '';
								if (state.patches) {
									state.patches.forEach((p, idx) => {
										patchHtml += ' 									<div class="patch-item"> 										<span>' + p.source_daw + ' (' + p.source_track + ')</span> 										<span class="patch-arrow">➔</span> 										<span>' + p.dest_daw + ' (' + p.dest_track + ')</span> 										<button onclick="removePatch(\'' + p.source_daw + '\', \'' + p.source_track + '\', \'' + p.dest_daw + '\', \'' + p.dest_track + '\')" style="background: none; border: none; color: #ff4444; cursor: pointer;">[X]</button> 									</div> 								';
									});
								}
								document.getElementById('routing').innerHTML = patchHtml || '<p style="color: #666">No active audio patches.</p>';

								// Render Jobs
								let jobHtml = '';
								if (state.jobs) {
									state.jobs.forEach(j => {
										jobHtml += ' 										<div class="patch-item"> 											<span>' + j.prompt + '</span> 											<span class="badge" style="width: 100px; background: #444; position: relative; overflow: hidden;"> 												<div style="background: #00ff88; width: ' + (j.progress*100) + '%%; height: 10px; border-radius: 5px;"></div> 											</span> 											<span>' + j.status + '</span> 										</div> 									';
									});
								}
								document.getElementById('jobs').innerHTML = jobHtml || '<p style="color: #666">No active generation jobs.</p>';

								// Render Blueprint
								let blueprint = 'graph LR\n';
								if (state.patches) {
									state.patches.forEach(p => {
										blueprint += '  ' + p.source_daw + ' --> ' + p.dest_daw + '\n';
									});
								}
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
												timelineHtml += '<div class="track-lane" style="top: ' + top + 'px; position: relative;"><span style="width: 100px; display: inline-block;">' + track.track + '</span>';
												track.clips.forEach(clip => {
													const left = clip.start * 20;
													const width = (clip.end - clip.start) * 20;
													timelineHtml += '<div class="clip-block" style="left: ' + (100+left) + 'px; width: ' + width + 'px;">' + clip.name + '</div>';
												});
												timelineHtml += '</div>';
												top += 40;
											});
										} catch(e) {}
									}
								}
								document.getElementById('timeline').innerHTML = timelineHtml || '<p style="color: #666; padding: 20px;">No arrangement data available.</p>';
							};
						}

						function switchRoom(newRoom) {
							currentRoom = newRoom;
							window.history.pushState({}, '', '?room=' + encodeURIComponent(currentRoom));
							connectWS();
						}

						function playNote(pitch) {
							callMcp('superdaw_write_midi', {
								track_id: '0',
								notes: [{pitch: pitch, velocity: 100, start_beat: 0, duration: 0.5}]
							});
						}

						async function callMcp(name, args) {
							return fetch('/api/call', {
								method: 'POST',
								headers: {'Content-Type': 'application/json'},
								body: JSON.stringify({room_id: currentRoom, name, arguments: args})
							});
						}

						async function loadPlugins() {
							const res = await callMcp('superdaw_list_plugins', {});
							const data = await res.json();
							const plugins = JSON.parse(data.content[0].text);
							let html = '<option>Select a plugin...</option>';
							plugins.forEach(p => {
								html += '<option value="' + p + '">' + p + '</option>';
							});
							document.getElementById('plugin-list').innerHTML = html;
						}

						async function loadPluginParams(pluginName) {
							if (pluginName === "Select a plugin...") return;
							const res = await callMcp('superdaw_get_plugin_params', {plugin_name: pluginName});
							const data = await res.json();
							const meta = JSON.parse(data.content[0].text);
							let html = '<h3>' + pluginName + ' Parameters</h3>';
							meta.parameters.forEach(p => {
								html += '<div style="margin: 5px 0; display: flex; justify-content: space-between; align-items: center;">' +
									'<span>' + p.name + '</span>' +
									'<input type="range" id="param-' + p.name + '" min="0" max="1" step="0.01" style="width: 200px;" onchange="setParam(\'' + pluginName + '\', \'' + p.name + '\', this.value)">' +
								'</div>';
							});
							document.getElementById('plugin-params').innerHTML = html;
						}

						function setParam(plugin, param, val) {
							callMcp('superdaw_set_plugin_parameter', {
								plugin_name: plugin,
								parameter_name: param,
								value: parseFloat(val),
								track_id: '0'
							});
						}

						function addPatch() {
							const srcDaw = document.getElementById('src-daw').value;
							const srcTrack = document.getElementById('src-track').value;
							const dstDaw = document.getElementById('dst-daw').value;
							const dstTrack = document.getElementById('dst-track').value;
							if (srcDaw && srcTrack && dstDaw && dstTrack) {
								callMcp('superdaw_patch_audio', {
									source_daw: srcDaw,
									source_track: srcTrack,
									dest_daw: dstDaw,
									dest_track: dstTrack
								});
							}
						}

						function removePatch(srcDaw, srcTrack, dstDaw, dstTrack) {
							callMcp('superdaw_unpatch_audio', {
								source_daw: srcDaw,
								source_track: srcTrack,
								dest_daw: dstDaw,
								dest_track: dstTrack
							});
						}

						function sendCC() {
							callMcp('superdaw_send_cc', {
								track_id: document.getElementById('cc-track').value,
								controller: parseInt(document.getElementById('cc-num').value),
								value: parseInt(document.getElementById('cc-val').value)
							});
						}

						function sendCustom() {
							let args = {};
							try { args = JSON.parse(document.getElementById('custom-args').value); } catch(e) {}
							callMcp('superdaw_custom_command', {
								command: document.getElementById('custom-cmd').value,
								args: args
							});
						}

						function importGenerative() {
							callMcp('superdaw_import_generative', {
								prompt: document.getElementById('gen-prompt').value,
								target_daw: 'active'
							});
						}

						function generateEuclidean() {
							callMcp('superdaw_generate_euclidean', {
								track_id: document.getElementById('euc-track').value,
								hits: parseInt(document.getElementById('euc-hits').value),
								steps: parseInt(document.getElementById('euc-steps').value),
								pitch: parseInt(document.getElementById('euc-pitch').value)
							});
						}

						connectWS();
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
	return orchestrator, mux
}

func (s *DashboardState) getOrCreateRoom(id string) *RoomState {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room, ok := s.rooms[id]; ok {
		return room
	}
	room := &RoomState{
		ID:      id,
		DAWs:    make(map[string]DAWState),
		Patches: []AudioPatch{},
		clients: make(map[*websocket.Conn]bool),
	}
	s.rooms[id] = room
	return room
}

// Global update methods that default to "default" room or can be updated to take roomID
func (s *DashboardState) UpdateDAW(roomID string, name string, playing bool, bpm float64) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	state := room.DAWs[name]
	state.Name = name
	state.IsPlaying = playing
	state.BPM = bpm
	room.DAWs[name] = state
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) UpdateJobs(roomID string, jobs []interface{}) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	room.Jobs = jobs
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) GetState(roomID string) RoomState {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.RLock()
	defer room.mu.RUnlock()
	return RoomState{
		ID:      room.ID,
		DAWs:    room.DAWs,
		Patches: room.Patches,
		Jobs:    room.Jobs,
	}
}

func (s *DashboardState) SetState(roomID string, state RoomState) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	room.DAWs = state.DAWs
	room.Patches = state.Patches
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) UpdateArrangement(roomID string, name string, arrangement string) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	state := room.DAWs[name]
	state.Arrangement = arrangement
	room.DAWs[name] = state
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) RemovePatch(roomID string, p AudioPatch) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	newPatches := []AudioPatch{}
	for _, patch := range room.Patches {
		if patch.SourceDAW == p.SourceDAW && patch.SourceTrack == p.SourceTrack &&
			patch.DestDAW == p.DestDAW && patch.DestTrack == p.DestTrack {
			continue
		}
		newPatches = append(newPatches, patch)
	}
	room.Patches = newPatches
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) AddPatch(roomID string, p AudioPatch) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	room.Patches = append(room.Patches, p)
	room.mu.Unlock()
	room.broadcast()
}

func (r *RoomState) broadcast() {
	r.mu.RLock()
	defer r.mu.RUnlock()

	msg := struct {
		RoomState
		UserCount int `json:"user_count"`
	}{
		RoomState: *r,
		UserCount: len(r.clients),
	}

	data, _ := json.Marshal(msg)
	for conn := range r.clients {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}
