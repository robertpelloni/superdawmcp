# Changelog

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
