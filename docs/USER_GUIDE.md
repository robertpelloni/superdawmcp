# SuperDAW-MCP User Guide

Welcome to SuperDAW-MCP, the Universal DAW Control Protocol. This guide will help you set up the MCP server, connect it to your favorite Digital Audio Workstations (DAWs), and start automating your music production workflows using AI and natural language.

## 1. Introduction

SuperDAW-MCP acts as a universal bridge between Large Language Models (LLMs) and DAWs like Ableton Live, REAPER, Bitwig Studio, Logic Pro, and more. It translates high-level commands (e.g., "Create a bass track and add a compressor") into the specific API calls required by each DAW.

### Architecture Overview

- **Core Daemon:** A high-performance Go application that runs the MCP server.
- **Dynamic Schemas:** Tool definitions (over 1000+ tools from various DAW integrations) are dynamically loaded from `data/schemas/`.
- **DAW Agents:** Lightweight scripts (Lua, Python, Java) that run inside your DAW to execute commands received from the core daemon.

## 2. Installation and Setup

### Prerequisites
- Go 1.23 or higher
- Python 3.9+ (for some DAW agents)
- Supported DAW (Ableton Live 11/12, REAPER 6+, Bitwig Studio 5+, FL Studio 21+, Logic Pro, Cubase 13, Ardour 8)

### Building the Server
Clone the repository and build the core daemon:
```bash
git clone https://github.com/robertpelloni/superdaw-mcp.git
cd superdaw-mcp
make build
```

### Installing DAW Agents

**REAPER:**
Copy `pkg/agents/reaper/superdaw_bridge.lua` to your REAPER Scripts folder and run it within REAPER.

**Ableton Live:**
Copy `pkg/agents/ableton/SuperDAW` to your Ableton MIDI Remote Scripts folder. Select "SuperDAW" as a Control Surface in Ableton's Link/MIDI preferences.

**Bitwig Studio:**
Copy `pkg/agents/bitwig/SuperDAW.bwextension` to your Extensions folder.

## 3. Connecting the MCP Server

Configure your AI assistant (e.g., Claude Desktop, Cursor) to use the SuperDAW-MCP server. Add the following to your MCP configuration file:

```json
{
  "mcpServers": {
    "superdaw": {
      "command": "/path/to/superdaw-mcp/bin/superdaw-mcp",
      "args": []
    }
  }
}
```

## 4. Usage and Workflows

Once connected, your AI assistant will have access to hundreds of DAW control tools.

### Example Natural Language Commands

- **Track Management:** "Create 4 new MIDI tracks and name them Kick, Snare, Bass, and Lead."
- **Mixing:** "Set the volume of the Bass track to -6dB and pan the Snare slightly to the left."
- **Transport:** "Set the project tempo to 128 BPM and start playing."
- **Generative:** "Generate a 4-bar Euclidean rhythm for the Kick track."
- **Routing:** "Route the output of the Lead track into the Reverb bus."

### Dynamic Tools Dispatch
Thanks to the dynamic schema system, specialized commands native to specific DAWs (e.g., triggering scenes in Ableton, executing custom actions in REAPER) are automatically parsed and routed via generic dispatch. You don't need to configure anything extra to use them!

## 5. Troubleshooting

- **Connection Issues:** Ensure the DAW agent is running and listening on the correct ports (default OSC: 11000/11001).
- **Missing Tools:** Check the `data/schemas/` directory to verify that the JSON schemas are present and properly formatted.
- **Latency:** SuperDAW is highly optimized (tool dispatch < 0.2ms). If you experience latency, check your system's audio buffer size or network routing.
