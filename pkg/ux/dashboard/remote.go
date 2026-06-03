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
						body { font-family: sans-serif; background: #000; color: #fff; text-align: center; padding: 20px; }
						.btn { display: block; width: 100%%; padding: 25px; margin: 10px 0; border: none; border-radius: 12px; font-size: 24px; font-weight: bold; cursor: pointer; }
						.play { background: #00ff88; color: #000; }
						.stop { background: #ff4444; color: #fff; }
						.fader { width: 100%%; height: 60px; margin: 30px 0; -webkit-appearance: none; background: #333; border-radius: 30px; }
						.fader::-webkit-slider-thumb { -webkit-appearance: none; width: 60px; height: 60px; background: #00ff88; border-radius: 50%%; }
						label { font-size: 18px; color: #888; text-transform: uppercase; letter-spacing: 2px; }
					</style>
				</head>
				<body>
					<h1>SuperDAW Remote</h1>
					<button class="btn play" onclick="callTool('superdaw_transport_control', {playing: true})">PLAY</button>
					<button class="btn stop" onclick="callTool('superdaw_transport_control', {playing: false})">STOP</button>

					<div style="margin-top: 40px;">
						<label>Master Volume</label>
						<input type="range" class="fader" min="0" max="1" step="0.01" value="0.8" oninput="callTool('superdaw_set_mixer', {track_id: '0', volume: parseFloat(this.value)})">
					</div>

					<script>
						async function callTool(name, args) {
							// In a real implementation, this would proxy through the MCP server.
							// For this remote, we assume an internal API endpoint is available.
							console.log("Calling tool:", name, args);
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

}
