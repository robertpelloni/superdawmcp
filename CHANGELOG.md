# CHANGELOG

## [v3.16.0] - 2026-11-08
### Added
- **Fine-Grained Parameter Metadata & Polarity Badging**: Added `IsBipolar`, `SmoothingMs`, `CanAutomate`, and `IsReadOnly` fields to `ParamMetadata` with helper methods `ClampValue` and `ResetToDefault` in `pkg/vst/scanner.go`.
- **Plugin Inspector Double-Click & Single Parameter Reset**: Added individual parameter reset buttons (`↺`) and double-click reset directly on parameter faders in Master Console (`pkg/ux/dashboard/ui.go`) and mobile remote (`pkg/ux/dashboard/remote.go`).
- **Visual Polarity & Slew Rate Indicators**: Added visual `± Bipolar` / `+ Uni` and `~Smoothing` pills on each parameter card in the Master Console.

## [v3.15.0] - 2026-11-08
### Added
- **Unified Master Console Preset Snapshot Banks & Continuous Morphing**: Fully integrated snapshot bank slot buttons (Slot A, Slot B, Slot C, Slot D) with visual active indicators and a continuous A ⟷ B morph slider into the Master Console Sound Design section in `pkg/ux/dashboard/ui.go`.
- **Dynamic Cross-Parameter State Morphing**: Implemented real-time parameter interpolation logic (`morphBankParameters`) crossfading between Slot A and Slot B with live percentage readout and DAW automation dispatch.

## [v3.14.0] - 2026-11-08
### Added
- **Mobile Touch Remote Fine-Nudge Controls**: Added tactile -0.05 and +0.05 step nudging buttons directly inside `mobile/src/screens/PluginScreen.tsx` for precise touch modulation without slider jumping.
- **Synchronized Mobile & Web Remote Stepping**: Unified stepping behavior across React Native and Web `/remote` interfaces.

## [v3.13.0] - 2026-11-08
### Added
- **Bidirectional Parameter Curve Conversion & Snapping**: Implemented `NormalizedToPhysical`, `PhysicalToNormalized`, and `SnapToStep` methods on `ParamMetadata` in `pkg/vst/scanner.go` supporting Logarithmic, Decibel, Sigmoidal S-Curve, and Stepped snapping.
- **Enhanced Curve Visualization**: Added distinct mathematical curve icons (📐 logarithmic, 🔊 decibel, 〰️ scurve, 📈 linear) and interactive badges in the Master Console Plugin Inspector (`pkg/ux/dashboard/ui.go`).
- **Comprehensive Curve Test Suite**: Added `TestParamMetadata_CurveConversionsAndSnapping` in `pkg/vst/scanner_test.go` verifying round-trip fidelity and boundary clamping.

## [v3.12.0] - 2026-11-08
### Added
- **VST3 Parameter Modulation Routing & Automation Engine**: Implemented `SetParamModulation` and `GetParamModulation` in `pkg/vst/scanner.go` supporting LFO 1/2, Envelope 1/2, Velocity, and Aftertouch modulation sources with bipolar depth scaling (-100% to +100%).
- **MCP Modulation Protocol Tools**: Registered `superdaw_set_plugin_param_mod` and `superdaw_get_plugin_param_mod` MCP tools for programmatic AI modulation patching.
- **Enhanced Master Console & Mobile Inspector Controls**: Integrated interactive modulation source dropdowns and bipolar depth sliders directly beneath each parameter card on the Master Console (`pkg/ux/dashboard/ui.go`) and mobile remote (`pkg/ux/dashboard/remote.go`).

## [v3.11.0] - 2026-11-08
### Added
- **Universal Link Synchronization & Studio Clock**: Expanded `LinkBridge` in `pkg/engine/link.go` with synchronized tempo ramping, quantization phase alignment, beat calculation, and peer monitoring.
- **MCP Link Protocol Tools**: Registered `superdaw_link_sync` and `superdaw_link_status` in protocol schemas and MCP message dispatch loop.
- **Master Console Universal Link Clock HUD**: Integrated real-time beat counter, 4/4 bar phase readout, tempo pulse LED, and bidirectional tempo sync into `pkg/ux/dashboard/ui.go`.

