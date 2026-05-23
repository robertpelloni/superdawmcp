# Deployment Guide

## System Requirements
- Go 1.23+
- Python 3.9+ (for Ableton and REAPER drivers)
- Java 11+ (for Bitwig driver)
- Ableton Live 11/12, REAPER 6+, Bitwig Studio 5+, or Ardour 8+

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
