# SuperDAW-Universal-MCP Memory

## Architectural Observations
- **Ableton Agent**: Implemented as a MIDI Remote Script. Due to Ableton's restricted Python 3 environment, we vendor the \`pythonosc\` library directly within the package to ensure reliable UDP communication.
- **Protocol Mapping**: The agent acts as a translation layer, mapping unified \`/superdaw/\` OSC commands to the native Live API and broadcasting state changes back with \`/superdaw/state/\` prefixes.
