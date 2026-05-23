# Changelog

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
