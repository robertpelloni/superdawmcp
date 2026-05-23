#!/bin/bash
echo "Installing SuperDAW Adapters..."
# Ableton
ADIR="$HOME/Music/Ableton/User Library/Remote Scripts/SuperDAW"
mkdir -p "$ADIR" && cp -r pkg/agents/ableton/SuperDAW/* "$ADIR/"
# REAPER
RDIR="$HOME/Library/Application Support/REAPER/Scripts/SuperDAW"
mkdir -p "$RDIR" && cp pkg/agents/reaper/superdaw_bridge.lua "$RDIR/"
mkdir -p "$HOME/Library/Application Support/REAPER/OSC"
cp pkg/agents/reaper/SuperDAW.ReaperOSC "$HOME/Library/Application Support/REAPER/OSC/"
echo "Adapters installed."
