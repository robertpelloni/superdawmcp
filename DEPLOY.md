# Deployment Guide

## System Requirements
- Go 1.23+
- Python 3.9+ (for Ableton, REAPER, and FL Studio drivers)
- Java 11+ (for Bitwig driver)
- Ableton Live 11/12, REAPER 6+, Bitwig Studio 5+, Ardour 8+, or FL Studio 21+

## Quick Start
1. Clone the repository with submodules:
   ```bash
   git clone --recursive https://github.com/robertpelloni/superdaw-mcp.git
   ```
2. Build the server:
   ```bash
   make build
   ```
3. Install DAW agents:
   - **Ableton Live:** Copy `pkg/agents/ableton/SuperDAW` to your MIDI Remote Scripts folder.
   - **REAPER:** Copy `pkg/agents/reaper/.ReaperOSC` to your REAPER resource path.
   - **Bitwig Studio:** Copy `pkg/agents/bitwig/SuperDAW.bwextension` to your Extensions folder.
   - **FL Studio:** Copy `pkg/agents/flstudio/device_SuperDAW.py` to your FL Studio MIDI Scripts folder.

## MCP Configuration
Add the following to your MCP client configuration (e.g., Claude Desktop config):
```json
{
  "mcpServers": {
    "superdaw": {
      "command": "/path/to/superdaw-mcp",
      "args": []
    }
  }
}
```

## Remote Studio Access (v2.7.0)
SuperDAW now supports remote network control via a TCP Gateway on port 12002.
To connect from a remote machine:
1. Ensure port 12002 is open on the server.
2. Use the 'remote_addr' parameter in the Python SDK:
   ```python
   client = SuperDAWClient(remote_addr="server_ip:12002")
   client.connect()
   ```
