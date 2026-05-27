#!/bin/bash
set -e

echo "------------------------------------------------"
echo "SuperDAW-MCP: Native Agent Installation Utility"
echo "------------------------------------------------"

# Detect OS
OS="unknown"
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "win32" ]]; then
    OS="windows"
fi

echo "Detected OS: $OS"

install_ableton() {
    local ADIR=""
    if [ "$OS" == "macos" ]; then
        ADIR="$HOME/Music/Ableton/User Library/Remote Scripts/SuperDAW"
    else
        ADIR="$APPDATA/Ableton/Live 11/User Library/Remote Scripts/SuperDAW"
    fi
    echo "Installing Ableton Agent to: $ADIR"
    mkdir -p "$ADIR"
    cp -r pkg/agents/ableton/SuperDAW/* "$ADIR/"
}

install_reaper() {
    local RDIR=""
    if [ "$OS" == "macos" ]; then
        RDIR="$HOME/Library/Application Support/REAPER/Scripts/SuperDAW"
    else
        RDIR="$APPDATA/REAPER/Scripts/SuperDAW"
    fi
    echo "Installing REAPER Bridge to: $RDIR"
    mkdir -p "$RDIR"
    cp pkg/agents/reaper/superdaw_bridge.lua "$RDIR/"

    local ODIR=""
    if [ "$OS" == "macos" ]; then
        ODIR="$HOME/Library/Application Support/REAPER/OSC"
    else
        ODIR="$APPDATA/REAPER/OSC"
    fi
    mkdir -p "$ODIR"
    cp pkg/agents/reaper/SuperDAW.ReaperOSC "$ODIR/"
}

install_logic() {
    if [ "$OS" == "macos" ]; then
        local LDIR="$HOME/Music/Audio Music Apps/Custom OSC"
        echo "Installing Logic Pro OSC Mapping to: $LDIR"
        mkdir -p "$LDIR"
        cp pkg/agents/logic/SuperDAW.logic_osc "$LDIR/"
    else
        echo "Skipping Logic Pro (macOS only)."
    fi
}

install_bitwig() {
    local BDIR=""
    if [ "$OS" == "macos" ]; then
        BDIR="$HOME/Documents/Bitwig Studio/Extensions"
    else
        BDIR="$USERPROFILE/Documents/Bitwig Studio/Extensions"
    fi
    echo "Installing Bitwig Extension to: $BDIR"
    mkdir -p "$BDIR"
    # Note: Requires JAR build, assuming exists or copying placeholder
    cp pkg/agents/bitwig/*.bwextension "$BDIR/" 2>/dev/null || echo "Warning: Bitwig .bwextension not found. Build it first."
}

# Execute installations
install_ableton
install_reaper
install_logic
install_bitwig

# Set permissions
chmod +x scripts/*.py 2>/dev/null || true
chmod +x examples/multi_daw_jam/*.py 2>/dev/null || true

echo "------------------------------------------------"
echo "Installation process complete."
echo "------------------------------------------------"
