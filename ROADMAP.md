# SuperDAW-MCP ROADMAP

## Current Milestone (v3.11.0 - Active)
- [x] **Universal Link Synchronization & Studio Clock**: Unified network clock arbiter (`LinkBridge`), tempo ramping, beat time calculation, and 4/4 phase quantization alignment.
- [x] **Link MCP Protocol Integration**: Added `superdaw_link_sync` and `superdaw_link_status` MCP tools with JSON schemas and gateway routing.
- [x] **Master Studio Console Link Clock HUD**: Real-time beat counter, phase bar readout, tempo pulse indicator, and peer network status.
- [x] **Desktop System Tray Controller (`scripts/superdaw_tray.py`)**: Platform status monitor, port inspection (HTTP 8081, OSC 12001, TCP 12002), and graceful termination integration via `/api/quit`.
- [x] **Unified Master Console Preset Snapshot Banks & Continuous Morphing**: Preset bank slots (A, B, C, D) and continuous parameter morph slider in Master Console UI.
- [x] **Mobile Touch Remote VST3 & Preset Morphing Integration**: Mobile remote and React Native `PluginScreen` with snapshot slots, tactile +/- step buttons, and continuous morph slider.
- [x] **Bidirectional Mathematical Curve Calibration**: High-precision `NormalizedToPhysical` and `PhysicalToNormalized` mapping for Logarithmic, Decibel ($10^{\text{dB}/20}$), S-Curve, and Linear responses.
- [x] **Consolidated Single-Page Dashboard**: Unified Master Studio Console categorizing features by high value (Transport, VST3 Inspector, Mixer, Audio Patching Matrix, Generative DSP, Arrangement Timeline, Virtual MIDI Performance).
- [x] **Advanced VST3 Parameter Discovery**: Deep scanning with moduleinfo.json & binary header validation, preset extraction, calibrated parameter ranges, and synchronous API endpoints (`/api/plugins`).
- [x] **Sound Design Quick Macro Knobs & State Tools**: One-click sound design macros (Brighten, Warmth, Space Out, Punch), musical parameter randomization, and factory reset.
- [x] **Human-Readable Readouts & MIDI Learn**: Live display readout formatting, interactive MIDI Learn simulation toggles, and tempo-synced musical division pickers.
- [x] **Desktop / Process Governance**: Full server shutdown & quit endpoint (`/api/quit`) with responsive controls for system tray & dashboard management.
- [x] **DAW Parity Hardening**: Unified interface implementation fixes across Logic Pro, Pro Tools, Ableton Live, Bitwig, FL Studio, REAPER, Ardour, and Cubase.

## Upcoming Milestones (v3.12.0 - Future)
- [ ] **Native CGO VST3 Deep-Scanning Bridge**: Direct `libvst3` dynamic loading (`dlopen`/`LoadLibrary`) for direct parameter queries without host DAW bridging.
- [ ] **Multi-User Realtime Collaborative Studios**: Operational WebRTC / WebSocket room synchronization for distributed jam sessions and shared parameter locking.
- [ ] **Expanded Mobile Remote**: Touch-screen clip launching grid and mixer fader banks in React Native.
