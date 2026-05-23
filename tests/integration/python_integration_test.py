import sys
import os
import time

# Add the local client library to the path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))

from superdaw_client import SuperDAWClient

def test_protocol_compliance():
    # Requires 'make build' first
    server_path = "./bin/superdaw-mcp"
    if not os.path.exists(server_path):
        print(f"Skipping test: {server_path} not found.")
        return

    client = SuperDAWClient(server_path)

    try:
        client.connect()
        print("Testing mixer control via Python...")
        client.set_mixer(track_id="1", volume=0.7)

        print("Testing track creation via Python...")
        client.create_track(name="Bass", track_type="midi")

        print("Testing transport via Python...")
        client.transport_control(playing=True, bpm=120.0)

        print("Python protocol compliance test: PASSED")

    except Exception as e:
        print(f"Python protocol compliance test: FAILED - {e}")
        sys.exit(1)
    finally:
        client.disconnect()

if __name__ == "__main__":
    test_protocol_compliance()
