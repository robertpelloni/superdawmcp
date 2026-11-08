package dashboard

import (
	"fmt"
	"net/http"
)

// RegisterMobileRemote sets up a touch-friendly remote control interface on /remote
func RegisterMobileRemote(mux *http.ServeMux) {
	mux.HandleFunc("/remote", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `
			<!DOCTYPE html>
			<html lang="en">
				<head>
					<meta charset="utf-8">
					<title>SuperDAW Mobile Remote</title>
					<meta name="viewport" content="width=device-width, initial-scale=1, maximum-scale=1, user-scalable=no">
					<style>
						* { box-sizing: border-box; -webkit-tap-highlight-color: transparent; }
						body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; background: #0f172a; color: #f8fafc; margin: 0; padding: 16px; }
						.header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 20px; }
						.title { font-size: 20px; font-weight: 800; color: #38bdf8; letter-spacing: -0.5px; }
						.status-badge { font-size: 11px; background: rgba(16, 185, 129, 0.2); color: #10b981; border: 1px solid #10b981; padding: 3px 8px; border-radius: 999px; }
						.card { background: #1e293b; border: 1px solid #334155; border-radius: 12px; padding: 16px; margin-bottom: 16px; }
						.card-title { font-size: 13px; font-weight: 700; text-transform: uppercase; color: #94a3b8; letter-spacing: 1px; margin-bottom: 12px; }
						.transport-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
						.btn { width: 100%; padding: 16px; border: none; border-radius: 8px; font-size: 16px; font-weight: 700; cursor: pointer; text-align: center; }
						.play { background: #10b981; color: #000; }
						.stop { background: #ef4444; color: #fff; }
						.macro-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 8px; margin-bottom: 12px; }
						.btn-macro { background: #334155; color: #e2e8f0; border: 1px solid #475569; padding: 12px 8px; border-radius: 6px; font-size: 12px; font-weight: 600; cursor: pointer; }
						.select-box { width: 100%; padding: 12px; background: #0f172a; border: 1px solid #475569; color: #f8fafc; border-radius: 6px; font-size: 14px; margin-bottom: 12px; outline: none; }
						.param-row { background: rgba(0,0,0,0.25); border: 1px solid #334155; border-radius: 8px; padding: 12px; margin-bottom: 10px; }
						.param-header { display: flex; justify-content: space-between; font-size: 13px; font-weight: 600; margin-bottom: 8px; }
						.fader { width: 100%; height: 36px; -webkit-appearance: none; background: #0f172a; border-radius: 18px; outline: none; }
						.fader::-webkit-slider-thumb { -webkit-appearance: none; width: 36px; height: 36px; background: #38bdf8; border-radius: 50%; cursor: pointer; }
						.badge { font-size: 10px; padding: 2px 6px; border-radius: 4px; background: #334155; color: #94a3b8; }
					</style>
				</head>
				<body>
					<div class="header">
						<div class="title">⚡ SuperDAW Remote</div>
						<div class="status-badge" id="conn-status">Online</div>
					</div>

					<!-- Transport Section -->
					<div class="card">
						<div class="card-title">Transport Control</div>
						<div class="transport-grid">
							<button class="btn play" onclick="callTool('superdaw_transport_control', {playing: true})">▶ PLAY</button>
							<button class="btn stop" onclick="callTool('superdaw_transport_control', {playing: false})">■ STOP</button>
						</div>
					</div>

					<!-- Mobile Plugin Inspector Section -->
					<div class="card">
						<div class="card-title">VST3 Touch Inspector</div>
						<select id="mobile-plugin-list" class="select-box" onchange="loadMobilePlugin(this.value)">
							<option value="">Select plugin...</option>
						</select>

						<div class="card-title" style="margin-top: 8px;">Sound Design Macros</div>
						<div class="macro-grid">
							<button class="btn-macro" onclick="triggerMobileMacro('brighten')">✨ Brighten</button>
							<button class="btn-macro" onclick="triggerMobileMacro('warmth')">🔥 Warmth</button>
							<button class="btn-macro" onclick="triggerMobileMacro('space')">🌌 Space Out</button>
							<button class="btn-macro" onclick="triggerMobileMacro('punch')">🥊 Punch</button>
						</div>

						<div id="mobile-param-list">
							<p style="color: #64748b; font-size: 13px; text-align: center;">Select a plugin to control parameters.</p>
						</div>
					</div>

					<script>
						let cachedPlugins = {};
						let activePluginName = "";

						async function callTool(name, args) {
							try {
								await fetch('/api/call', {
									method: 'POST',
									headers: {'Content-Type': 'application/json'},
									body: JSON.stringify({name, arguments: args})
								});
							} catch(e) {
								console.error("Tool dispatch failed:", e);
							}
						}

						async function initPlugins() {
							try {
								const res = await fetch('/api/plugins');
								cachedPlugins = await res.json();
								const sel = document.getElementById('mobile-plugin-list');
								sel.innerHTML = '<option value="">Select plugin...</option>';
								for (const name in cachedPlugins) {
									sel.innerHTML += '<option value="' + name + '">' + name + '</option>';
								}
								// Auto-select first plugin if available
								const keys = Object.keys(cachedPlugins);
								if (keys.length > 0) {
									sel.value = keys[0];
									loadMobilePlugin(keys[0]);
								}
							} catch(e) {
								console.error("Plugin fetch failed:", e);
							}
						}

						function loadMobilePlugin(name) {
							activePluginName = name;
							if (!name || !cachedPlugins[name]) {
								document.getElementById('mobile-param-list').innerHTML = '<p style="color: #64748b; font-size: 13px; text-align: center;">Select a plugin to control parameters.</p>';
								return;
							}
							const p = cachedPlugins[name];
							let html = '';
							(p.parameters || []).forEach(param => {
								const val = param.default_value !== undefined ? param.default_value : 0.5;
								const unit = param.unit ? ' ' + param.unit : '';
								html += ' \
									<div class="param-row"> \
										<div class="param-header"> \
											<span>' + param.name + '</span> \
											<span class="badge" id="val-' + param.index + '">' + val.toFixed(2) + unit + '</span> \
										</div> \
										<input type="range" class="fader" min="' + (param.min_value !== undefined ? param.min_value : 0) + '" max="' + (param.max_value !== undefined ? param.max_value : 1) + '" step="' + (param.step || 0.01) + '" value="' + val + '" oninput="syncMobileParam(\'' + name + '\', \'' + param.name + '\', this.value, ' + param.index + ', \'' + (param.unit || '') + '\')"> \
									</div> \
								';
							});
							document.getElementById('mobile-param-list').innerHTML = html;
						}

						function syncMobileParam(pluginName, paramName, value, index, unit) {
							const num = parseFloat(value);
							const badge = document.getElementById('val-' + index);
							if (badge) badge.innerText = num.toFixed(2) + (unit ? ' ' + unit : '');
							callTool('superdaw_set_plugin_parameter', {
								daw: 'ableton',
								track_id: '0',
								plugin_name: pluginName,
								parameter_name: paramName,
								value: num
							});
						}

						function triggerMobileMacro(macroType) {
							if (!activePluginName || !cachedPlugins[activePluginName]) return;
							const p = cachedPlugins[activePluginName];
							(p.parameters || []).forEach(param => {
								let newVal = null;
								const pName = param.name.toLowerCase();
								if (macroType === 'brighten' && (pName.includes('cutoff') || pName.includes('high cut'))) {
									newVal = Math.min((param.max_value || 1.0), 0.85);
								} else if (macroType === 'warmth' && (pName.includes('drive') || pName.includes('damping'))) {
									newVal = 0.65;
								} else if (macroType === 'space' && (pName.includes('reverb') || pName.includes('wet') || pName.includes('delay'))) {
									newVal = 0.70;
								} else if (macroType === 'punch' && (pName.includes('attack'))) {
									newVal = 0.01;
								}
								if (newVal !== null) {
									syncMobileParam(activePluginName, param.name, newVal, param.index, param.unit || '');
								}
							});
							loadMobilePlugin(activePluginName);
						}

						window.onload = initPlugins;
					</script>
				</body>
			</html>
		`)
	})
}