## [v3.10.0] - 2026-11-08
### Added
- **Desktop System Tray Controller (`scripts/superdaw_tray.py`)**: Implemented status inspection, port monitoring (HTTP 8081, OSC 12001, TCP 12002), and graceful shutdown integration via `/api/quit`.
- **Studio Health Diagnostic Fix (`scripts/studio_health.py`)**: Fixed `get_studio_status` integration to properly report active DAW connections and transport statuses.

## [v3.9.0] - 2026-11-07
### Added
- **Unified Master Console Preset Snapshot Banks & Continuous Morphing**: Added snapshot slot buttons (Slot A, Slot B, Slot C, Slot D) and continuous parameter interpolation morph slider in `pkg/ux/dashboard/ui.go`.
- **Touch Mobile Remote Preset Morphing Integration**: Fully wired React Native `mobile/src/screens/PluginScreen.tsx` with snapshot slots A/B/C/D, continuous morph slider, and tactile step buttons.
- **Robust Stdio Mode Detection**: Added command-line flag awareness (`--headless`/`--daemon`) in `cmd/superdaw/main.go` to support automated testing while guaranteeing background server persistence.

## [v3.8.0] - 2026-11-07
### Added
- **Bidirectional Mathematical Curve Engine**: Implemented `NormalizedToPhysical` and `PhysicalToNormalized` across Logarithmic, Decibel, and S-Curves.
- **Preset State Snapshots Engine**: Added `PresetManager` with snapshot bank slots (A, B, C, D) and continuous state morphing.

## [v3.7.0] - 2026-11-07
### Added
- **Consolidated Unified Dashboard Page**: Merged separate subpages into a high-visibility, single responsive command console with priority-ranked cards and tooltips.

## [v3.6.0] - 2026-11-06
### Added
- **Human-Readable Parameter Readout & Display Value Formatting**: Added `FormatDisplayValue` in `pkg/vst/scanner.go` formatting raw parameter floats into musical frequencies ("kHz/Hz"), decibel gains ("dB"), percentages, envelope time units ("ms/s"), and waveform shape names.
- **Dynamic Display Readout & MIDI Learn in Plugin Inspector**: Live formatted value readouts next to parameter inputs, interactive MIDI Learn simulation buttons for hardware CC controller mapping, and tempo-synced musical division pickers (1/4, 1/8d, 1/16, etc.).

## [v3.5.0] - 2026-11-06
### Added
- **Touch Mobile Remote VST3 Touch Inspector**: Added touch-friendly plugin selector, macro buttons (Brighten, Warmth, Space Out, Punch), and live parameter faders on the `/remote` endpoint (`pkg/ux/dashboard/remote.go`).
- **React Native Mobile Plugin Inspector Screen**: Implemented `mobile/src/screens/PluginScreen.tsx` with VST3 parameter modulation sliders, macro buttons, and wired `setPluginParam` in `mobile/src/hooks/useMcp.ts` and `mobile/App.tsx`.
- **Daemon Background Persistence**: Refined daemon standard input handling to prevent premature exit on EOF and sustain background HTTP, OSC, and TCP gateway operations.

## [v3.4.0] - 2026-11-06
### Added
- **Enhanced VST3 Parameter Mapping Accuracy**: Added logarithmic, decibel, and linear response curve classifications and precision step resolutions (`Step`, `Curve`) across synths, dynamics processors (compressors, limiters), and filters.
- **Plugin Inspector Sound Design & State Tools**: Added A/B Comparison state toggling (storing and toggling parameter snapshots A vs B), controlled musical parameter randomization (`±25%` safe margin), and factory default reset.
- **Dynamic Step & Curve Labels**: Parameter sliders now dynamically adapt their step resolution and display parameter curve types (e.g. logarithmic frequency, dB gain) alongside unit badges.

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
