# SuperDAW-MCP ROADMAP

## Current Milestone (v3.9.0 - Active)
- [x] **Unified Master Console Preset Snapshot Banks & Continuous Morphing**: Preset bank slots (A, B, C, D) and continuous parameter morph slider in Master Console UI.
- [x] **Mobile Touch Remote VST3 & Preset Morphing Integration**: Mobile remote and React Native `PluginScreen` with snapshot slots, tactile +/- step buttons, and continuous morph slider.
- [x] **Bidirectional Mathematical Curve Calibration**: High-precision `NormalizedToPhysical` and `PhysicalToNormalized` mapping for Logarithmic, Decibel ($10^{\text{dB}/20}$), S-Curve, and Linear responses.
- [x] **Consolidated Single-Page Dashboard**: Unified Master Studio Console categorizing features by high value (Transport, VST3 Inspector, Mixer, Audio Patching Matrix, Generative DSP, Arrangement Timeline, Virtual MIDI Performance).
- [x] **Advanced VST3 Parameter Discovery**: Deep scanning with moduleinfo.json & binary header validation, preset extraction, calibrated parameter ranges, and synchronous API endpoints (`/api/plugins`).
- [x] **Sound Design Quick Macro Knobs & State Tools**: One-click sound design macros (Brighten, Warmth, Space Out, Punch), musical parameter randomization, and factory reset.
- [x] **Human-Readable Readouts & MIDI Learn**: Live display readout formatting, interactive MIDI Learn simulation toggles, and tempo-synced musical division pickers.
- [x] **Desktop / Process Governance**: Full server shutdown & quit endpoint (`/api/quit`) with responsive controls for system tray & dashboard management.
- [x] **DAW Parity Hardening**: Unified interface implementation fixes across Logic Pro, Pro Tools, Ableton Live, Bitwig, FL Studio, REAPER, Ardour, and Cubase.

## Upcoming Milestones (v3.10.0 - Future)
- [ ] **Hardware System Tray Daemon**: Native desktop systray icon (macOS menu bar / Windows system tray / Linux AppIndicator) with background toggle and process manager.
- [ ] **Native CGO VST3 Deep-Scanning Bridge**: Direct `libvst3` dynamic loading (`dlopen`/`LoadLibrary`) for direct parameter queries without host DAW bridging.
- [ ] **Multi-User Realtime Collaborative Studios**: Operational WebRTC / WebSocket room synchronization for distributed jam sessions and shared parameter locking.
- [ ] **Expanded Mobile Remote**: Touch-screen clip launching grid and mixer fader banks in React Native.
