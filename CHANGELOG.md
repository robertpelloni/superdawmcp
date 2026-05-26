# Changelog

## [2.4.0] - 2024-11-21
### Added
- **Arrangement Discovery:** Real-time track/clip layout broadcasting for Ableton and REAPER.
- **JACK Audio Routing:** Native system-level audio patching backend using `jack_connect`.
- **Global OSC Gateway:** Unified OSC input listener on port 12001 for hardware controllers.
- **Enhanced Web Dashboard:** New Arrangement View timeline, Virtual MIDI Keyboard, and Studio Blueprint (signal flow) visualization.
- **Multi-Language SDKs:** Added high-level clients for Ruby, Rust, TypeScript, and Sonic Pi.
- **Performance Suite:** High-frequency LFO automation (10Hz+) and studio stress testing tools.
- **Hardware Support:** First-class integration for Novation Launchpad.
- **Session Persistence:** Save and load complete studio configurations to JSON.
- **Generative AI Job System:** Asynchronous background job monitoring for stem separation and music generation.
- **OBS Integration:** Transparent web overlay endpoint for live streaming.

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
