# SuperDAW-Universal-MCP Memory

## Architectural Observations
- **Ableton Agent**: Implemented as a MIDI Remote Script. Due to Ableton's restricted Python 3 environment, we vendor the \`pythonosc\` library directly within the package to ensure reliable UDP communication.
- **Protocol Mapping**: The agent acts as a translation layer, mapping unified \`/superdaw/\` OSC commands to the native Live API and broadcasting state changes back with \`/superdaw/state/\` prefixes.
- **VST3 Discovery (v3.1.0)**: Enhanced heuristics in `pkg/vst/scanner.go` provide high-confidence mapping for common parameters (LFO, Filters) across generic plugins.
- **Unified Inspector**: The dashboard now features a `Plugin Inspector` that leverages MCP tool `superdaw_get_plugin_params` for real-time control.
- **Bidirectional Feedback (v3.1.0)**: Supports real-time parameter sync for Ableton Live (OSC), REAPER (File-polling), and Bitwig Studio (JSON-RPC broadcast).
- **Generative Theory (v3.1.0)**: The core engine now includes a `theory.go` module supporting Tonic-Predominant-Dominant chord patterns and scale-aware MIDI generation.
