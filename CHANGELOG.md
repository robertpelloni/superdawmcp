# Changelog

## [2.3.0] - 2024-11-20
### Added
- Interactive SuperDAW Shell in `scripts/superdaw_shell.py` for real-time multi-DAW orchestration.
- Logic Pro and Cubase agent installation support in `scripts/install_adapters.sh`.
- Native agent stubs for Logic Pro (`SuperDAW.logic_osc`) and Cubase (`SuperDAW_Cubase.js`).
- Integrated Logic Pro and Cubase drivers into the core orchestration daemon.

## [2.2.0] - 2024-11-20
### Added
- Logic Pro and Cubase high-level Go drivers implemented.
- Bidirectional State Synchronization for REAPER (cached on driver side).
- Native Bitwig Studio Agent (Java) for real-time command processing.
- Universal protocol mapping for Logic OSC and Cubase MIDI Remote.
- Integrated all 7 supported DAWs into the core orchestration daemon.

## [2.1.0] - 2024-11-20
### Added
- Logic Pro and Cubase high-level CLI adapters.
- Interactive SuperDAW Shell for real-time multi-DAW orchestration.
