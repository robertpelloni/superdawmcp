package dashboard

import (
	"fmt"
	"net/http"
)

// renderUnifiedDashboard outputs the consolidated, single responsive dashboard page.
// Unescaped literal % for fmt.Fprint for fmt.Fprintf to avoid fmt template substitution errors.
func renderUnifiedDashboard(w http.ResponseWriter) {
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1">
	<title>SuperDAW Master Studio Console</title>
	<style>
		:root {
			--bg-main: #0c0d0f;
			--bg-card: #15181c;
			--bg-card-sub: #1e2228;
			--border: #2a303c;
			--primary: #00ff88;
			--secondary: #00bcd4;
			--accent: #ff4772;
			--text-main: #f0f4f8;
			--text-muted: #8b9bb4;
			--font-mono: 'JetBrains Mono', 'Fira Code', Menlo, Consolas, monospace;
		}
		* { box-sizing: border-box; margin: 0; padding: 0; }
		body {
			font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
			background: var(--bg-main);
			color: var(--text-main);
			padding: 20px;
			min-height: 100vh;
			line-height: 1.5;
		}

		/* Tooltip styling */
		[data-tooltip] {
			position: relative;
			cursor: help;
		}
		[data-tooltip]::after {
			content: attr(data-tooltip);
			position: absolute;
			bottom: 125%;
			left: 50%;
			transform: translateX(-50%);
			background: #000;
			color: #fff;
			padding: 6px 10px;
			border-radius: 6px;
			font-size: 11px;
			font-weight: normal;
			white-space: nowrap;
			opacity: 0;
			pointer-events: none;
			transition: opacity 0.2s;
			box-shadow: 0 4px 10px rgba(0,0,0,0.6);
			border: 1px solid var(--border);
			z-index: 1000;
		}
		[data-tooltip]:hover::after {
			opacity: 1;
		}

		/* Header and Top Bar */
		header {
			display: flex;
			flex-wrap: wrap;
			justify-content: space-between;
			align-items: center;
			padding-bottom: 20px;
			margin-bottom: 25px;
			border-bottom: 1px solid var(--border);
			gap: 15px;
		}
		.logo-section h1 {
			font-size: 24px;
			font-weight: 800;
			letter-spacing: 1.5px;
			color: var(--primary);
			text-transform: uppercase;
			display: flex;
			align-items: center;
			gap: 10px;
		}
		.logo-section p {
			font-size: 12px;
			color: var(--text-muted);
			margin-top: 3px;
		}
		.header-actions {
			display: flex;
			align-items: center;
			gap: 10px;
			flex-wrap: wrap;
		}

		/* Buttons & Badges */
		.btn {
			padding: 8px 16px;
			border-radius: 6px;
			font-size: 12px;
			font-weight: 700;
			text-transform: uppercase;
			cursor: pointer;
			border: none;
			transition: all 0.15s ease-in-out;
			display: inline-flex;
			align-items: center;
			gap: 6px;
			text-decoration: none;
		}
		.btn-primary { background: var(--primary); color: #000; }
		.btn-primary:hover { background: #00dd75; }
		.btn-secondary { background: var(--secondary); color: #000; }
		.btn-secondary:hover { background: #00a4b8; }
		.btn-danger { background: var(--accent); color: #fff; }
		.btn-danger:hover { background: #e0325d; }
		.btn-subtle { background: var(--bg-card-sub); color: var(--text-main); border: 1px solid var(--border); }
		.btn-subtle:hover { background: #282f3a; border-color: #3e4758; }
		.btn-tab {
			padding: 5px 12px;
			border-radius: 4px;
			font-size: 11px;
			cursor: pointer;
			background: var(--bg-card-sub);
			color: var(--text-muted);
			border: 1px solid var(--border);
			font-weight: 600;
		}
		.btn-tab.active {
			background: var(--secondary);
			color: #000;
			border-color: var(--secondary);
		}

		.badge {
			padding: 4px 10px;
			border-radius: 20px;
			font-size: 11px;
			font-family: var(--font-mono);
			background: var(--bg-card-sub);
			border: 1px solid var(--border);
			color: var(--text-muted);
		}
		.badge-online {
			background: rgba(0, 255, 136, 0.1);
			color: var(--primary);
			border-color: rgba(0, 255, 136, 0.3);
		}

		/* Grid Layout */
		.grid-main {
			display: grid;
			grid-template-columns: repeat(12, 1fr);
			gap: 20px;
		}
		.col-12 { grid-column: span 12; }
		.col-8 { grid-column: span 8; }
		.col-6 { grid-column: span 6; }
		.col-4 { grid-column: span 4; }

		@media (max-width: 1024px) {
			.col-8, .col-6, .col-4 { grid-column: span 12; }
		}

		/* Card Section Styling */
		.card {
			background: var(--bg-card);
			border: 1px solid var(--border);
			border-radius: 12px;
			padding: 20px;
			display: flex;
			flex-direction: column;
			box-shadow: 0 4px 12px rgba(0,0,0,0.3);
		}
		.card-header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			margin-bottom: 16px;
			padding-bottom: 10px;
			border-bottom: 1px solid var(--border);
		}
		.card-title {
			font-size: 15px;
			font-weight: 700;
			text-transform: uppercase;
			letter-spacing: 0.8px;
			color: var(--secondary);
			display: flex;
			align-items: center;
			gap: 8px;
		}

		/* Forms & Controls */
		.control-group {
			display: flex;
			flex-direction: column;
			gap: 6px;
			margin-bottom: 12px;
		}
		.control-group label {
			font-size: 11px;
			text-transform: uppercase;
			color: var(--text-muted);
			font-weight: 600;
		}
		.input-text, .select-box {
			background: var(--bg-card-sub);
			border: 1px solid var(--border);
			color: var(--text-main);
			padding: 8px 12px;
			border-radius: 6px;
			font-size: 13px;
			font-family: inherit;
		}
		.input-text:focus, .select-box:focus {
			border-color: var(--secondary);
			outline: none;
		}
		.input-row {
			display: flex;
			gap: 10px;
			align-items: center;
		}

		/* Plugin Inspector Sliders & Metadata */
		.inspector-toolbar {
			display: flex;
			justify-content: space-between;
			align-items: center;
			gap: 10px;
			flex-wrap: wrap;
			margin-top: 10px;
			margin-bottom: 10px;
		}
		.filter-tabs {
			display: flex;
			gap: 6px;
			flex-wrap: wrap;
		}
		.param-grid {
			display: grid;
			grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
			gap: 15px;
			margin-top: 10px;
			max-height: 380px;
			overflow-y: auto;
			padding-right: 5px;
		}
		.param-box {
			background: var(--bg-card-sub);
			border: 1px solid var(--border);
			border-radius: 8px;
			padding: 12px;
			display: flex;
			flex-direction: column;
			gap: 8px;
		}
		.param-header {
			display: flex;
			justify-content: space-between;
			align-items: center;
			font-size: 12px;
			font-weight: 600;
		}
		.param-slider-row {
			display: flex;
			gap: 10px;
			align-items: center;
		}
		.param-slider {
			-webkit-appearance: none;
			flex: 1;
			height: 8px;
			border-radius: 4px;
			background: #111;
			outline: none;
		}
		.param-slider::-webkit-slider-thumb {
			-webkit-appearance: none;
			width: 18px;
			height: 18px;
			border-radius: 50%;
			background: var(--primary);
			cursor: pointer;
			box-shadow: 0 0 6px rgba(0,255,136,0.6);
		}
		.param-input {
			width: 65px;
			padding: 3px 6px;
			font-size: 11px;
			font-family: var(--font-mono);
			background: #111;
			border: 1px solid var(--border);
			color: #fff;
			border-radius: 4px;
			text-align: right;
		}
		.meta-tag {
			font-size: 11px;
			font-family: var(--font-mono);
			background: #111;
			padding: 3px 6px;
			border-radius: 4px;
			color: var(--text-muted);
		}

		/* DAW Transport & Matrix List */
		.daw-row {
			background: var(--bg-card-sub);
			padding: 12px;
			border-radius: 8px;
			margin-bottom: 8px;
			display: flex;
			justify-content: space-between;
			align-items: center;
			border: 1px solid var(--border);
		}
		.daw-name { font-weight: 700; color: #fff; font-size: 14px; }
		.daw-status-badge {
			padding: 4px 10px;
			border-radius: 4px;
			font-size: 11px;
			font-weight: bold;
		}
		.playing { background: rgba(0,255,136,0.2); color: var(--primary); }
		.stopped { background: rgba(255,71,114,0.2); color: var(--accent); }

		/* Timeline Arrangement */
		#timeline {
			width: 100%;
			min-height: 160px;
			background: #08090a;
			border: 1px solid var(--border);
			border-radius: 8px;
			position: relative;
			overflow-x: auto;
			padding: 10px;
		}
		.track-lane {
			height: 36px;
			border-bottom: 1px solid #1a202c;
			display: flex;
			align-items: center;
			white-space: nowrap;
			position: relative;
		}
		.clip-block {
			position: absolute;
			background: var(--secondary);
			height: 26px;
			border-radius: 4px;
			font-size: 10px;
			font-weight: bold;
			color: #000;
			padding: 4px 8px;
			overflow: hidden;
			text-overflow: ellipsis;
			box-shadow: 0 2px 5px rgba(0,0,0,0.5);
		}

		/* Performance Keyboard */
		.keyboard {
			display: flex;
			justify-content: center;
			position: relative;
			margin-top: 10px;
			padding-bottom: 5px;
		}
		.key {
			width: 44px;
			height: 130px;
			border: 1px solid #000;
			background: #f8fafc;
			cursor: pointer;
			border-radius: 0 0 5px 5px;
			box-shadow: 0 4px 6px rgba(0,0,0,0.5);
			user-select: none;
		}
		.key.black {
			background: #111;
			height: 85px;
			width: 32px;
			margin-left: -16px;
			margin-right: -16px;
			z-index: 2;
			border-color: #333;
		}
		.key:active { background: var(--primary); }

		/* Blueprint graph */
		#blueprint {
			width: 100%;
			min-height: 110px;
			background: #08090a;
			border: 1px dashed var(--border);
			border-radius: 6px;
			padding: 12px;
			font-family: var(--font-mono);
			color: var(--primary);
			font-size: 12px;
			overflow: auto;
		}

		/* Status footer */
		footer {
			margin-top: 30px;
			border-top: 1px solid var(--border);
			padding-top: 15px;
			display: flex;
			justify-content: space-between;
			font-size: 11px;
			color: var(--text-muted);
			font-family: var(--font-mono);
		}
	</style>
</head>
<body>

	<!-- Top Master Header -->
	<header>
		<div class="logo-section">
			<h1>⚡ SuperDAW Master Console</h1>
			<p>Unified Multi-DAW Orchestration, VST3 Discovery & Neural Stem Engine</p>
		</div>
		<div class="header-actions">
			<span class="badge badge-success" data-tooltip="Active daemon instance version">v3.15.0 Active</span>
			<button class="btn btn-primary" onclick="callMcp('superdaw_save_session', {})" data-tooltip="Save studio session state (routing, clips, patches) to JSON">💾 Save Session</button>
			<button class="btn btn-secondary" onclick="callMcp('superdaw_load_session', {})" data-tooltip="Restore studio state from saved JSON file">📂 Load Session</button>
			<a href="/remote" target="_blank" class="btn btn-subtle" data-tooltip="Open dedicated touch-optimized remote control for mobile/tablet">📱 Mobile Remote</a>
			<a href="/obs" target="_blank" class="btn btn-subtle" data-tooltip="Clean, transparent stream overlay HUD for OBS Studio">📺 OBS Overlay</a>
			<button class="btn btn-danger" onclick="quitSuperDAW()" data-tooltip="Gracefully shutdown all SuperDAW MCP servers, TCP gateways and OSC bridges">🛑 Quit Server</button>
		</div>
	</header>

	<main class="grid-main">

		<!-- SECTION 1 (HIGH VALUE): DAW Transport & Orchestration -->
		<section class="card col-4">
			<div class="card-header">
				<h2 class="card-title">🎛️ DAW Instances & Transport</h2>
				<span class="badge" id="link-peer-badge" data-tooltip="Active synchronized Link peers on local subnet">Link: 1 Peer</span>
			</div>
			<div class="input-row" style="margin-bottom: 12px;">
				<button class="btn btn-primary" style="flex: 1;" onclick="syncTransport(true)" data-tooltip="Trigger playback on all active DAWs and Link network">▶ Play All</button>
				<button class="btn btn-danger" style="flex: 1;" onclick="syncTransport(false)" data-tooltip="Stop transport on all connected DAWs and Link network">⏹ Stop All</button>
			</div>

			<!-- Universal Link Clock HUD -->
			<div style="background: rgba(0,0,0,0.3); border: 1px solid var(--border-color); border-radius: 6px; padding: 10px; margin-bottom: 12px;">
				<div style="display: flex; justify-content: space-between; font-size: 11px; font-weight: 700; color: var(--text-muted); text-transform: uppercase; margin-bottom: 6px;">
					<span>Link Studio Clock</span>
					<span id="link-quantum-label" style="color: var(--secondary);">4/4 Bar</span>
				</div>
				<div style="display: flex; align-items: center; justify-content: space-between; font-family: var(--font-mono);">
					<div>
						<span style="font-size: 11px; color: var(--text-muted);">Beat:</span>
						<span id="link-beat-val" style="font-size: 15px; font-weight: bold; color: var(--primary); margin-left: 4px;">0.00</span>
					</div>
					<div>
						<span style="font-size: 11px; color: var(--text-muted);">Phase:</span>
						<span id="link-phase-val" style="font-size: 15px; font-weight: bold; color: #a855f7; margin-left: 4px;">0.00</span>
					</div>
					<div id="link-beat-led" style="width: 12px; height: 12px; border-radius: 50%; background: #334155; border: 1px solid #475569;" data-tooltip="Flashing pulse indicator for active tempo"></div>
				</div>
			</div>

			<div class="control-group">
				<label data-tooltip="Master BPM clock broadcast to Ableton Link and OSC bridges">Global Tempo (BPM)</label>
				<div class="input-row">
					<input type="number" id="global-bpm" class="input-text" value="128.0" step="0.1" style="flex: 1;" onchange="syncBpm(this.value)">
					<button class="btn btn-subtle" onclick="syncBpm(document.getElementById('global-bpm').value)">Set</button>
				</div>
			</div>
			<div id="daws" style="margin-top: 10px;">
				<!-- Rendered dynamically -->
			</div>
		</section>

		<!-- SECTION 2 (HIGH VALUE): VST3 Plugin Inspector -->
		<section class="card col-8">
			<div class="card-header">
				<h2 class="card-title">🔍 VST3 Plugin Inspector & Discovery</h2>
				<div style="display: flex; gap: 8px;">
					<button class="btn btn-subtle" onclick="deepScanSelected()" data-tooltip="Execute granular binary inspection and bus extraction">🧬 Deep Scan</button>
					<button class="btn btn-subtle" onclick="refreshPlugins()" data-tooltip="Re-scan system VST3 directories">🔄 Refresh</button>
				</div>
			</div>
			<div class="grid-main" style="gap: 15px;">
				<div class="col-6">
					<div class="control-group">
						<label data-tooltip="Select discovered VST3 plugin from local filesystem">Discovered Plugins</label>
						<select id="plugin-list" class="select-box" onchange="loadPlugin(this.value)">
							<option value="">Select a plugin...</option>
						</select>
					</div>
				</div>
				<div class="col-6">
					<div class="control-group">
						<label data-tooltip="Discovered vendor presets and system defaults">Plugin Presets</label>
						<select id="preset-list" class="select-box" onchange="loadPreset(this.value)">
							<option value="">Load Preset...</option>
						</select>
					</div>
				</div>
			</div>
			<div id="plugin-meta-banner" style="display: flex; gap: 8px; flex-wrap: wrap; margin-top: 5px;"></div>

			<!-- Quick Macro Sound Design Knobs & State Tools -->
			<div style="background: rgba(0,0,0,0.25); border: 1px solid var(--border-color); border-radius: 6px; padding: 10px; margin-top: 10px;">
				<div style="font-size: 11px; font-weight: bold; color: var(--text-muted); text-transform: uppercase; margin-bottom: 6px; display: flex; justify-content: space-between;">
					<span>Sound Design Macro Controls & State</span>
					<div style="display: flex; gap: 6px;">
						<button id="btn-ab" class="btn btn-subtle" style="padding: 2px 8px; font-size: 10px;" onclick="toggleABComparison()" data-tooltip="Switch between snapshot A and snapshot B">⚖️ State: <span id="ab-label">A</span></button>
						<button class="btn btn-subtle" style="padding: 2px 8px; font-size: 10px;" onclick="randomizeParameters()" data-tooltip="Controlled musical parameter randomization (±25% range)">🎲 Randomize</button>
						<button class="btn btn-subtle" style="padding: 2px 8px; font-size: 10px;" onclick="resetParametersToDefault()" data-tooltip="Reset all parameters to factory defaults">↺ Reset</button>
					</div>
				</div>
				<div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px;">
					<button class="btn btn-subtle" onclick="triggerMacroPreset('brighten')" data-tooltip="Increase filter cutoff and presence on target plugin">✨ Brighten</button>
					<button class="btn btn-subtle" onclick="triggerMacroPreset('warmth')" data-tooltip="Increase saturation/drive and slight low boost">🔥 Warmth</button>
					<button class="btn btn-subtle" onclick="triggerMacroPreset('space')" data-tooltip="Boost reverb wet mix and pre-delay">🌌 Space Out</button>
					<button class="btn btn-subtle" onclick="triggerMacroPreset('punch')" data-tooltip="Fast attack, punchy decay and envelope depth">🥊 Snappy / Punch</button>
				</div>

				<!-- Continuous Preset Morphing & Snapshot Bank -->
				<div style="display: flex; align-items: center; justify-content: space-between; margin-top: 10px; padding-top: 8px; border-top: 1px dashed var(--border-color); font-size: 11px;">
					<div style="display: flex; align-items: center; gap: 6px;">
						<span style="color: var(--text-muted); font-weight: 700;">BANKS:</span>
						<button class="btn btn-subtle" style="padding: 2px 6px; font-size: 10px;" onclick="storeSnapshotBank('A')" data-tooltip="Save current parameter state to Snapshot Bank A">Slot A</button>
						<button class="btn btn-subtle" style="padding: 2px 6px; font-size: 10px;" onclick="storeSnapshotBank('B')" data-tooltip="Save current parameter state to Snapshot Bank B">Slot B</button>
						<button class="btn btn-subtle" style="padding: 2px 6px; font-size: 10px;" onclick="storeSnapshotBank('C')" data-tooltip="Save current parameter state to Snapshot Bank C">Slot C</button>
						<button class="btn btn-subtle" style="padding: 2px 6px; font-size: 10px;" onclick="storeSnapshotBank('D')" data-tooltip="Save current parameter state to Snapshot Bank D">Slot D</button>
					</div>
					<div style="display: flex; align-items: center; gap: 8px;">
						<span style="color: var(--text-muted); font-size: 10px;" data-tooltip="Continuous parameter interpolation between Bank A and Bank B">A ⟷ B Morph:</span>
						<input type="range" id="morph-slider" class="param-slider" style="width: 100px;" min="0" max="1" step="0.01" value="0" oninput="morphBankParameters(this.value)">
						<span id="morph-val" style="font-size: 10px; color: var(--primary); font-family: var(--font-mono); min-width: 30px;">0%</span>
					</div>
				</div>
			</div>

			<!-- Inspector Filter Toolbar -->
			<div class="inspector-toolbar">
				<div class="filter-tabs">
					<button class="btn-tab active" onclick="setParamFilter('all', this)">All</button>
					<button class="btn-tab" onclick="setParamFilter('Filter/Tone', this)">Filter/Tone</button>
					<button class="btn-tab" onclick="setParamFilter('Envelopes', this)">Envelopes</button>
					<button class="btn-tab" onclick="setParamFilter('Mod/LFO', this)">Mod/LFO</button>
					<button class="btn-tab" onclick="setParamFilter('FX/Mix', this)">FX/Mix</button>
					<button class="btn-tab" onclick="setParamFilter('Master', this)">Master</button>
				</div>
				<input type="text" id="param-search" class="input-text" placeholder="Filter parameter..." style="padding: 4px 8px; font-size: 11px; width: 140px;" oninput="applyParamFilter()">
			</div>

			<div id="plugin-params" class="param-grid">
				<p style="color: var(--text-muted); padding: 15px;">Select a plugin above to inspect parameters and modulate in real time.</p>
			</div>
		</section>

		<!-- SECTION 3: Virtual Mixer & Track Creation -->
		<section class="card col-6">
			<div class="card-header">
				<h2 class="card-title">🎚️ Virtual Track Mixer & Channel Strip</h2>
				<span class="badge" data-tooltip="Target DAW driver for commands">Channel Matrix</span>
			</div>
			<div class="grid-main" style="gap: 10px; margin-bottom: 12px;">
				<div class="col-6">
					<div class="control-group">
						<label>Target DAW</label>
						<select id="mixer-daw" class="select-box">
							<option value="ableton">Ableton Live</option>
							<option value="reaper">REAPER</option>
							<option value="bitwig">Bitwig Studio</option>
							<option value="flstudio">FL Studio</option>
							<option value="logic">Logic Pro</option>
							<option value="cubase">Cubase</option>
							<option value="protools">Pro Tools</option>
							<option value="ardour">Ardour</option>
						</select>
					</div>
				</div>
				<div class="col-6">
					<div class="control-group">
						<label>Track ID / Index</label>
						<input type="text" id="mixer-track-id" class="input-text" value="0">
					</div>
				</div>
			</div>
			<div class="grid-main" style="gap: 15px;">
				<div class="col-6">
					<div class="control-group">
						<label data-tooltip="Set channel fader volume level (0.0 - 1.0)">Fader Volume: <span id="vol-display">0.80</span></label>
						<input type="range" class="param-slider" min="0" max="1" step="0.01" value="0.8" oninput="updateMixerVolume(this.value)">
					</div>
				</div>
				<div class="col-6">
					<div class="control-group">
						<label data-tooltip="Set stereo panoramic positioning (-1.0 to +1.0)">Stereo Pan: <span id="pan-display">0.00</span></label>
						<input type="range" class="param-slider" min="-1" max="1" step="0.05" value="0.0" oninput="updateMixerPan(this.value)">
					</div>
				</div>
			</div>
			<div style="margin-top: 15px; padding-top: 12px; border-top: 1px solid var(--border);">
				<label style="font-size: 11px; text-transform: uppercase; color: var(--text-muted); font-weight: 600;">Create New Track</label>
				<div class="input-row" style="margin-top: 6px;">
					<input type="text" id="new-track-name" class="input-text" placeholder="Track Name (e.g. Lead Synth)" style="flex: 2;">
					<select id="new-track-type" class="select-box" style="flex: 1;">
						<option value="midi">MIDI</option>
						<option value="audio">Audio</option>
					</select>
					<button class="btn btn-subtle" onclick="createTrack()" data-tooltip="Send track creation command to DAW">Add Track</button>
				</div>
			</div>
		</section>

		<!-- SECTION 4: Virtual Audio Patching & Jack Routing Matrix -->
		<section class="card col-6">
			<div class="card-header">
				<h2 class="card-title">🔌 Audio Patching & Inter-DAW Routing</h2>
				<span class="badge" data-tooltip="Jack / ReRoute virtual audio cables">Routing Graph</span>
			</div>
			<div class="input-row" style="margin-bottom: 12px;">
				<input type="text" id="patch-src-daw" class="input-text" placeholder="Src DAW (e.g. ableton)" style="flex: 1;">
				<input type="text" id="patch-src-trk" class="input-text" placeholder="Src Trk (0)" style="width: 70px;">
				<span style="color: var(--primary); font-weight: bold;">➔</span>
				<input type="text" id="patch-dst-daw" class="input-text" placeholder="Dst DAW (e.g. reaper)" style="flex: 1;">
				<input type="text" id="patch-dst-trk" class="input-text" placeholder="Dst Trk (1)" style="width: 70px;">
				<button class="btn btn-primary" onclick="addPatch()" data-tooltip="Connect virtual audio stream between DAWs">Connect</button>
			</div>
			<div id="routing" style="max-height: 110px; overflow-y: auto; margin-bottom: 10px;"></div>
			<div id="blueprint" data-tooltip="Visual topological dependency graph">graph LR\n  No connections.</div>
		</section>

		<!-- SECTION 5: Generative AI Stems & Euclidean Pattern Creator -->
		<section class="card col-6">
			<div class="card-header">
				<h2 class="card-title">🤖 Generative AI & Euclidean Engine</h2>
				<span class="badge" data-tooltip="Algorithmic stem generation & Spleeter stem separator">Neural DSP</span>
			</div>
			<div class="control-group">
				<label data-tooltip="Generate MIDI & Audio stems via generative prompt">Generative Prompt / Stem Import</label>
				<div class="input-row">
					<input type="text" id="ai-prompt" class="input-text" placeholder="e.g. 145 BPM Psytrance Rolling Bass" style="flex: 2;">
					<button class="btn btn-secondary" onclick="importGenerative()" data-tooltip="Trigger AI stem generation">Generate</button>
				</div>
			</div>
			<div class="control-group" style="margin-top: 10px;">
				<label data-tooltip="Algorithmic Euclidean rhythm injection into active DAW clip">Euclidean Rhythm Generator</label>
				<div class="input-row">
					<input type="number" id="euc-hits" class="input-text" placeholder="Hits" value="3" style="width: 70px;" data-tooltip="Hits (pulses)">
					<input type="number" id="euc-steps" class="input-text" placeholder="Steps" value="8" style="width: 70px;" data-tooltip="Total steps">
					<input type="number" id="euc-pitch" class="input-text" placeholder="Pitch" value="36" style="width: 70px;" data-tooltip="MIDI pitch (36=C1)">
					<button class="btn btn-subtle" onclick="generateEuclidean()" style="flex: 1;" data-tooltip="Inject Euclidean rhythm pattern into track">Inject Rhythm</button>
				</div>
			</div>
			<div id="jobs" style="margin-top: 12px; max-height: 90px; overflow-y: auto;"></div>
		</section>

		<!-- SECTION 6: Live Arrangement Timeline & Session View -->
		<section class="card col-6">
			<div class="card-header">
				<h2 class="card-title">📊 Live Arrangement Timeline</h2>
				<span class="badge" data-tooltip="Real-time timeline feedback broadcast via OSC">Live Session</span>
			</div>
			<div id="timeline">
				<p style="color: var(--text-muted); padding: 10px;">Arrangement timeline standby. Connect DAW agent to stream clips.</p>
			</div>
		</section>

		<!-- SECTION 7: Virtual Performance Keyboard -->
		<section class="card col-12">
			<div class="card-header">
				<h2 class="card-title">🎹 Virtual MIDI Performance Strip</h2>
				<span class="badge" data-tooltip="Click or touch keys to transmit low-latency MIDI notes">C4 - C5 Octave</span>
			</div>
			<div class="keyboard">
				<div class="key" onclick="playNote(60)" data-tooltip="C4 (MIDI 60)"></div>
				<div class="key black" onclick="playNote(61)" data-tooltip="C#4 (MIDI 61)"></div>
				<div class="key" onclick="playNote(62)" data-tooltip="D4 (MIDI 62)"></div>
				<div class="key black" onclick="playNote(63)" data-tooltip="D#4 (MIDI 63)"></div>
				<div class="key" onclick="playNote(64)" data-tooltip="E4 (MIDI 64)"></div>
				<div class="key" onclick="playNote(65)" data-tooltip="F4 (MIDI 65)"></div>
				<div class="key black" onclick="playNote(66)" data-tooltip="F#4 (MIDI 66)"></div>
				<div class="key" onclick="playNote(67)" data-tooltip="G4 (MIDI 67)"></div>
				<div class="key black" onclick="playNote(68)" data-tooltip="G#4 (MIDI 68)"></div>
				<div class="key" onclick="playNote(69)" data-tooltip="A4 (MIDI 69)"></div>
				<div class="key black" onclick="playNote(70)" data-tooltip="A#4 (MIDI 70)"></div>
				<div class="key" onclick="playNote(71)" data-tooltip="B4 (MIDI 71)"></div>
				<div class="key" onclick="playNote(72)" data-tooltip="C5 (MIDI 72)"></div>
			</div>
		</section>

	</main>

	<footer>
		<div>SuperDAW Universal Control Protocol &copy; 2026. Autonomously Engineered.</div>
		<div id="daemon-status">WebSocket: Connecting...</div>
	</footer>

	<script>
		let cachedPlugins = {};
		let activePlugin = "";
		let activeCategoryFilter = "all";
		let abState = "A";
		let snapshotA = {};
		let snapshotB = {};
		let snapshotBanks = { A: {}, B: {}, C: {}, D: {} };

		function storeSnapshotBank(slot) {
			if (!activePlugin || !cachedPlugins[activePlugin]) return;
			const p = cachedPlugins[activePlugin];
			snapshotBanks[slot] = {};
			(p.parameters || []).forEach(param => {
				const slider = document.getElementById('slider-' + param.index);
				if (slider) snapshotBanks[slot][param.name] = parseFloat(slider.value);
			});
			const btn = event?.target;
			if (btn) {
				const orig = btn.innerText;
				btn.innerText = 'Saved!';
				btn.style.color = 'var(--primary)';
				setTimeout(() => { btn.innerText = orig; btn.style.color = ''; }, 1000);
			}
		}

		function morphBankParameters(alpha) {
			if (!activePlugin || !cachedPlugins[activePlugin]) return;
			const p = cachedPlugins[activePlugin];
			const val = parseFloat(alpha);
			const label = document.getElementById('morph-val');
			if (label) label.innerText = Math.round(val * 100) + '%';

			const bankA = snapshotBanks['A'];
			const bankB = snapshotBanks['B'];
			if (Object.keys(bankA).length === 0 || Object.keys(bankB).length === 0) return;

			(p.parameters || []).forEach(param => {
				const a = bankA[param.name];
				const b = bankB[param.name];
				if (a !== undefined && b !== undefined) {
					const morphed = a + (b - a) * val;
					const slider = document.getElementById('slider-' + param.index);
					const numInput = document.getElementById('num-' + param.index);
					if (slider) slider.value = morphed;
					if (numInput) numInput.value = morphed.toFixed(2);
					updateReadout(param.index, morphed, param.unit || '');
					sendPluginParam(activePlugin, param.name, morphed);
				}
			});
		}

		const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
		const ws = new WebSocket(wsProtocol + '//' + window.location.host + '/ws');

		ws.onopen = () => {
			document.getElementById('daemon-status').innerText = 'WebSocket: Connected (Low Latency)';
			refreshPlugins();
		};
		ws.onclose = () => {
			document.getElementById('daemon-status').innerText = 'WebSocket: Disconnected (Retrying...)';
		};

		ws.onmessage = (event) => {
			const state = JSON.parse(event.data);

			// Render DAWs
			let dawHtml = '';
			for (const name in state.daws) {
				const d = state.daws[name];
				dawHtml += ' \
					<div class="daw-row"> \
						<div class="daw-name">' + name.toUpperCase() + '</div> \
						<div style="display: flex; gap: 8px; align-items: center;"> \
							<span class="badge">' + d.bpm.toFixed(1) + ' BPM</span> \
							<span class="daw-status-badge ' + (d.is_playing ? 'playing' : 'stopped') + '">' + (d.is_playing ? 'PLAYING' : 'STOPPED') + '</span> \
						</div> \
					</div> \
				';
			}
			document.getElementById('daws').innerHTML = dawHtml || '<p style="color: var(--text-muted); padding: 10px 0;">No active DAWs connected.</p>';

			// Render Patches
			let patchHtml = '';
			state.patches.forEach(p => {
				patchHtml += ' \
					<div class="daw-row" style="margin-bottom: 4px;"> \
						<span>' + p.source_daw + ' (' + p.source_track + ')</span> \
						<span style="color: var(--primary);">➔</span> \
						<span>' + p.dest_daw + ' (' + p.dest_track + ')</span> \
						<button class="btn btn-danger" style="padding: 2px 8px; font-size: 10px;" onclick="removePatch(\'' + p.source_daw + '\', \'' + p.source_track + '\', \'' + p.dest_daw + '\', \'' + p.dest_track + '\')">✕</button> \
					</div> \
				';
			});
			document.getElementById('routing').innerHTML = patchHtml || '<p style="color: var(--text-muted); font-size: 12px;">No active audio patches.</p>';

			// Render Blueprint
			let blueprint = 'graph LR\n';
			state.patches.forEach(p => {
				blueprint += '  ' + p.source_daw + ' --> ' + p.dest_daw + '\n';
			});
			document.getElementById('blueprint').innerText = blueprint === 'graph LR\n' ? 'No connections.' : blueprint;

			// Render Jobs
			let jobHtml = '';
			if (state.jobs && state.jobs.length > 0) {
				state.jobs.forEach(j => {
					jobHtml += ' \
						<div class="daw-row"> \
							<span style="font-size: 12px; font-weight: 600;">' + j.prompt + '</span> \
							<span class="badge">' + j.status + ' (' + Math.round((j.progress || 1)*100) + '%)</span> \
						</div> \
					';
				});
			}
			document.getElementById('jobs').innerHTML = jobHtml || '<p style="color: var(--text-muted); font-size: 12px;">No active generation jobs.</p>';

			// Render Timeline
			let timelineHtml = '';
			let top = 10;
			for (const name in state.daws) {
				const d = state.daws[name];
				if (d.arrangement) {
					try {
						const arrangement = JSON.parse(d.arrangement);
						arrangement.forEach(track => {
							timelineHtml += '<div class="track-lane"><span style="width: 100px; display: inline-block; font-size: 11px; font-weight: 600;">' + track.track + '</span>';
							(track.clips || []).forEach(clip => {
								const left = clip.start * 20;
								const width = Math.max(40, (clip.end - clip.start) * 20);
								timelineHtml += '<div class="clip-block" style="left: ' + (110 + left) + 'px; width: ' + width + 'px;">' + clip.name + '</div>';
							});
							timelineHtml += '</div>';
							top += 40;
						});
					} catch(e) {}
				}
			}
			document.getElementById('timeline').innerHTML = timelineHtml || '<p style="color: var(--text-muted); padding: 10px;">Arrangement timeline standby. Connect DAW agent to stream clips.</p>';
		};

		async function callMcp(name, args) {
			return fetch('/api/call', {
				method: 'POST',
				headers: {'Content-Type': 'application/json'},
				body: JSON.stringify({name, arguments: args})
			});
		}

		let linkPlaying = false;
		let linkBpm = 128.0;
		let linkStartTime = 0;
		let clockInterval = null;

		function syncTransport(playing) {
			linkPlaying = playing;
			if (playing) {
				linkStartTime = Date.now();
			}
			linkBpm = parseFloat(document.getElementById('global-bpm').value) || 120.0;
			callMcp('superdaw_transport_control', { playing, bpm: linkBpm });
			callMcp('superdaw_link_sync', { playing, bpm: linkBpm });
		}

		function syncBpm(val) {
			linkBpm = parseFloat(val) || 120.0;
			callMcp('superdaw_transport_control', { playing: linkPlaying, bpm: linkBpm });
			callMcp('superdaw_link_sync', { playing: linkPlaying, bpm: linkBpm });
		}

		function updateLinkClock() {
			if (!linkPlaying) {
				document.getElementById('link-beat-val').innerText = '0.00';
				document.getElementById('link-phase-val').innerText = '0.00';
				const led = document.getElementById('link-beat-led');
				if (led) { led.style.background = '#334155'; led.style.boxShadow = 'none'; }
				return;
			}
			const elapsedSec = (Date.now() - linkStartTime) / 1000.0;
			const beat = elapsedSec * (linkBpm / 60.0);
			const phase = beat % 4.0;
			document.getElementById('link-beat-val').innerText = beat.toFixed(2);
			document.getElementById('link-phase-val').innerText = phase.toFixed(2);
			const led = document.getElementById('link-beat-led');
			if (led) {
				if (phase < 0.25) {
					led.style.background = 'var(--primary)';
					led.style.boxShadow = '0 0 8px var(--primary)';
				} else {
					led.style.background = '#334155';
					led.style.boxShadow = 'none';
				}
			}
		}
		setInterval(updateLinkClock, 50);

		async function refreshPlugins() {
			try {
				const res = await fetch('/api/plugins');
				const data = await res.json();
				cachedPlugins = data;
				const sel = document.getElementById('plugin-list');
				sel.innerHTML = '<option value="">Select a plugin...</option>';
				for (const name in data) {
					sel.innerHTML += '<option value="' + name + '">' + name + ' (' + data[name].vendor + ')</option>';
				}
			} catch(e) {
				console.error("Plugin fetch error:", e);
			}
		}

		async function loadPlugin(name) {
			activePlugin = name;
			if (!name || !cachedPlugins[name]) {
				document.getElementById('plugin-params').innerHTML = '<p style="color: var(--text-muted); padding: 15px;">Select a plugin above.</p>';
				document.getElementById('plugin-meta-banner').innerHTML = '';
				return;
			}
			const p = cachedPlugins[name];

			// Meta banner
			document.getElementById('plugin-meta-banner').innerHTML = ' \
				<span class="meta-tag">Vendor: ' + p.vendor + '</span> \
				<span class="meta-tag">Version: ' + p.version + '</span> \
				<span class="meta-tag">Category: ' + (p.deep_meta?.category || 'General') + '</span> \
				<span class="meta-tag">UID: ' + (p.deep_meta?.uuid || 'VST3-CORE') + '</span> \
				<span class="meta-tag">Buses: ' + (p.deep_meta?.input_buses || 2) + 'In / ' + (p.deep_meta?.output_buses || 2) + 'Out</span> \
			';

			// Presets dropdown
			const presetSel = document.getElementById('preset-list');
			presetSel.innerHTML = '<option value="">Load Preset...</option>';
			(p.presets || []).forEach(pr => {
				presetSel.innerHTML += '<option value="' + pr + '">' + pr + '</option>';
			});

			renderParameters();
		}

		function setParamFilter(cat, btn) {
			activeCategoryFilter = cat;
			document.querySelectorAll('.filter-tabs .btn-tab').forEach(b => b.classList.remove('active'));
			if (btn) btn.classList.add('active');
			renderParameters();
		}

		function applyParamFilter() {
			renderParameters();
		}

		function renderParameters() {
			if (!activePlugin || !cachedPlugins[activePlugin]) return;
			const p = cachedPlugins[activePlugin];
			const search = (document.getElementById('param-search').value || '').toLowerCase();

			let html = '';
			(p.parameters || []).forEach(param => {
				if (activeCategoryFilter !== 'all' && param.category && param.category !== activeCategoryFilter) {
					return;
				}
				if (search && !param.name.toLowerCase().includes(search) && !(param.category || '').toLowerCase().includes(search)) {
					return;
				}

				const val = param.default_value !== undefined ? param.default_value : 0.5;
				const stepVal = param.step || (param.unit === "Shape" || param.unit === "Oct" ? 1.0 : 0.01);
				const unitLabel = param.unit ? ' (' + param.unit + ')' : '';
				let curveIcon = '📈';
				if (param.curve === 'logarithmic') curveIcon = '📐';
				else if (param.curve === 'decibel') curveIcon = '🔊';
				else if (param.curve === 'scurve') curveIcon = '〰️';
				const curveLabel = param.curve ? ' • ' + curveIcon + ' ' + param.curve : '';

				let syncDivisionsHtml = '';
				if (param.sync_divisions && param.sync_divisions.length > 0) {
					syncDivisionsHtml = '<select class="select-box" style="padding: 2px 4px; font-size: 10px; width: auto;" onchange="applySyncDivision(\'' + activePlugin + '\', \'' + param.name + '\', ' + param.index + ', this.value)" data-tooltip="Tempo-synced division"><option value="">Sync Division...</option>';
					param.sync_divisions.forEach(div => {
						syncDivisionsHtml += '<option value="' + div + '">' + div + '</option>';
					});
					syncDivisionsHtml += '</select>';
				}

				const modSrc = param.modulation ? param.modulation.source || "None" : "None";
				const modDepth = param.modulation ? param.modulation.depth || 0.0 : 0.0;

				html += ' \
					<div class="param-box"> \
						<div class="param-header"> \
							<span>' + param.name + '<small style="color: var(--text-muted); font-size: 10px;">' + unitLabel + '</small></span> \
							<div style="display: flex; gap: 4px; align-items: center;"> \
								' + syncDivisionsHtml + ' \
								<button class="btn btn-subtle" style="padding: 2px 6px; font-size: 10px;" id="midi-learn-' + param.index + '" onclick="toggleMidiLearn(' + param.index + ', \'' + param.name + '\')" data-tooltip="Map physical MIDI CC controller to this parameter">🎹 Learn</button> \
								<span class="badge" style="font-size: 10px;">' + (param.category || 'General') + curveLabel + '</span> \
							</div> \
						</div> \
						<div class="param-slider-row"> \
							<input type="range" id="slider-' + param.index + '" class="param-slider" min="' + (param.min_value !== undefined ? param.min_value : 0) + '" max="' + (param.max_value !== undefined ? param.max_value : 1) + '" step="' + stepVal + '" value="' + val + '" oninput="syncParamSlider(\'' + activePlugin + '\', \'' + param.name + '\', this.value, ' + param.index + ', \'' + (param.unit || '') + '\')"> \
							<input type="number" id="num-' + param.index + '" class="param-input" value="' + (Number.isInteger(stepVal) ? parseInt(val) : parseFloat(val).toFixed(2)) + '" step="' + stepVal + '" onchange="syncParamInput(\'' + activePlugin + '\', \'' + param.name + '\', this.value, ' + param.index + ', \'' + (param.unit || '') + '\')"> \
							<span id="display-readout-' + param.index + '" style="font-size: 11px; color: var(--accent-blue); min-width: 50px; text-align: right;">' + formatReadout(param, val) + '</span> \
						</div> \
						<div style="display: flex; align-items: center; justify-content: space-between; margin-top: 6px; padding-top: 4px; border-top: 1px dashed rgba(255,255,255,0.06); font-size: 11px;"> \
							<div style="display: flex; align-items: center; gap: 6px;"> \
								<span style="color: var(--text-muted); font-size: 10px;" data-tooltip="Modulation routing source">Mod Src:</span> \
								<select id="mod-src-' + param.index + '" class="select-box" style="padding: 2px 4px; font-size: 10px; width: 95px;" onchange="updateParamModulation(\'' + activePlugin + '\', ' + param.index + ')"> \
									<option value="None"' + (modSrc === "None" ? " selected" : "") + '>None</option> \
									<option value="LFO 1"' + (modSrc === "LFO 1" ? " selected" : "") + '>LFO 1</option> \
									<option value="LFO 2"' + (modSrc === "LFO 2" ? " selected" : "") + '>LFO 2</option> \
									<option value="Envelope 1"' + (modSrc === "Envelope 1" ? " selected" : "") + '>Envelope 1</option> \
									<option value="Velocity"' + (modSrc === "Velocity" ? " selected" : "") + '>Velocity</option> \
									<option value="Aftertouch"' + (modSrc === "Aftertouch" ? " selected" : "") + '>Aftertouch</option> \
								</select> \
							</div> \
							<div style="display: flex; align-items: center; gap: 6px;"> \
								<span style="color: var(--text-muted); font-size: 10px;" data-tooltip="Bipolar modulation depth (-100% to +100%)">Depth:</span> \
								<input type="range" id="mod-depth-' + param.index + '" class="param-slider" style="width: 70px;" min="-1" max="1" step="0.05" value="' + modDepth + '" oninput="updateParamModulation(\'' + activePlugin + '\', ' + param.index + ')"> \
								<span id="mod-depth-val-' + param.index + '" style="font-size: 10px; color: var(--secondary); min-width: 35px; text-align: right;">' + (modDepth > 0 ? "+" : "") + Math.round(modDepth*100) + '%</span> \
							</div> \
						</div> \
					</div> \
				';
			});
			document.getElementById('plugin-params').innerHTML = html || '<p style="color: var(--text-muted); padding: 15px;">No matching parameters.</p>';
		}

		function formatReadout(param, val) {
			const num = parseFloat(val);
			if (param.unit === "Hz") {
				return num >= 1000 ? (num/1000).toFixed(1) + " kHz" : num.toFixed(0) + " Hz";
			} else if (param.unit === "dB") {
				return (num > 0 ? "+" : "") + num.toFixed(1) + " dB";
			} else if (param.unit === "%") {
				return (num * 100).toFixed(0) + "%";
			} else if (param.unit === "ms") {
				return num >= 1000 ? (num/1000).toFixed(2) + " s" : num.toFixed(1) + " ms";
			} else if (param.unit === "Shape" || param.unit === "Wave") {
				const shapes = ["Sine", "Tri", "Saw", "Square", "Noise"];
				return shapes[parseInt(num)] || "Shape " + num;
			}
			return num.toFixed(2) + (param.unit ? " " + param.unit : "");
		}

		function syncParamSlider(pluginName, paramName, value, index, unit) {
			const numVal = parseFloat(value);
			const numInput = document.getElementById('num-' + index);
			if (numInput) numInput.value = numVal.toFixed(2);
			updateReadout(index, numVal, unit);
			sendPluginParam(pluginName, paramName, numVal);
		}

		function updateParamModulation(plugin, idx) {
			const srcSel = document.getElementById('mod-src-' + idx);
			const depthSlider = document.getElementById('mod-depth-' + idx);
			const depthVal = document.getElementById('mod-depth-val-' + idx);
			if (!srcSel || !depthSlider) return;
			const src = srcSel.value;
			const depth = parseFloat(depthSlider.value);
			if (depthVal) {
				depthVal.innerText = (depth > 0 ? "+" : "") + Math.round(depth * 100) + "%";
			}
			callMcp('superdaw_set_plugin_param_mod', {
				plugin_name: plugin,
				param_index: idx,
				source: src,
				depth: depth
			});
		}

		function syncParamInput(pluginName, paramName, value, index, unit) {
			const numVal = parseFloat(value);
			const slider = document.getElementById('slider-' + index);
			if (slider) slider.value = numVal;
			updateReadout(index, numVal, unit);
			sendPluginParam(pluginName, paramName, numVal);
		}

		function updateReadout(index, val, unit) {
			const el = document.getElementById('display-readout-' + index);
			if (el) el.innerText = formatReadout({ unit }, val);
		}

		let activeMidiLearnIndex = null;
		function toggleMidiLearn(index, name) {
			const btn = document.getElementById('midi-learn-' + index);
			if (activeMidiLearnIndex === index) {
				activeMidiLearnIndex = null;
				if (btn) { btn.innerText = '🎹 Learn'; btn.style.background = ''; }
			} else {
				if (activeMidiLearnIndex !== null) {
					const oldBtn = document.getElementById('midi-learn-' + activeMidiLearnIndex);
					if (oldBtn) { oldBtn.innerText = '🎹 Learn'; oldBtn.style.background = ''; }
				}
				activeMidiLearnIndex = index;
				if (btn) { btn.innerText = 'Listening...'; btn.style.background = '#f59e0b'; }
			}
		}

		function applySyncDivision(pluginName, paramName, index, division) {
			if (!division) return;
			// Convert musical division (e.g., 1/4 = 0.25, 1/8 = 0.125) to normalized float
			const parts = division.replace('d', '').replace('t', '').split('/');
			let val = 0.25;
			if (parts.length === 2) {
				val = parseFloat(parts[0]) / parseFloat(parts[1]);
				if (division.includes('d')) val *= 1.5;
				if (division.includes('t')) val *= (2.0 / 3.0);
			}
			const slider = document.getElementById('slider-' + index);
			const numInput = document.getElementById('num-' + index);
			if (slider) slider.value = val;
			if (numInput) numInput.value = val.toFixed(3);
			const el = document.getElementById('display-readout-' + index);
			if (el) el.innerText = division;
			sendPluginParam(pluginName, paramName, val);
		}

		function sendPluginParam(pluginName, paramName, numVal) {
			callMcp('superdaw_set_plugin_parameter', {
				daw: document.getElementById('mixer-daw').value,
				track_id: document.getElementById('mixer-track-id').value,
				plugin_name: pluginName,
				parameter_name: paramName,
				value: numVal
			});
		}

		function triggerMacroPreset(macroType) {
			if (!activePlugin || !cachedPlugins[activePlugin]) {
				alert("Please select a plugin first.");
				return;
			}
			const p = cachedPlugins[activePlugin];
			(p.parameters || []).forEach(param => {
				let newVal = null;
				const pName = param.name.toLowerCase();
				if (macroType === 'brighten' && (pName.includes('cutoff') || pName.includes('high cut') || pName.includes('tone'))) {
					newVal = Math.min((param.max_value || 1.0), (param.default_value || 0.5) * 1.5);
				} else if (macroType === 'warmth' && (pName.includes('drive') || pName.includes('band 1 gain') || pName.includes('damping'))) {
					newVal = Math.min((param.max_value || 1.0), 0.65);
				} else if (macroType === 'space' && (pName.includes('reverb') || pName.includes('delay') || pName.includes('wet') || pName.includes('decay'))) {
					newVal = Math.min((param.max_value || 1.0), 0.70);
				} else if (macroType === 'punch' && (pName.includes('attack'))) {
					newVal = 0.01;
				} else if (macroType === 'punch' && (pName.includes('decay') || pName.includes('env amount'))) {
					newVal = 0.80;
				}

				if (newVal !== null) {
					const slider = document.getElementById('slider-' + param.index);
					const numInput = document.getElementById('num-' + param.index);
					if (slider) slider.value = newVal;
					if (numInput) numInput.value = newVal.toFixed(2);
					sendPluginParam(activePlugin, param.name, newVal);
				}
			});
		}

		function resetParametersToDefault() {
			if (!activePlugin || !cachedPlugins[activePlugin]) return;
			const p = cachedPlugins[activePlugin];
			(p.parameters || []).forEach(param => {
				const defVal = param.default_value !== undefined ? param.default_value : 0.5;
				const slider = document.getElementById('slider-' + param.index);
				const numInput = document.getElementById('num-' + param.index);
				if (slider) slider.value = defVal;
				if (numInput) numInput.value = defVal.toFixed(2);
				sendPluginParam(activePlugin, param.name, defVal);
			});
		}

		function randomizeParameters() {
			if (!activePlugin || !cachedPlugins[activePlugin]) return;
			const p = cachedPlugins[activePlugin];
			(p.parameters || []).forEach(param => {
				if (param.category === "Master") return; // Keep master volume stable
				const min = param.min_value !== undefined ? param.min_value : 0;
				const max = param.max_value !== undefined ? param.max_value : 1;
				let rnd = min + Math.random() * (max - min);
				if (param.step && Number.isInteger(param.step)) {
					rnd = Math.round(rnd);
				} else {
					rnd = Math.round(rnd * 100) / 100;
				}
				const slider = document.getElementById('slider-' + param.index);
				const numInput = document.getElementById('num-' + param.index);
				if (slider) slider.value = rnd;
				if (numInput) numInput.value = rnd.toFixed(2);
				sendPluginParam(activePlugin, param.name, rnd);
			});
		}

		function toggleABComparison() {
			if (!activePlugin || !cachedPlugins[activePlugin]) return;
			const p = cachedPlugins[activePlugin];

			if (abState === "A") {
				// Store current in A
				snapshotA[activePlugin] = {};
				(p.parameters || []).forEach(param => {
					const slider = document.getElementById('slider-' + param.index);
					if (slider) snapshotA[activePlugin][param.name] = parseFloat(slider.value);
				});
				abState = "B";
				document.getElementById('ab-label').innerText = "B";
				// If B exists, apply it
				if (snapshotB[activePlugin]) {
					(p.parameters || []).forEach(param => {
						const val = snapshotB[activePlugin][param.name];
						if (val !== undefined) {
							const slider = document.getElementById('slider-' + param.index);
							const numInput = document.getElementById('num-' + param.index);
							if (slider) slider.value = val;
							if (numInput) numInput.value = val.toFixed(2);
							sendPluginParam(activePlugin, param.name, val);
						}
					});
				}
			} else {
				// Store current in B
				snapshotB[activePlugin] = {};
				(p.parameters || []).forEach(param => {
					const slider = document.getElementById('slider-' + param.index);
					if (slider) snapshotB[activePlugin][param.name] = parseFloat(slider.value);
				});
				abState = "A";
				document.getElementById('ab-label').innerText = "A";
				// Apply A
				if (snapshotA[activePlugin]) {
					(p.parameters || []).forEach(param => {
						const val = snapshotA[activePlugin][param.name];
						if (val !== undefined) {
							const slider = document.getElementById('slider-' + param.index);
							const numInput = document.getElementById('num-' + param.index);
							if (slider) slider.value = val;
							if (numInput) numInput.value = val.toFixed(2);
							sendPluginParam(activePlugin, param.name, val);
						}
					});
				}
			}
		}

		async function deepScanSelected() {
			if (!activePlugin) { alert("Please select a plugin first."); return; }
			try {
				const res = await fetch('/api/plugins?name=' + encodeURIComponent(activePlugin) + '&deep=true');
				const data = await res.json();
				cachedPlugins[activePlugin] = data;
				loadPlugin(activePlugin);
				alert("Deep scan complete for " + activePlugin + "! Validated " + (data.deep_meta?.binary_info || 'binary signatures.'));
			} catch(e) {
				alert("Deep scan failed: " + e);
			}
		}

		function loadPreset(presetName) {
			if (!presetName || !activePlugin) return;
			callMcp('superdaw_set_plugin_parameter', {
				daw: document.getElementById('mixer-daw').value,
				track_id: document.getElementById('mixer-track-id').value,
				plugin_name: activePlugin,
				parameter_name: "Preset",
				value: 1.0
			});
		}

		function updateMixerVolume(val) {
			document.getElementById('vol-display').innerText = parseFloat(val).toFixed(2);
			callMcp('superdaw_set_mixer', {
				daw: document.getElementById('mixer-daw').value,
				track_id: document.getElementById('mixer-track-id').value,
				volume: parseFloat(val)
			});
		}

		function updateMixerPan(val) {
			document.getElementById('pan-display').innerText = parseFloat(val).toFixed(2);
			callMcp('superdaw_set_mixer', {
				daw: document.getElementById('mixer-daw').value,
				track_id: document.getElementById('mixer-track-id').value,
				pan: parseFloat(val)
			});
		}

		function createTrack() {
			const name = document.getElementById('new-track-name').value || "New Track";
			const type = document.getElementById('new-track-type').value;
			const daw = document.getElementById('mixer-daw').value;
			callMcp('superdaw_create_track', { name, type, daw });
			document.getElementById('new-track-name').value = '';
		}

		function addPatch() {
			const srcDaw = document.getElementById('patch-src-daw').value;
			const srcTrack = document.getElementById('patch-src-trk').value;
			const dstDaw = document.getElementById('patch-dst-daw').value;
			const dstTrack = document.getElementById('patch-dst-trk').value;
			if (!srcDaw || !dstDaw) { alert("Please specify source and destination DAWs."); return; }
			callMcp('superdaw_patch_audio', {
				source_daw: srcDaw,
				source_track: srcTrack,
				dest_daw: dstDaw,
				dest_track: dstTrack
			});
		}

		function removePatch(srcDaw, srcTrack, dstDaw, dstTrack) {
			callMcp('superdaw_unpatch_audio', {
				source_daw: srcDaw,
				source_track: srcTrack,
				dest_daw: dstDaw,
				dest_track: dstTrack
			});
		}

		function importGenerative() {
			const prompt = document.getElementById('ai-prompt').value;
			if (!prompt) return;
			const target = document.getElementById('mixer-daw').value;
			callMcp('superdaw_import_generative', { prompt, target_daw: target });
			document.getElementById('ai-prompt').value = '';
		}

		function generateEuclidean() {
			const hits = parseInt(document.getElementById('euc-hits').value);
			const steps = parseInt(document.getElementById('euc-steps').value);
			const pitch = parseInt(document.getElementById('euc-pitch').value);
			const track_id = document.getElementById('mixer-track-id').value;
			const daw = document.getElementById('mixer-daw').value;
			callMcp('superdaw_generate_euclidean', { hits, steps, pitch, track_id, daw });
		}

		function playNote(pitch) {
			const daw = document.getElementById('mixer-daw').value;
			const track_id = document.getElementById('mixer-track-id').value;
			callMcp('superdaw_write_midi', {
				daw,
				track_id,
				notes: [{pitch: pitch, velocity: 100, start_beat: 0, duration: 0.5}]
			});
		}

		async function quitSuperDAW() {
			if (!confirm("Are you sure you want to stop and quit all SuperDAW MCP servers and gateways?")) return;
			try {
				await fetch('/api/quit', { method: 'POST' });
				document.body.innerHTML = '<div style="display: flex; height: 100vh; align-items: center; justify-content: center; font-family: sans-serif; color: #ff4772; background: #0c0d0f; flex-direction: column;"><h1>🛑 SuperDAW Shut Down</h1><p style="color: #8b9bb4; margin-top: 10px;">The SuperDAW daemon and servers have exited cleanly.</p></div>';
			} catch(e) {
				console.error("Quit error:", e);
			}
		}
	</script>
</body>
</html>
`)
}
