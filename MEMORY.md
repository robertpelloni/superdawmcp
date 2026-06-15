# SuperDAW-Universal-MCP Memory

## Architectural Observations
- **Ableton Agent**: Implemented as a MIDI Remote Script. Due to Ableton's restricted Python 3 environment, we vendor the `pythonosc` library directly within the package to ensure reliable UDP communication.
- **Protocol Mapping**: The agent acts as a translation layer, mapping unified `/superdaw/` OSC commands to the native Live API and broadcasting state changes back with `/superdaw/state/` prefixes.
- **Connection Manager**: Decouples driver state from global orchestration, allowing dynamic registration and routing for multiple DAW instances.
- **Plugin Inspector (v3.1.0)**: Uses a centralized VST3 scanner with heuristic-based parameter mapping to provide a unified UI for plugin control.
- **Dashboard & CommandBus**: Internal channel used to execute tools from the Web UI/Remote without corrupting stdout, maintaining MCP protocol integrity.
