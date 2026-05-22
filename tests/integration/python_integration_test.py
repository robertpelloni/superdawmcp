import sys
import os
import time
import subprocess
import json

# Add the local client library to the path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))

from superdaw_client import SuperDAWClient

def test_protocol_compliance():
    # Note: This test requires the SuperDAW-MCP binary to be built at bin/superdaw-mcp
    server_path = "./bin/superdaw-mcp"
    if not os.path.exists(server_path):
        print(f"Skipping test: {server_path} not found. Run 'make build' first.")
        return

    client = SuperDAWClient(server_path)

    try:
        client.connect()
        print("Testing mixer control...")
        client.set_mixer(track_id="1", volume=0.7)

        print("Testing track creation...")
        client.create_track(name="Bass", track_type="midi")

        print("Testing transport...")
        client.transport_control(playing=True, bpm=120.0)

        print("Python protocol compliance test: PASSED")

    except Exception as e:
        print(f"Python protocol compliance test: FAILED - {e}")
        sys.exit(1)
    finally:
        client.disconnect()

if __name__ == "__main__":
    test_protocol_compliance()
