# SuperDAW-MCP Deployment Guide

## Prerequisites
- Go 1.21+
- Python 3.10+ (for agents and SDKs)
- Node.js 18+ (for TS SDK)
- JACK Audio Connection Kit (optional, for audio routing)

## Installation
1. Clone the repository and submodules:
   `git clone --recursive https://github.com/robertpelloni/superdawmcp`
2. Build the core daemon:
   `make build`
3. Install DAW agents:
   `./scripts/install_adapters.sh`

## Running
Start the core daemon:
`./bin/superdaw-mcp`

The MCP server will listen on `stdin/stdout`.
The Web Dashboard will be available at `http://127.0.0.1:8081`.
The Mobile Remote will be available at `http://127.0.0.1:8081/remote`.
The TCP Gateway (for Remote SDKs) will listen on port `12002`.
