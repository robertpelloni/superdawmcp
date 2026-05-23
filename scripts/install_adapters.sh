#!/bin/bash
# SuperDAW Client Adapter Installer

# 1. Ableton Live
ABLETON_REMOTE_SCRIPTS_DIR="$HOME/Music/Ableton/User Library/Remote Scripts"
if [ -d "$ABLETON_REMOTE_SCRIPTS_DIR" ]; then
    echo "Installing Ableton Adapter..."
    mkdir -p "$ABLETON_REMOTE_SCRIPTS_DIR/SuperDAW"
    cp -r pkg/agents/ableton/* "$ABLETON_REMOTE_SCRIPTS_DIR/SuperDAW/"
fi

# 2. REAPER
REAPER_SCRIPTS_DIR="$HOME/Library/Application Support/REAPER/Scripts"
if [ -d "$REAPER_SCRIPTS_DIR" ]; then
    echo "Installing REAPER Adapter..."
    mkdir -p "$REAPER_SCRIPTS_DIR/SuperDAW"
    cp pkg/agents/reaper/superdaw_bridge.lua "$REAPER_SCRIPTS_DIR/SuperDAW/"
fi

# 3. Bitwig Studio
BITWIG_EXT_DIR="$HOME/Bitwig Studio/Extensions"
if [ -d "$BITWIG_EXT_DIR" ]; then
    echo "Installing Bitwig Adapter (source)..."
    mkdir -p "$BITWIG_EXT_DIR/SuperDAW"
    cp -r pkg/agents/bitwig/* "$BITWIG_EXT_DIR/SuperDAW/"
fi

echo "Installation Complete."
