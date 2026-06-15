# Changelog

## [3.1.0] - 2024-11-21
### Added
- **VST3 Plugin Inspector:** New UI component in the dashboard for deep-scanning and controlling VST3 plugins.
- **Enhanced Parameter Heuristics:** Improved mapping for common plugin parameters (Cutoff, Resonance, etc.) across different vendors.
- **Mobile Remote UX:** Touch-friendly remote interface on `/remote` for transport and mixer control.
- **Logic Pro & Cubase Parity:** Refined drivers for 100% feature parity with Ableton/REAPER.

## [3.0.0] - 2024-11-21
### Added
- **Multi-Instance Connection Manager:** Core engine now supports concurrent connections to multiple instances of the same DAW.
- **Heterogeneous DAW Registry:** Dynamic registration and routing for all 8 supported engines.
- **Enhanced Bidirectional Telemetry:** Standardized state feedback across REAPER, Logic Pro, and FL Studio agents.

## [2.9.0] - 2024-11-21
### Synchronized
- **Executive Protocol:** Performed full upstream sync and recursive submodule sanitization.
- **Branch Reconciliation:** Merged all upstream progress into main and aligned feature tracking.
