#!/usr/bin/env python3
import sys
import os
import time

# Add pkg/client/py to sys.path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'pkg', 'client', 'py')))

from superdaw_client import SuperDAWClient, LogicPro

def main():
    server_path = "./bin/superdaw-mcp"
    client = SuperDAWClient(server_path)

    print("Logic Pro Bridge: Connecting to SuperDAW core...")
    client.connect()

    # Use the specialized LogicPro subclass
    logic = LogicPro(client)

    try:
        print("Logic Pro: Resetting transport and mixer...")
        logic.stop()
        logic.set_volume("0", 0.75) # Master

        print("Logic Pro: Starting playback at 128 BPM...")
        logic.play(bpm=128.0)

        # Wait for feedback
        time.sleep(2)
        state = logic.get_transport()
        print(f"Logic Pro State: {state}")

        print("Logic Pro: Transitioning to Bridge marker...")
        # Since LogicPro class is basic now, we use the base call for markers
        # In a real build, we'd add fire_scene to LogicPro if it maps to Markers
        client._call("tools/call", {"name": "superdaw_fire_scene", "arguments": {"scene_index": 2, "daw": "logic"}})

        print("Integration successful.")

    except Exception as e:
        print(f"Error: {e}")
    finally:
        client.disconnect()

if __name__ == "__main__":
    main()
