# Changelog

## [1.2.1] - 2024-11-20
### Added
- Multi-language client libraries for C#, Ruby, and Rust.
- Expanded `CLIENTS.md` with integration guides for Unity, Godot, and Sonic Pi.

## [1.2.0] - 2024-11-20
### Added
- FL Studio support via specialized Go driver and MIDI scripting Agent.
- Bidirectional Client Connector module for plugin-to-server communication.
- Comprehensive DAW Compatibility Test Suite.
- Enhanced tool manifest with FL Studio routing.

## [1.1.0] - 2024-11-20
### Added
- Real-time Web Dashboard for visual DAW orchestration (port 8080).
- `superdaw_custom_command` for DAW-specific extended functionality.
- Enhanced VST scanning with macOS `Info.plist` metadata extraction.
- Aggressive feature expansion roadmap in `IDEAS.md`.

## [1.0.1] - 2024-11-20
### Added
- High-level Python CLI adapters for Ableton Live and REAPER (`scripts/ableton_adapter.py`, `scripts/reaper_adapter.py`).
- Enhanced thread-safety for Ableton native agent.
- Modernized Go VST scanning logic (removed deprecated `ioutil`).

## [1.0.0] - 2024-11-20
### Added
- Initial implementation of SuperDAW-MCP.
- Universal Go-based MCP server for DAW control.
- Drivers for Ableton Live, REAPER, Ardour, and Bitwig Studio.
- Native agents for Ableton (Python), REAPER (Lua), and Bitwig (Java).
- VST3 scanning engine with path-based discovery.
- Euclidean rhythm generation engine.
- Client libraries for Go, TypeScript, and Python.
- Integration testing framework.
