# Changelog

## [3.2.0] - 2024-11-21
### Added
- **Multi-user Collaborative Sessions**: Introduced WebSocket-based Rooms allowing shared studio state, chat, and event logging.
- **Deep Bidirectional VST Sync**: Real-time parameter feedback for Ableton, Reaper, and Bitwig wired to Plugin Inspector.
- **Music Theory Engine**: New `theory.go` for scale validation and harmonic progression generation.
- **Collaborative Dashboard**: Added Chat, Event Log, and User Presence tracking (active track indicators).
- **Mobile Remote v2.0**: Touch-optimized mixer, transport, and scene launcher with room support.

## [3.1.0] - 2024-11-21
### Added
- **Plugin Inspector:** Integrated a new UI component in the Web Dashboard for real-time VST3 parameter control.
- **Enhanced VST3 Heuristics:** Expanded automated parameter discovery in `pkg/vst/scanner.go` for LFO, Filters, and effects.
- **Improved Versioning:** Centralized versioning synchronized across VERSION.md and Dashboard.
