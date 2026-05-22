import sys
import os
import time
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))
from superdaw_client import SuperDAWClient
def run_workflow(client, daw_name):
    print(f"--- Testing {daw_name} ---")
    client.transport_control(playing=False, bpm=124.0, daw=daw_name)
    client.create_track(name="Synth", track_type="midi", daw=daw_name)
    client.set_mixer(track_id="1", volume=0.7, pan=0.2, daw=daw_name)
def main():
    if len(sys.argv) < 2: sys.exit(1)
    client = SuperDAWClient(sys.argv[1])
    client.connect()
    run_workflow(client, "ableton")
    run_workflow(client, "reaper")
    client.disconnect()
if __name__ == "__main__": main()
