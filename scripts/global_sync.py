#!/usr/bin/env python3
"""
SuperDAW Global Transport Sync
Listens for state updates from a 'Master' DAW and projects them to all other DAWs.
Usage: python3 global_sync.py --master ableton
"""
import sys
import os
import argparse
import time

# Add client to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "pkg", "client", "py")))
from superdaw_client.client import SuperDAWClient

class GlobalSync:
    def __init__(self, server_path, master_daw):
        self.client = SuperDAWClient(server_path)
        self.master_daw = master_daw
        self.targets = ["ableton", "reaper", "bitwig", "ardour"]
        if master_daw in self.targets:
            self.targets.remove(master_daw)

    def start(self):
        self.client.connect()
        print(f"Global Sync Active. Master: {self.master_daw} -> Targets: {self.targets}")

        # Register for transport updates
        self.client.on_notification("superdaw/transport_update", self._on_update)

        try:
            while True:
                time.sleep(1)
        except KeyboardInterrupt:
            self.client.disconnect()

    def _on_update(self, params):
        if params.get("daw") != self.master_daw:
            return

        print(f"[SYNC] Master {self.master_daw} update: {params}")

        # Project to targets
        playing = params.get("playing")
        tempo = params.get("tempo")

        for target in self.targets:
            self.client.transport_control(playing=playing, bpm=tempo, daw=target)

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--server", default="./bin/superdaw-mcp")
    parser.add_argument("--master", default="ableton")
    args = parser.parse_args()

    GlobalSync(args.server, args.master).start()
