# Changelog

## [1.8.0] - 2024-11-20
### Added
- Ableton Link Bridge: Unified network clock for all connected DAWs.
- Mobile-Optimized Web Remote for session control on the go.
- Universal Feature Parity Report (`docs/PARITY_REPORT.md`).
- Multi-DAW tempo and transport synchronization via network clock.
- Updated documentation and roadmap for v1.8.0.

## [1.7.0] - 2024-11-20
### Added
- Logic Pro and Cubase high-level CLI adapters (`scripts/logic_adapter.py`, `scripts/cubase_adapter.py`).
- Interactive SuperDAW Shell (`scripts/superdaw_shell.py`) for real-time multi-DAW orchestration.
- Expanded `install_adapters.sh` with Logic Pro and Cubase support.
- Fully synchronized all 7 DAW drivers into the universal orchestration CLI suite.

## [1.6.0] - 2024-11-20
### Added
- Logic Pro and Cubase support via specialized Go drivers.
- WebSocket-powered Web Dashboard for real-time state synchronization.
- Generative MIDI Music Engine with style-based composition heuristics.
- `superdaw_generate_music` tool added to MCP manifest.
- Centralized version governance synchronized to v1.6.0.

## [1.5.0] - 2024-11-20
### Added
- Unified DAW Remote Control GUI (`scripts/superdaw_remote.py`) using tkinter.
- Real-time Mixer, Transport, and Audio Routing controls in the Remote GUI.
- Enhanced discovery logic in `SuperDAWOrchestrator`.
- VST Bridge example in C++ (`examples/vst_bridge/RemoteControlVST.cpp`).
- Updated universal protocol with support for GUI-based orchestration.

## [1.4.0] - 2024-11-20
### Added
- Universal Audio Routing: Virtual patching between DAWs via `superdaw_patch_audio`.
- Generative AI Stem Import: Integrated prompt-based stem generation with `superdaw_import_generative`.
- Enhanced Dashboard: Visual "Audio Routing Matrix" and real-time state visualization.
- Production-ready Spleeter CLI integration in `pkg/engine/stems.go`.
- Multi-DAW session deployment support in `install_adapters.sh`.

## [1.3.0] - 2024-11-20
### Added
- Repository Synchronization Protocol: All 25+ submodules verified and updated.
- Automated installation logic for Bitwig and FL Studio agents in `scripts/install_adapters.sh`.
- High-level `SuperDAWOrchestrator` Python class for multi-DAW session management.
- Python E2E compatibility test suite in `tests/e2e/`.
- Integrated all DAW-specific extensions into the main architecture.

## [1.2.3] - 2024-11-20
### Added
- High-performance C++ (header-only) and standard Java client libraries.
- Comprehensive Stress Test suite for high-frequency DAW control.
- Centralized versioning: Go daemon now reads version from `VERSION.md`.

## [1.2.2] - 2024-11-20
### Added
- High-level Python CLI adapters for Bitwig Studio and FL Studio (`scripts/bitwig_adapter.py`, `scripts/flstudio_adapter.py`).
- Completed CLI adapter suite for all major supported DAWs.

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
