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

        # 1. Global Sync
        bpm = 126.0
        print(f"Setting global tempo to {bpm} BPM...")
        client.transport_control(playing=False, bpm=bpm, daw="ableton")
        client.transport_control(playing=False, bpm=bpm, daw="reaper")

        # 2. Ableton Setup (Drums)
        print("Configuring Ableton (Drums)...")
        client.create_track(name="Drums", track_type="midi", daw="ableton")
        client.set_mixer(track_id="0", volume=0.8, daw="ableton")

        # 3. REAPER Setup (Bass)
        print("Configuring REAPER (Bass)...")
        client.create_track(name="Bass", track_type="midi", daw="reaper")
        client.set_mixer(track_id="0", volume=0.7, pan=-0.3, daw="reaper")

        # 4. Start Jam
        print("Starting synchronized playback...")
        client.transport_control(playing=True, daw="ableton")
        client.transport_control(playing=True, daw="reaper")

        print("Multi-DAW Jam in progress... (5 seconds)")
        time.sleep(5)

        print("Stopping all DAWs.")
        client.transport_control(playing=False, daw="ableton")
        client.transport_control(playing=False, daw="reaper")

    except Exception as e:
        print(f"Error during jam: {e}")
    finally:
        client.disconnect()
        print("Disconnected.")

if __name__ == "__main__":
    main()
