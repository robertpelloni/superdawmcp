# SuperDAW-Universal-MCP Deployment

## Prerequisites
- Go 1.24+
- Targeted DAWs (Ableton Live, REAPER, etc.)
- Python 3.x (for in-DAW agents like AbletonOSC)

## Setup Instructions

### 1. Clone the Repository
```bash
git clone https://github.com/robertpelloni/superdaw-mcp
cd superdaw-mcp
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Build the Server
```bash
go build -o bin/superdaw-mcp cmd/superdaw/main.go
```

### 4. Configure DAWs
- **Ableton Live**: Install the `AbletonOSC` MIDI Remote Script from `third_party/ableton-osc`.
- **REAPER**: Enable the Web Control surface on port 8080 and ensure OSC is listening on port 8000.

### 5. Run the MCP Server
```bash
./bin/superdaw-mcp
```

## Environment Variables
- `REAPER_WEB_PORT`: Default 8080
- `ABLETON_OSC_PORT`: Default 11000
- `VST3_PATH`: Path to scan for VST3 plugins
