# Changelog

## [2.1.0] - 2024-11-20
### Added
- Logic Pro and Cubase high-level CLI adapters (`scripts/logic_adapter.py`, `scripts/cubase_adapter.py`).
- Interactive SuperDAW Shell (`scripts/superdaw_shell.py`) for real-time multi-DAW orchestration.
- Expanded `install_adapters.sh` with Logic Pro and Cubase support.
- Fully synchronized all 7 DAW drivers into the universal orchestration CLI suite.

## [2.0.0] - 2024-11-20
### Added
- Major release: Universal Multi-DAW Orchestration stable.
- Native Bitwig Studio Agent (Java Extension) implemented.
- Enhanced VST Scanner with path-based vendor heuristics.
- Universal Project Importer (`pkg/engine/importer.go`) for multi-DAW session templates.
- Repository Synchronization Protocol fully integrated into build lifecycle.
- Bidirectional state sync verified for primary drivers.

## [1.9.0] - 2024-11-20
### Added
- Bidirectional State Synchronization for Ableton Live.
- Thread-safe state caching in Ableton and REAPER drivers.
- Real-time transport and tempo feedback from DAW agents to Go core.
