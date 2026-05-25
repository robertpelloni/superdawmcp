#!/usr/bin/env python3
"""
Cross-DAW Automation Example
1. Create a MIDI drum pattern in Ableton.
2. Route audio to REAPER for specialized mixing.
3. Synchronize both DAWs and play.
"""
import sys
import os
import time

# Add client to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))
from superdaw_client import SuperDAWClient, SuperDAWOrchestrator

def main():
    server_path = "../../bin/superdaw-mcp"
    client = SuperDAWClient(server_path)
    client.connect()
    orch = SuperDAWOrchestrator(client)

    print("--- Starting Studio Automation Flow ---")

    # 1. Setup Ableton
    print("Initializing Ableton (MIDI Source)...")
    client.create_track("Drums", "midi", daw="ableton")
    client.generate_euclidean("0", hits=5, steps=16, pitch=36, daw="ableton") # Kick

    # 2. Setup REAPER
    print("Initializing REAPER (Processing Hub)...")
    client.create_track("Ableton_Return", "audio", daw="reaper")

    # 3. Patch Audio
    print("Establishing Virtual Audio Patch...")
    client.patch_audio(source_daw="ableton", source_track="Drums",
                       dest_daw="reaper", dest_track="Ableton_Return")

    # 4. Synchronized Playback
    print("Synchronizing and Playing...")
    orch.studio_play(bpm=128.0)

    print("Automation complete. Studio is now live.")
    time.sleep(2)
    client.disconnect()

if __name__ == "__main__":
    main()
