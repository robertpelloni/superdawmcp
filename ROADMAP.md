# SuperDAW-MCP ROADMAP

## Current Milestone (v3.15.0 - Active)
- [x] **Unified Master Console Preset Snapshot Banks & Continuous Morphing**: Integrated snapshot slots (Slot A, Slot B, Slot C, Slot D) with active selection highlights and dynamic A ⟷ B continuous parameter morph slider in `pkg/ux/dashboard/ui.go`.
- [x] **Mobile Touch Remote Fine-Nudge Stepping**: Added tactile -0.05 and +0.05 parameter adjustment buttons in `mobile/src/screens/PluginScreen.tsx`.
- [x] **Bidirectional Parameter Curve Conversion & Snapping**: Precision normalized-to-physical and physical-to-normalized conversion across logarithmic, decibel, S-curve, and stepped snapping.
- [x] **Curve Visualization & Badging**: Dedicated mathematical curve icons in Master Console parameter cards.
- [x] **VST3 Parameter Modulation Routing & Automation Engine**: Bipolar depth scaling (-1.0 to +1.0) and modulation sources (LFO 1/2, Envelope 1/2, Velocity, Aftertouch).
- [x] **MCP Modulation Protocol Tools**: Added `superdaw_set_plugin_param_mod` and `superdaw_get_plugin_param_mod` tools.
- [x] **Master Console & Mobile Remote Mod Controls**: Real-time modulation source dropdown and bipolar depth slider beneath each parameter card.
- [x] **Universal Link Synchronization & Studio Clock**: Unified network clock arbiter (`LinkBridge`), tempo ramping, beat time calculation, and 4/4 phase quantization alignment.
- [x] **Link MCP Protocol Integration**: Added `superdaw_link_sync` and `superdaw_link_status` MCP tools with JSON schemas and gateway routing.
- [x] **Master Studio Console Link Clock HUD**: Real-time beat counter, phase bar readout, tempo pulse indicator, and peer network status.
- [x] **Desktop System Tray Controller (`scripts/superdaw_tray.py`)**: Platform status monitor, port inspection (HTTP 8081, OSC 12001, TCP 12002), and graceful termination integration via `/api/quit`.
- [x] **Mobile Touch Remote VST3 & Preset Morphing Integration**: Mobile remote and React Native `PluginScreen` with snapshot slots, tactile +/- step buttons, and continuous morph slider.
- [x] **Bidirectional Mathematical Curve Calibration**: High-precision `NormalizedToPhysical` and `PhysicalToNormalized` mapping for Logarithmic, Decibel ($10^{\text{dB}/20}$), S-Curve, and Linear responses.
- [x] **Consolidated Single-Page Dashboard**: Unified Master Studio Console categorizing features by high value (Transport, VST3 Inspector, Mixer, Audio Patching Matrix, Generative DSP, Arrangement Timeline, Virtual MIDI Performance).
- [x] **Advanced VST3 Parameter Discovery**: Deep scanning with moduleinfo.json & binary header validation, preset extraction, calibrated parameter ranges, and synchronous API endpoints (`/api/plugins`).
- [x] **Sound Design Quick Macro Knobs & State Tools**: One-click sound design macros (Brighten, Warmth, Space Out, Punch), musical parameter randomization, and factory reset.
- [x] **Human-Readable Readouts & MIDI Learn**: Live display readout formatting, interactive MIDI Learn simulation toggles, and tempo-synced musical division pickers.
- [x] **Desktop / Process Governance**: Full server shutdown & quit endpoint (`/api/quit`) with responsive controls for system tray & dashboard management.
- [x] **DAW Parity Hardening**: Unified interface implementation fixes across Logic Pro, Pro Tools, Ableton Live, Bitwig, FL Studio, REAPER, Ardour, and Cubase.

## Upcoming Milestones (v3.16.0 - Future)
- [ ] **Native CGO VST3 Deep-Scanning Bridge**: Direct `libvst3` dynamic loading (`dlopen`/`LoadLibrary`) for direct parameter queries without host DAW bridging.
- [ ] **Multi-User Realtime Collaborative Studios**: Operational WebRTC / WebSocket room synchronization for distributed jam sessions and shared parameter locking.
- [ ] **Expanded Mobile Remote**: Touch-screen clip launching grid and mixer fader banks in React Native.
