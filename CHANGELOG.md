# Changelog

## [3.1.0] - 2024-11-21
### Added
- **Plugin Inspector:** Integrated a new UI component in the Web Dashboard for real-time VST3 parameter control.
- **Enhanced VST3 Heuristics:** Expanded automated parameter discovery in `pkg/vst/scanner.go` for LFO, Filters, and effects.
- **Improved Versioning:** Centralized versioning synchronized across VERSION.md and Dashboard.

## [3.0.0] - 2024-11-21
### Added
- **Multi-Instance Connection Manager:** Core engine now supports concurrent connections to multiple instances of the same DAW (e.g., `reaper-1`, `reaper-2`).
- **Heterogeneous DAW Registry:** Dynamic registration and routing for all 8 supported engines.
- **Enhanced Bidirectional Telemetry:** Standardized state feedback (`/superdaw/state/*`) across REAPER, Logic Pro, and FL Studio agents.
- **Specialized SDK Adapters:** High-level Python and TypeScript classes for idiomatic DAW control (`AbletonLive`, `Reaper`, etc.).
- **V3.0.0 Refactor:** Major internal reorganization to decouple driver state from global orchestration.

## [2.9.0] - 2024-11-21
### Synchronized
- **Executive Protocol:** Performed full upstream sync and recursive submodule sanitization.
- **Branch Reconciliation:** Merged all upstream progress into main and aligned feature tracking.
- **Documentation Sync:** Updated Roadmap and TODO based on the latest architectural audit.

## [2.8.0] - 2024-11-21
### Added
- **Logic Pro OSC Feedback:** Real-time state tracking for playing and tempo via port 12101 listener.
- **Python DAW Subclasses:** New specialized classes (AbletonLive, LogicPro, Reaper, etc.) for intuitive SDK usage.
- **Robust Installer:** OS-aware native agent installation script with Windows/macOS path handling.
- **Enhanced Memory & Documentation:** Synthesized architectural learnings and protocol parity reports.

## [2.7.0] - 2024-11-21
### Added
- **TCP MCP Gateway:** Standalone network access on port 12002 for remote control without subprocess constraints.
- **Remote Python SDK:** Support for TCP socket transport in 'SuperDAWClient'.
- **Production Sequence Logic:** Automated multi-DAW orchestration example.
- **MIDI CC Automation:** Universal support for sending MIDI Control Change messages across engines.

## [2.5.0] - 2024-11-21
### Added
- **Universal Feature Analysis:** Integrated 30+ architectural reference repositories to establish a comprehensive feature blueprint (see FEATURES.md).
- **Deep VST Discovery:** Enhanced VST scanner with name-based parameter discovery heuristics and preset metadata extraction.
- **Pro Tools Integration:** Added native OSC driver and CLI adapter for Avid Pro Tools.
- **Universal Studio CLI:** New 'scripts/universal_adapter.py' for studio-wide command execution across all active DAWs.
- **Extended Toolset:** Implemented 'superdaw_set_plugin_parameter', 'superdaw_fire_scene', and full bidirectional monitoring.

## [2.4.0] - 2024-11-21
### Added
- Arrangement Discovery: Real-time track/clip layout broadcasting for Ableton and REAPER.
- JACK Audio Routing: Native system-level audio patching backend.
- Global OSC Gateway: Unified OSC input listener on port 12001.
- Enhanced Web Dashboard: Arrangement View, Virtual MIDI Keyboard, and Studio Blueprint visualization.
- Multi-Language SDKs: Added clients for Ruby, Rust, TS, Sonic Pi.
