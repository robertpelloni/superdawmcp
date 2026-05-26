#!/usr/bin/env python3
import sys
import os
import time

# Add pkg/client/py to sys.path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..', 'pkg', 'client', 'py')))

from superdaw_client.client import SuperDAWClient

def main():
    # Connect to a remote SuperDAW server over the network
    # Assumes server is running with: go run cmd/superdaw/main.go
    remote_server = "127.0.0.1:12002"
    client = SuperDAWClient(remote_addr=remote_server)

    print(f"Connecting to remote SuperDAW at {remote_server}...")
    try:
        client.connect()
        print("Connected! Executing production sequence...")

        # 1. Sync all DAWs to production tempo
        client.transport_control(playing=True, bpm=124.0, daw="ableton")
        client.transport_control(playing=True, bpm=124.0, daw="reaper")

        # 2. Automated Mixer fade-in
        for i in range(10):
            vol = i / 10.0
            client.set_mixer(track_id="0", volume=vol, daw="ableton")
            time.sleep(0.5)

        print("Production sequence complete.")

    except Exception as e:
        print(f"Error: {e}")
    finally:
        client.disconnect()

if __name__ == "__main__":
    main()
