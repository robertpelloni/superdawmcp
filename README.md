# SuperDAW-MCP

A Model Context Protocol (MCP) server that exposes DAW (digital audio workstation) control through one semantic interface, so LLMs, automation scripts, and hardware controllers can drive different DAWs with a DAW-agnostic toolset: transport, mixer, tempo, arrangement discovery, and routing.

## Status

**Active hobby project, no stable release.** Personal studio-automation project by Robert Pelloni (psytrance production + MCP tooling). Version in tree: v3.1.0 (`VERSION.md`). Last commit: 2026-11-06 (submodule re-pin); last substantive feature work: v3.1.0 VST3 Inspector, merged 2026-06-15.

## What's in here

- **Go server core** (`go.mod`, `cmd/`, `pkg/`) with Python client/tooling alongside.
- **DAW bridges** — per-DAW adapters of varying maturity. Ableton Live (Python remote script) and REAPER (Lua/web) are the most developed. The full target list from the design docs: Ableton Live 11/12, REAPER 6/7, Bitwig Studio 5, Logic Pro (OSC), Pro Tools (OSC), FL Studio 21 (MIDI script), Cubase 13 (MIDI remote), Ardour 8 (OSC). Treat non-Ableton/REAPER bridges as in-progress unless `FEATURES.md` / `CLIENTS.md` say otherwise.
- **Psytrance production scripts** (`build_psytrance*.py`, `make_psytrance.py`, `automate_surge_xt.py`, arrangement enhancers) — the working tooling this project grew out of.
- `examples/`, `mobile/`, `docs/`, `tests/`, plus many DAW-MCP submodules (see `SUBMODULES.md` / `SUBMODULE_AUDIT.md`).

## Quick start

See `DEPLOY.md` for installation and environment setup. Python SDK sketch (from the project docs):

```python
from superdaw_client.client import SuperDAWClient
from superdaw_client.orchestrator import SuperDAWOrchestrator

client = SuperDAWClient("./superdaw")
client.connect()
studio = SuperDAWOrchestrator(client)
studio.sync_tempo(128.0)
studio.play_all()
```

## Docs in this repo

`FEATURES.md`, `CLIENTS.md`, `VISION.md`, `ROADMAP.md`, `TODO.md`, `CHANGELOG.md`, `HANDOFF.md`, `SUBMODULES.md`.
