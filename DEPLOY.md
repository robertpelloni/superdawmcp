# SuperDAW-Universal-MCP Deployment

## Prerequisites
- Go 1.24+
- Python 3.x (for Ableton)
- Java 11+ (for Bitwig)

## Installation

### 1. Build the Server
```bash
make build
```

### 2. Install DAW Adapters
Run the automated installation script:
```bash
./scripts/install_adapters.sh
```

#### Manual Steps per DAW:
- **Ableton Live**: Ensure "SuperDAW" is selected as a Control Surface in Link/MIDI preferences.
- **REAPER**: Copy `pkg/agents/reaper/SuperDAW.ReaperOSC` to your REAPER/OSC folder and select it in Control/OSC/web settings (Port 8000).
- **Bitwig Studio**: Copy `pkg/agents/bitwig/superdaw-bitwig.bwextension` to your Extensions folder.
- **Ardour**: Copy `pkg/agents/ardour/superdaw.osc` to your Ardour/osc folder and enable OSC in preferences.

## Running the Server
```bash
./bin/superdaw-mcp
```
