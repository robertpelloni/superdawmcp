import sys
import os
import time

# Add local SDK to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', '..', 'pkg', 'client', 'py')))

from superdaw_client.client import SuperDAWClient
from superdaw_client.orchestrator import SuperDAWOrchestrator

def production_routine():
    server_path = "./superdaw"
    if not os.path.exists(server_path):
        server_path = "../../bin/superdaw-mcp"

    print(f"Connecting to SuperDAW Studio via {server_path}...")
    client = SuperDAWClient(server_path)
    client.connect()

    orchestrator = SuperDAWOrchestrator(client)

    try:
        print("\n--- Phase 1: Studio Setup ---")
        orchestrator.sync_tempo(128.0)

        print("Creating arrangement structure...")
        orchestrator.ableton.create_track("Drums", "audio")
        orchestrator.reaper.create_track("Bass", "midi")
        orchestrator.logic.create_track("Strings", "audio")

        print("\n--- Phase 2: Mixing and Routing ---")
        orchestrator.ableton.set_volume("0", 0.7)
        orchestrator.reaper.set_volume("0", 0.9)

        # Patching Ableton Drums to REAPER Sidechain
        client.patch_audio("ableton", "out1", "reaper", "in3")

        print("\n--- Phase 3: Synchronized Performance ---")
        print("Starting all engines...")
        orchestrator.play_all()

        time.sleep(5)

        print("Dropping into bridge section...")
        orchestrator.ableton.set_volume("0", 0.3) # Fade drums

        time.sleep(5)

        print("\n--- Phase 4: Session Finalization ---")
        orchestrator.stop_all()

        print("\nFinal Studio Report:")
        print(orchestrator.get_studio_status())

    finally:
        client.disconnect()

if __name__ == "__main__":
    production_routine()
