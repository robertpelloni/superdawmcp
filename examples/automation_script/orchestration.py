import sys
import os
import time

# Add the local client library to the path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))

from superdaw_client import SuperDAWClient

def main():
    if len(sys.argv) < 2:
        print("Usage: python3 orchestration.py <path_to_superdaw_mcp>")
        sys.exit(1)

    server_path = sys.argv[1]
    client = SuperDAWClient(server_path)

    try:
        print("Connecting to SuperDAW...")
        client.connect()

        # 1. Setup Ableton Live session
        print("Initializing Ableton Live mixer...")
        client.set_mixer(track_id="0", volume=0.8, daw="ableton")
        client.set_mixer(track_id="1", volume=0.5, daw="ableton")

        # 2. Setup REAPER session
        print("Configuring REAPER tracks...")
        client.set_mixer(track_id="1", volume=0.7, pan=-0.5, daw="reaper")

        # 3. Inject MIDI into Ableton
        print("Injecting MIDI sequence into Ableton track 0...")
        notes = [
            {"pitch": 60, "velocity": 100, "start_beat": 0.0, "duration": 0.5},
            {"pitch": 64, "velocity": 90, "start_beat": 1.0, "duration": 0.5},
            {"pitch": 67, "velocity": 110, "start_beat": 2.0, "duration": 1.0},
        ]
        client.write_midi(track_id="0", notes=notes, daw="ableton")

        print("Orchestration complete.")

    except Exception as e:
        print(f"Error during orchestration: {e}")
    finally:
        client.disconnect()
        print("Disconnected.")

if __name__ == "__main__":
    main()
