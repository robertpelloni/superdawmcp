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
- **REAPER**: Ensure the SuperDAW OSC configuration is added in Control/OSC/web settings (Port 8000).
- **Bitwig Studio**: Copy `pkg/agents/bitwig/superdaw-bitwig.bwextension` to your Extensions folder.
- **Ardour**: Enable OSC in preferences. SuperDAW communicates with Ardour's native OSC port (default 3819).

## Running the Server
```bash
./bin/superdaw-mcp
```
