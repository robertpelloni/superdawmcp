# SuperDAW-Universal-MCP Ideas

## Aggressive Feature Ideas
- **Neural Plugin Mapper**: Use an LLM to automatically map obscure VST parameter names to common semantic labels (e.g., "Filter Cutoff", "Drive").
- **Multi-DAW Sync (SuperLink)**: Orchestrate sample-accurate sync across different DAWs running on the same machine using `ableton-link`.
- **Natural Language DAW Bridge**: A dedicated sub-module that translates vague musical requests ("Make it sound more underwater") into specific FX chain mutations (low-pass filter, reverb).
- **Universal Preset Browser**: A global database of presets for common VSTs that can be recalled from the MCP server regardless of which DAW is active.
- **Auto-Routing Engine**: Automatically handle sidechaining and complex bus routing across DAWs by abstracting the routing matrix.

## Refactoring / Arch-Porting
- **Native C Bridge**: Port critical timing components to C/Rust to ensure zero-jitter MIDI injection.
- **WASM Agents**: Explore running in-DAW agents as WebAssembly for better cross-platform compatibility and security.
- **GraphQL for Audio**: Expose the entire DAW state as a GraphQL schema for complex queries from frontend clients.
