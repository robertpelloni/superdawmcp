# SuperDAW-Universal-MCP Memory

## Architectural Observations
- **Go Daemon**: Selected for its low-latency concurrency and single-binary distribution. It serves as the JSON-RPC server for MCP.
- **Protocol Mapping**:
  - **Ableton**: Uses OSC via `AbletonOSC` and `pylive` frameworks.
  - **REAPER**: Uses a hybrid of OSC and the Web API (`wwr`). Reference to `total-reaper-mcp`'s file bridge for future implementation of complex ReaScript calls.
  - **Ardour**: Uses native OSC paths (`/strip`, `/transport`).
- **VST Metadata**: A dedicated scanning module extracts parameter schemas to facilitate LLM-driven plugin control without manual mapping.

## Design Preferences
- **JSON-RPC 2.0**: Standard for MCP.
- **Thread-Safety**: VST scanner and state tracking use RWMutex for safe concurrent access.
- **Modular Drivers**: Each DAW has its own driver implementation to isolate protocol-specific logic.

## Implementation Details
- `DAWDriver` interface in `pkg/daw/driver.go` is the core contract.
- OSC communication uses `github.com/hypebeast/go-osc/osc`.
- Internal routing in `cmd/superdaw/main.go` currently defaults to Ableton but is designed to support multiple active targets.
