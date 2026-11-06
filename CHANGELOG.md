# CHANGELOG

## [v3.3.0] - 2026-11-06
### Added
- **Refined VST3 Discovery & Sound Design Macro Controls**: Expanded VST3 heuristic parameter generation with multi-oscillator parameters (Osc 1 Waveform, Osc 1 Octave, Osc 2 Detune, Sub Osc Level) and binary signature heuristics.
- **Interactive Macro Knobs in Plugin Inspector**: Added real-time sound design macro buttons (✨ Brighten, 🔥 Warmth, 🌌 Space Out, 🥊 Snappy / Punch) connected to the active plugin parameter bus.
- **Granular Category Filtering & Bidirectional Synchronization**: Filter parameter controls by functional category (Filter/Tone, Envelopes, Mod/LFO, FX/Mix, Master) with live search and slider-input sync.

## [v3.2.0] - 2026-11-06
### Added
- **Consolidated Master Studio Console**: Redesigned the entire dashboard into a single, unified responsive control surface organized by high-value functional domains.
- **Deep VST3 Parameter Discovery**: Implemented moduleinfo.json parsing, binary header probing heuristics, preset enumeration, and calibrated parameter ranges.
- **Categorized Parameter Inspector UX**: Added functional group filter tabs (Filter/Tone, Envelopes, Mod/LFO, FX/Mix, Master), instant search bar, and bidirectional numeric input sync with parameter sliders.
- **Synchronous Plugin Query API**: Added `/api/plugins` endpoint supporting direct querying of all plugins, single plugin parameter sets, and triggerable deep scans (`?deep=true`).
- **Graceful Server Shutdown & Process Control**: Added `/api/quit` endpoint and dashboard button to shut down background servers, TCP gateways, and OSC bridges cleanly.
- **Interactive VST3 Control Surface**: Dynamic parameter sliders with real-time numeric readouts, unit labels, and preset selector wired directly to DAW parameter automation.
- **Comprehensive Tooltip System**: Added contextual tooltip descriptions across all console controls, transport buttons, patching fields, and generative parameters.

### Fixed
- Fixed compilation failures in `pkg/daw/logic.go` and `pkg/daw/protools.go` caused by duplicate method definitions (`SendCC`, `SetPluginParameter`) and redundant imports.
- Archived conflicting root-level `test_*.go` scripts into `archive/test_scripts.tar.gz` to restore clean Go workspace testing.
- Resolved git merge conflict markers in `TODO.md`.

## [v3.1.0] - 2026-06-17
### Added
- Initial VST3 heuristic parameter discovery and basic Plugin Inspector module.
- Ableton Live 10 Standard instrument mapping and deferred loading.
- Submodule audit and clean consolidation into native core drivers.

## [v3.0.0] - 2026-04-12
### Added
- Multi-instance connection manager and bidirectional OSC telemetry.
- TCP JSON-RPC Gateway on port 12002.
