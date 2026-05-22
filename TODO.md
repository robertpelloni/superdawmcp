# SuperDAW-MCP TODO

## Immediate Tasks
- [ ] Initialize Go module and install dependencies (`go-osc`, `uuid`).
- [ ] Define `DAWDriver` interface in `pkg/daw/driver.go`.
- [ ] Implement MCP protocol handlers in `pkg/mcp/protocol.go`.
- [ ] Implement Ableton Live OSC client in `pkg/daw/ableton.go`.
- [ ] Implement REAPER HTTP/File client in `pkg/daw/reaper.go`.
- [ ] Wire up main loop in `cmd/superdaw/main.go`.

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
