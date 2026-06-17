# IDEAS.md
## Aggressive Pivots & Refactoring
1. **Audio Routing Engine**: Transition from a basic CLI wrapper for JACK to a dedicated C-based routing engine using `libjack` or `pipewire-native` for sub-millisecond latency and automated jitter compensation.
2. **DAW Agent Porting**: Implement a universal C++ "SuperDAW Agent" that can be compiled as a VST3 plugin itself, allowing it to bypass native API limitations in restricted environments like Ableton.
3. **WebAssembly Dashboard**: Port the entire Dashboard logic to a WASM-based Go implementation for better performance in heavy-orchestration scenarios.
4. **LLM Native Protocol**: Instead of JSON-RPC over TCP, implement a binary protocol optimized for token-efficient tool calls from local LLMs.
5. **Generative Logic**: Integrate a local Stable Audio or AudioLDM engine directly into the Go core for offline stem generation.
