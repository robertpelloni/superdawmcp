# SuperDAW Architectural Vision

## Core Philosophy
SuperDAW is designed to be the definitive universal interface between AI agents (via the Model Context Protocol - MCP) and professional Digital Audio Workstations (DAWs). Rather than fragmenting DAW control across disparate, incompatible plugins or proprietary scripts, SuperDAW unifies them under a standard model:
- **Universal Control**: One standardized tool interface (`superdaw_*`) controls Ableton Live, REAPER, Bitwig Studio, Logic Pro, Cubase, FL Studio, Pro Tools, and Ardour.
- **Layered Abstraction**: The core engine bridges generic standard transport/mixer controls while preserving DAW-exclusive strengths (e.g. session clips, ReaScript action matrices, Bitwig remote control pages).
- **Consolidated User Experience**: The Master Studio Console consolidates all operational workflows (Transport, VST3 parameter modulation, audio routing, stem separation, arrangement view, virtual performance) into a single, cohesive, responsive interface without fragmented subpages.
- **Bi-Directional Telemetry**: Changes made inside the physical DAW reflect instantly back to the AI context and dashboard via OSC, WebSockets, and TCP JSON-RPC gateways.

## User-Satisfaction Design
- **Zero-Friction Discovery**: Users and agents can inspect available VST3 plugins, their parameter ranges, units, and presets without opening third-party plugin windows.
- **High-Value Hierarchy**: The UI elevates high-frequency operational controls (Transport, Plugin Inspector, Mixer) while keeping deeper topological routing and neural stem workflows immediately accessible.
- **Contextual Tooltips**: Every control is self-documenting via interactive tooltip markers explaining parameter behavior, ranges, and underlying driver protocols.
- **Defensive Sound Design & Non-Destructive Morphing**: Intelligent parameter lock safeguards prevent accidental overwriting of sensitive tuning and volume settings during broad macro adjustments, randomizations, and cross-snapshot morphing.
- **Semantic Organization & Instant Accessibility**: Semantic tagging and custom favoriting enable rapid navigation across hundreds of plugin parameters, ensuring critical controls are instantly reachable on both desktop consoles and touch mobile remotes.
- **Dynamic Macro Modulation & Multi-Target Routing**: Macro parameters modulate multiple plugin parameters simultaneously according to calibrated bipolar weights (-100% to +100%), allowing cohesive single-action sound shaping (e.g. brightening both cutoff and resonance, or punching up envelope punch while trimming attack).
- **Physical MIDI Hardware Mapping & Bidirectional Control**: Parameters seamlessly link to hardware Continuous Controllers (MIDI CC 0-127) with conflict detection and automated unmapping, mirroring tactile studio controller surfaces directly into the unified interface.
- **Non-Destructive Parameter History & Deep Undo/Redo**: Full state-tracked parameter edit journal preserving granular parameter tweaking history across sound design sessions, empowering instant rollbacks (`↶ Undo`) and forward restoration (`↷ Redo`) across web consoles and mobile remotes.
