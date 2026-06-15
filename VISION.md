# SuperDAW-Universal-MCP Vision

## Ultimate Goal
To provide a unified, DAW-agnostic semantic interface for Digital Audio Workstations. SuperDAW-MCP acts as an orchestration layer that allows LLMs and other clients to control multiple DAWs (Ableton Live, REAPER, Ardour, Bitwig, etc.) using a standardized set of tools and commands.

## Core Concepts
1. **Universal Abstraction**: Define a common language for transport, mixer, and MIDI operations that works across all DAWs.
2. **Hybrid Architecture**: Use a high-performance Go daemon to handle MCP and routing, communicating with in-DAW agents via OSC, TCP/JSON, or MIDI.
3. **Multi-Instance Orchestration**: Support concurrent connections to multiple DAW engines (e.g., controlling Ableton and REAPER simultaneously).
4. **Bidirectional State Sync**: Real-time telemetry from DAW agents back to the core daemon and dashboard.
5. **Universal Plugin Control**: Deep scanning and control of VST3 plugins across all supported engines.
