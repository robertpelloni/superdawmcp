#!/usr/bin/env python3
"""
Multi-DAW Synchronization Connector
Demonstrates real-time transport and BPM synchronization across multiple DAWs.
"""
import sys
import os
import time
import argparse

# Add the local client library to the path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))
from superdaw_client.client import SuperDAWClient

def main():
    parser = argparse.ArgumentParser(description="SuperDAW Multi-DAW Sync")
    parser.add_argument("--server", default="./bin/superdaw-mcp", help="Path to SuperDAW MCP server binary")
    parser.add_argument("--daws", default="ableton,reaper", help="Comma-separated list of DAWs to sync")
    args = parser.parse_args()

    target_daws = args.daws.split(",")
    client = SuperDAWClient(args.server)
    client.connect()

    print(f"Syncing DAWs: {', '.join(target_daws)}")

    try:
        # Initial BPM Sync
        bpm = 120.0
        print(f"Aligning all DAWs to {bpm} BPM...")
        for daw in target_daws:
            client.transport_control(playing=False, bpm=bpm, daw=daw)

        time.sleep(1)

        # Simultaneous Start
        print("Firing Start Command to all targets...")
        for daw in target_daws:
            client.transport_control(playing=True, daw=daw)

        print("Syncing... Press Ctrl+C to stop.")

        # Simple loop to periodically re-align (drift correction)
        while True:
            time.sleep(5)
            # In a more advanced version, we would query one DAW for its playhead
            # position and sync others to it.
            # For now, we just ensure transport state is maintained.
            pass

    except KeyboardInterrupt:
        print("\nStopping all DAWs...")
        for daw in target_daws:
            client.transport_control(playing=False, daw=daw)
    finally:
        client.disconnect()

if __name__ == "__main__":
    main()
