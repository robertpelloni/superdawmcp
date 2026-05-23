#!/bin/bash
echo "Installing SuperDAW Adapters and Connectors..."

# Ableton Remote Script
ADIR="$HOME/Music/Ableton/User Library/Remote Scripts/SuperDAW"
mkdir -p "$ADIR" && cp -r pkg/agents/ableton/SuperDAW/* "$ADIR/"

# REAPER Bridge
RDIR="$HOME/Library/Application Support/REAPER/Scripts/SuperDAW"
mkdir -p "$RDIR" && cp pkg/agents/reaper/superdaw_bridge.lua "$RDIR/"
mkdir -p "$HOME/Library/Application Support/REAPER/OSC"
cp pkg/agents/reaper/SuperDAW.ReaperOSC "$HOME/Library/Application Support/REAPER/OSC/"

# Make examples executable
chmod +x examples/midi_bridge/midi_to_mcp.py
chmod +x examples/multi_daw_jam/sync_connector.py
chmod +x examples/multi_daw_jam/orchestration.py

echo "Installation complete."
