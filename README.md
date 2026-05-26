# SuperDAW-MCP: The Universal DAW Control Protocol

SuperDAW-MCP is a unified Model Context Protocol (MCP) server that standardizes the control of multiple Digital Audio Workstations (DAWs) through a single semantic interface. It allows LLMs, automation scripts, and hardware controllers to interact with professional audio software using a DAW-agnostic toolset.

## Supported DAWs
- **Ableton Live 11/12** (Python Remote Script)
- **REAPER 6/7** (Lua Bridge + Web API)
- **Bitwig Studio 5** (Java Extension)
- **Logic Pro** (OSC Bridge)
- **Pro Tools** (OSC Bridge)
- **FL Studio 21** (MIDI Script)
- **Cubase 13** (MIDI Remote JS)
- **Ardour 8** (Native OSC)

## Key Features
- **Unified Transport & Mixer Control:** Play, stop, set BPM, and adjust volume/pan across any connected DAW.
- **Arrangement Discovery:** Real-time track and clip layout broadcasting (Ableton/REAPER).
- **JACK Audio Routing:** System-level virtual audio patching between DAW engines.
- **Global OSC Gateway:** Control your entire studio from hardware controllers via port 12001.
- **Web Dashboard:** Interactive timeline visualization, virtual MIDI keyboard, and signal flow monitoring.
- **Multi-Language SDKs:** Ready-to-use libraries for Python, Go, TypeScript, Rust, and more.

## Quick Start
See [DEPLOY.md](DEPLOY.md) for installation and environment setup instructions.

## Architecture
SuperDAW-MCP uses a high-performance Go core daemon to route standardized JSON-RPC commands to lightweight native "Agents" running inside each DAW. This ensures low-latency execution and deep access to proprietary DAW APIs.

---
*Built with Google Jules and the open-source audio community.*
