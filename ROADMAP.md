# SuperDAW-MCP ROADMAP

## Current Milestone (v3.20.0 - Active)
- [x] **VST3 Parameter Tagging & Favorites Engine**: Added `IsFavorite` boolean flag and `Tags []string` to `ParamMetadata` in `pkg/vst/scanner.go`, alongside `SetParamFavorite`, `AddParamTag`, and `GetFavoriteParams` on `*Scanner`. Automatically seeds standard semantic tags (`Gain`, `Master`, `Filter`, `Tone`, `Pitch`, `Tuning`, `Mod`, `Time`, `Space`, `FX`) and default favorites for fast access.
- [x] **Master Console Favorites Filter & Tag Pill Badges**: Integrated a dedicated `⭐ Favorites` filter tab in the Plugin Inspector toolbar in `pkg/ux/dashboard/ui.go`. Parameters display interactive favorite star toggle buttons (`⭐`/`☆`) and semantic tag pills next to their names. Parameter search matches custom tag names.
- [x] **Mobile Touch Remote Favorites Filtering**: Added `isFavorite` and `tags` support to `mobile/src/screens/PluginScreen.tsx` with a tactile `⭐ Favorites` pill toggle button and interactive star toggles on touch parameter cards.
- [x] **Parameter Lock Safeguards & Batch Update Engine (v3.19.0)**: Added `IsLocked` flag to `ParamMetadata`, and implemented thread-safe `SetParamLock`, `GetParamLock`, and `BatchUpdateParams` on `*Scanner` in `pkg/vst/scanner.go`. `BatchUpdateParams` guarantees that locked parameters remain untouched while unlocked values are safely clamped within calibrated `[MinValue, MaxValue]` bounds.
- [x] **Master Console Lock Safeguard Controls (v3.19.0)**: Integrated interactive lock toggle buttons (`🔒`/`🔓`) across every parameter card in `pkg/ux/dashboard/ui.go`. Guarded single reset, macro presets, randomization, batch reset, and continuous A ⟷ B morphing to strictly respect parameter lock states.
- [x] **Mobile Touch Remote Parameter Locking (v3.19.0)**: Added lock toggle buttons, disabled slider interactions for locked parameters, and guard enforcement in `mobile/src/screens/PluginScreen.tsx`.
- [x] **Visual Parameter Transfer Curve Sparklines (v3.18.0)**: Added `GenerateCurvePoints` to `ParamMetadata` in `pkg/vst/scanner.go` with real-time `<canvas>` sparkline curve graphs across parameter cards in `pkg/ux/dashboard/ui.go`.
- [x] **Mobile Touch Remote Snapshot Banks & Continuous Morphing (v3.17.0)**: Fully integrated snapshot banks (Slot A, Slot B, Slot C, Slot D) and real-time A ⟷ B morph slider (0-100%) in `mobile/src/screens/PluginScreen.tsx`.
- [x] **Mobile Precision Parameter Reset (v3.17.0)**: Added tactile individual parameter reset buttons (`↺`) to `mobile/src/screens/PluginScreen.tsx` restoring calibrated defaults.
- [x] **VST3 Parameter Polarity, Slew Rate & Precision Default Reset (v3.16.0)**: Added `IsBipolar`, `SmoothingMs`, `CanAutomate`, and `IsReadOnly` metadata flags with `ClampValue` and `ResetToDefault` methods in `pkg/vst/scanner.go`.
- [x] **Master Console & Mobile Remote Granular Reset (v3.16.0)**: Direct double-click on fader and single-parameter reset buttons (`↺`) restoring exact calibrated defaults without resetting the whole plugin.
- [x] **Visual Polarity & Smoothing Badges (v3.16.0)**: Added `± Bipolar` and `~Smoothing` pills to parameter headers in Master Console UI.
- [x] **Unified Master Console Preset Snapshot Banks & Continuous Morphing (v3.15.0)**: Integrated snapshot slots (Slot A, Slot B, Slot C, Slot D) with active selection highlights and dynamic A ⟷ B continuous parameter morph slider in `pkg/ux/dashboard/ui.go`.
- [x] **Mobile Touch Remote Fine-Nudge Stepping (v3.14.0)**: Added tactile -0.05 and +0.05 parameter adjustment buttons in `mobile/src/screens/PluginScreen.tsx`.
- [x] **Bidirectional Parameter Curve Conversion & Snapping (v3.13.0)**: Precision normalized-to-physical and physical-to-normalized conversion across logarithmic, decibel, S-curve, and stepped snapping.
- [x] **Curve Visualization & Badging (v3.13.0)**: Dedicated mathematical curve icons in Master Console parameter cards.
- [x] **VST3 Parameter Modulation Routing & Automation Engine (v3.12.0)**: Bipolar depth scaling (-1.0 to +1.0) and modulation sources (LFO 1/2, Envelope 1/2, Velocity, Aftertouch).
- [x] **MCP Modulation Protocol Tools (v3.12.0)**: Added `superdaw_set_plugin_param_mod` and `superdaw_get_plugin_param_mod` tools.
- [x] **Master Console & Mobile Remote Mod Controls (v3.12.0)**: Real-time modulation source dropdown and bipolar depth slider beneath each parameter card.
- [x] **Universal Link Synchronization & Studio Clock (v3.11.0)**: Unified network clock arbiter (`LinkBridge`), tempo ramping, beat time calculation, and 4/4 phase quantization alignment.
- [x] **Link MCP Protocol Integration (v3.11.0)**: Added `superdaw_link_sync` and `superdaw_link_status` MCP tools with JSON schemas and gateway routing.
- [x] **Master Studio Console Link Clock HUD (v3.11.0)**: Real-time beat counter, phase bar readout, tempo pulse indicator, and peer network status.
- [x] **Desktop System Tray Controller (`scripts/superdaw_tray.py`) (v3.10.0)**: Platform status monitor, port inspection (HTTP 8081, OSC 12001, TCP 12002), and graceful termination integration via `/api/quit`.
- [x] **Mobile Touch Remote VST3 & Preset Morphing Integration (v3.9.0)**: Mobile remote and React Native `PluginScreen` with snapshot slots, tactile +/- step buttons, and continuous morph slider.
- [x] **Bidirectional Mathematical Curve Calibration (v3.8.0)**: High-precision `NormalizedToPhysical` and `PhysicalToNormalized` mapping for Logarithmic, Decibel ($10^{\text{dB}/20}$), S-Curve, and Linear responses.
- [x] **Consolidated Single-Page Dashboard (v3.7.0)**: Unified Master Studio Console categorizing features by high value (Transport, VST3 Inspector, Mixer, Audio Patching Matrix, Generative DSP, Arrangement Timeline, Virtual MIDI Performance).
- [x] **Advanced VST3 Parameter Discovery (v3.2.0)**: Deep scanning with moduleinfo.json & binary header validation, preset extraction, calibrated parameter ranges, and synchronous API endpoints (`/api/plugins`).
- [x] **Sound Design Quick Macro Knobs & State Tools (v3.3.0 - v3.4.0)**: One-click sound design macros (Brighten, Warmth, Space Out, Punch), musical parameter randomization, and factory reset.
- [x] **Human-Readable Readouts & MIDI Learn (v3.6.0)**: Live display readout formatting, interactive MIDI Learn simulation toggles, and tempo-synced musical division pickers.
- [x] **Desktop / Process Governance (v3.2.0)**: Full server shutdown & quit endpoint (`/api/quit`) with responsive controls for system tray & dashboard management.
- [x] **DAW Parity Hardening**: Unified interface implementation fixes across Logic Pro, Pro Tools, Ableton Live, Bitwig, FL Studio, REAPER, Ardour, and Cubase.

## Upcoming Milestones (v3.20.0 - Future)
- [ ] **Native CGO VST3 Deep-Scanning Bridge**: Direct `libvst3` dynamic loading (`dlopen`/`LoadLibrary`) for direct parameter queries without host DAW bridging.
- [ ] **Multi-User Realtime Collaborative Studios**: Operational WebRTC / WebSocket room synchronization for distributed jam sessions and shared parameter locking.
- [ ] **Expanded Mobile Remote**: Touch-screen clip launching grid and mixer fader banks in React Native.
