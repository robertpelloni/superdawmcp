# Changelog

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
