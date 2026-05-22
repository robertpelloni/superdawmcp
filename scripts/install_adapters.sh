#!/bin/bash
# SuperDAW Client Adapter Installer

# 1. Ableton Live
ABLETON_REMOTE_SCRIPTS_DIR="$HOME/Music/Ableton/User Library/Remote Scripts"
if [ -d "$ABLETON_REMOTE_SCRIPTS_DIR" ]; then
    echo "Installing Ableton Adapter..."
    mkdir -p "$ABLETON_REMOTE_SCRIPTS_DIR/SuperDAW"
    cp -r pkg/agents/ableton/* "$ABLETON_REMOTE_SCRIPTS_DIR/SuperDAW/"
else
    echo "Ableton Remote Scripts directory not found. Skipping."
fi

# 2. REAPER
REAPER_SCRIPTS_DIR="$HOME/Library/Application Support/REAPER/Scripts"
if [ -d "$REAPER_SCRIPTS_DIR" ]; then
    echo "Installing REAPER Adapter..."
    mkdir -p "$REAPER_SCRIPTS_DIR/SuperDAW"
    cp pkg/agents/reaper/superdaw_bridge.lua "$REAPER_SCRIPTS_DIR/SuperDAW/"
else
    echo "REAPER Scripts directory not found. Skipping."
fi

echo "Installation Complete."
