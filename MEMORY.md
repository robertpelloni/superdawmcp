# Codebase Memory & Architectural Observations

- **Unified Dashboard Architecture**: The web console lives in `pkg/ux/dashboard/` and is served on port 8081. It avoids template variable collision by escaping `%` signs (`%%`).
- **Synchronous Plugin Querying**: In addition to asynchronous MCP bus dispatch (`/api/call`), the dashboard provides synchronous REST query endpoints (`/api/plugins?name=...&deep=true`) for instant parameter population.
- **Process Management**: Graceful server termination is exposed via `/api/quit`, triggering clean daemon teardown across HTTP, WebSocket, OSC, and TCP socket listeners.
- **Root Directory Cleanliness**: Do not leave runnable `main` package files at the repository root (`test_*.go`), as they break package builds (`go test ./...`). Auxiliary scratch scripts must be archived in `archive/`.
- **DAW Driver Completeness**: All drivers must implement every method on `daw.DAWDriver` (`SendCC`, `SetPluginParameter`, etc.) even if stubbed for uncooperative DAWs, avoiding duplicate declaration collisions.
- **Version Governance**: `VERSION.md` is the single source of truth for versioning (currently `v3.9.0`), read dynamically by both documentation and the Go binary at runtime.
- **Snapshot Banks & Continuous Morphing**: Managed across web Master Studio Console (`pkg/ux/dashboard/ui.go`), web mobile remote (`pkg/ux/dashboard/remote.go`), and React Native remote (`mobile/src/screens/PluginScreen.tsx`) using slots A, B, C, D and linear interpolation.
