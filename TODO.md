# SuperDAW-MCP TODO

## Immediate Tasks
- [x] Initialize Go module and install dependencies (`go-osc`, `uuid`).
- [x] Define `DAWDriver` interface in `pkg/daw/driver.go`.
- [x] Implement MCP protocol handlers in `pkg/mcp/protocol.go`.
- [x] Implement Ableton Live OSC client in `pkg/daw/ableton.go`.
- [x] Implement REAPER HTTP/File client in `pkg/daw/reaper.go`.
- [x] Wire up main loop in `cmd/superdaw/main.go`.
- [x] Develop high-level Client Libraries (Go, TypeScript, Python).
- [x] Build automated integration test suite with Mock DAW.
- [x] Create example connectors (MIDI bridge, Python script).
- [x] Implement Ardour driver client in `pkg/daw/ardour.go`.
- [x] Refine REAPER Lua bridge for Pan, Mute, Solo.

## Features to Implement
- [ ] `superdaw_set_mixer`: Set volume, pan, mute across DAWs.
- [ ] `superdaw_write_midi`: Inject MIDI note arrays into clips.
- [ ] `superdaw_get_tracks`: Retrieve track lists and status.
- [ ] `superdaw_transport`: Play, Stop, Record, Tempo control.

## Documentation
- [ ] Detailed `VISION.md` outlining the universal DAW control goal.
- [ ] `MEMORY.md` for architectural observations.
- [ ] `DEPLOY.md` for setup and environment instructions.

## Verification
- [ ] Verify latency is < 2ms for local loopback commands.
- [ ] Test JSON-RPC routing accuracy.
