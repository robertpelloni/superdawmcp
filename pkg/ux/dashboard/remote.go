package dashboard

import (
	"fmt"
	"net/http"
)

// RegisterMobileRemote sets up a touch-friendly remote control interface on /remote
func RegisterMobileRemote(mux *http.ServeMux) {
	mux.HandleFunc("/remote", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<html>
				<head>
					<title>SuperDAW Mobile Remote</title>
					<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
					<style>
						body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; background: #000; color: #fff; text-align: center; padding: 10px; margin: 0; overflow-x: hidden; }
						.header { padding: 20px 0; background: #111; border-bottom: 1px solid #333; margin-bottom: 20px; display: flex; flex-direction: column; align-items: center; }
						.section { background: #1a1a1a; margin: 10px; padding: 15px; border-radius: 12px; border: 1px solid #333; }
						h2 { font-size: 14px; color: #888; text-transform: uppercase; letter-spacing: 2px; margin-top: 0; }
						.btn { display: inline-block; width: 45%%; padding: 25px 0; margin: 5px; border: none; border-radius: 12px; font-size: 20px; font-weight: bold; cursor: pointer; transition: transform 0.1s; }
						.btn:active { transform: scale(0.95); }
						.play { background: #00ff88; color: #000; }
						.stop { background: #ff4444; color: #fff; }
						.scene-btn { width: 30%%; padding: 15px 0; background: #333; color: #00bcd4; font-size: 16px; }
						.fader-container { margin: 20px 0; }
						.fader { width: 100%%; height: 50px; -webkit-appearance: none; background: #222; border-radius: 25px; outline: none; }
						.fader::-webkit-slider-thumb { -webkit-appearance: none; width: 50px; height: 50px; background: #00ff88; border-radius: 50%%; box-shadow: 0 0 10px rgba(0,255,136,0.5); }
						.pan-fader::-webkit-slider-thumb { background: #00bcd4; }
						select { width: 100%%; padding: 15px; background: #222; color: #fff; border: 1px solid #444; border-radius: 8px; font-size: 16px; margin-bottom: 10px; }
						.room-badge { background: #00ff88; color: #000; padding: 5px 15px; border-radius: 20px; font-size: 12px; margin-top: 5px; font-weight: bold; }
					</style>
				</head>
				<body>
					<div class="header">
						<h1 style="margin: 0; font-size: 24px; color: #00ff88;">SuperDAW Remote</h1>
						<div id="room-display" class="room-badge">Room: default</div>
					</div>

					<div class="section">
						<h2>Transport</h2>
						<button class="btn play" onclick="callTool('superdaw_transport_control', {playing: true})">PLAY</button>
						<button class="btn stop" onclick="callTool('superdaw_transport_control', {playing: false})">STOP</button>
					</div>

					<div class="section">
						<h2>Scene Launcher</h2>
						<div style="display: flex; flex-wrap: wrap; justify-content: center;">
							<button class="btn scene-btn" onclick="callTool('superdaw_fire_scene', {scene_index: 0})">SCENE 1</button>
							<button class="btn scene-btn" onclick="callTool('superdaw_fire_scene', {scene_index: 1})">SCENE 2</button>
							<button class="btn scene-btn" onclick="callTool('superdaw_fire_scene', {scene_index: 2})">SCENE 3</button>
							<button class="btn scene-btn" onclick="callTool('superdaw_fire_scene', {scene_index: 3})">SCENE 4</button>
							<button class="btn scene-btn" onclick="callTool('superdaw_fire_scene', {scene_index: 4})">SCENE 5</button>
							<button class="btn scene-btn" onclick="callTool('superdaw_fire_scene', {scene_index: 5})">SCENE 6</button>
						</div>
					</div>

					<div class="section">
						<h2>Mixer Control</h2>
						<select id="track-id" onchange="updateFaderLabel()">
							<option value="0">Master</option>
						</select>
						<div class="fader-container">
							<div id="vol-label" style="margin-bottom: 10px; color: #00ff88;">Master Volume: 0.80</div>
							<input type="range" id="vol-fader" class="fader" min="0" max="1" step="0.01" value="0.8" oninput="handleVol(this.value)">
						</div>
						<div class="fader-container">
							<div id="pan-label" style="margin-bottom: 10px; color: #00bcd4;">Master Pan: 0.00</div>
							<input type="range" id="pan-fader" class="fader pan-fader" min="-1" max="1" step="0.01" value="0" oninput="handlePan(this.value)">
						</div>
					</div>

					<script>
						let currentRoom = new URLSearchParams(window.location.search).get('room') || 'default';
						document.getElementById('room-display').innerText = 'Room: ' + currentRoom;

						let ws;
						function connectWS() {
							if (ws) ws.close();
							ws = new WebSocket('ws://' + window.location.host + '/ws?room=' + currentRoom);
							ws.onmessage = (event) => {
								const state = JSON.parse(event.data);
								console.log("State updated:", state);

								const trackSelector = document.getElementById('track-id');
								let currentVal = trackSelector.value;
								const savedTrack = localStorage.getItem('superdaw_last_track');
								if (!currentVal && savedTrack) currentVal = savedTrack;

								let options = '<option value="0">Master</option>';
								for (const name in state.daws) {
									const d = state.daws[name];
									if (d.arrangement) {
										try {
											const arrangement = JSON.parse(d.arrangement);
											arrangement.forEach((track, idx) => {
												options += '<option value="' + (idx+1) + '">' + track.track + '</option>';
											});
										} catch(e) {}
									}
								}
								trackSelector.innerHTML = options;
								trackSelector.value = currentVal;
								updateFaderLabel();
							};
						}

						// Restore saved track selection
						window.onload = () => {
							const savedTrack = localStorage.getItem('superdaw_last_track');
							if (savedTrack) {
								document.getElementById('track-id').value = savedTrack;
								updateFaderLabel();
							}
							connectWS();
						};

						function updateFaderLabel() {
							const v = document.getElementById('vol-fader').value;
							const p = document.getElementById('pan-fader').value;
							const trackSelector = document.getElementById('track-id');
							localStorage.setItem('superdaw_last_track', trackSelector.value);
							const name = trackSelector.options[trackSelector.selectedIndex]?.text || 'Master';
							document.getElementById('vol-label').innerText = name + ' Volume: ' + parseFloat(v).toFixed(2);
							document.getElementById('pan-label').innerText = name + ' Pan: ' + parseFloat(p).toFixed(2);
						}

						function handleVol(val) {
							updateFaderLabel();
							const track = document.getElementById('track-id').value;
							callTool('superdaw_set_mixer', {track_id: track, volume: parseFloat(val)});
						}

						function handlePan(val) {
							updateFaderLabel();
							const track = document.getElementById('track-id').value;
							callTool('superdaw_set_mixer', {track_id: track, pan: parseFloat(val)});
						}

						async function callTool(name, args) {
							console.log("Calling tool in room " + currentRoom + ":", name, args);
							fetch('/api/call', {
								method: 'POST',
								headers: {'Content-Type': 'application/json'},
								body: JSON.stringify({room_id: currentRoom, name, arguments: args})
							});
						}
					</script>
				</body>
			</html>
		`)
	})

}
