# SuperDAW-MCP

A Model Context Protocol (MCP) server that exposes DAW (digital audio workstation) control through one semantic interface, so LLMs, automation scripts, and hardware controllers can drive different DAWs with a DAW-agnostic toolset: transport, mixer, tempo, arrangement discovery, and routing.

## What's implemented vs. targeted

The server core is **Go** (`go.mod`, `cmd/`, `pkg/`) with Python client/tooling alongside. DAW bridges are per-DAW adapters of varying maturity — Ableton Live (Python remote script) and REAPER (Lua/web) are the most developed; the full target list from the project's design docs is: Ableton Live 11/12, REAPER 6/7, Bitwig Studio 5, Logic Pro (OSC), Pro Tools (OSC), FL Studio 21 (MIDI script), Cubase 13 (MIDI remote), Ardour 8 (OSC). Treat non-Ableton/REAPER bridges as in-progress unless `FEATURES.md` / `CLIENTS.md` say otherwise.

Also in this tree: a set of Python psytrance-production scripts (`build_psytrance*.py`, `make_psytrance.py`, `automate_surge_xt.py`, arrangement enhancers) — the working tooling this project grew out of — plus `examples/`, `mobile/`, and `docs/`.

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

## Status

Personal studio-automation project (Robert Pelloni — psytrance production + MCP tooling). Actively hacked on; no stable release. Docs: `FEATURES.md`, `CLIENTS.md`, `VISION.md`, `ROADMAP.md`, `TODO.md`, `SUBMODULES.md`.
