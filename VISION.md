# SuperDAW-Universal-MCP Vision

## Ultimate Goal
To provide a unified, DAW-agnostic semantic interface for Digital Audio Workstations. SuperDAW-MCP acts as an orchestration layer that allows LLMs and other clients to control multiple DAWs (Ableton Live, REAPER, Ardour, Bitwig, etc.) using a standardized set of tools and commands.

## Core Concepts
1. **Universal Abstraction**: Define a common language for transport, mixer, and MIDI operations that works across all DAWs.
2. **Hybrid Architecture**: Use a high-performance Go daemon to handle MCP and routing, communicating with in-DAW agents via OSC, TCP/JSON, or MIDI.
3. **Low Latency**: Optimize for real-time interaction with network roundtrips under 2ms.
4. **Extensibility**: Easily add new DAWs by implementing the `DAWDriver` interface.
