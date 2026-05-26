#!/usr/bin/env python3
import sys
import os
import argparse
import json

# Add pkg/client/py to sys.path to import superdaw_client
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'pkg', 'client', 'py')))

from superdaw_client.client import SuperDAWClient

class UniversalStudioAdapter:
    def __init__(self, server_path):
        self.client = SuperDAWClient(server_path)
        self.client.connect()
        self.daws = ["ableton", "reaper", "bitwig", "flstudio", "logic", "cubase", "protools", "ardour"]

    def close(self):
        self.client.disconnect()

    def run_studio_wide(self, operation, **kwargs):
        print(f"Executing Studio-Wide Operation: {operation.upper()}")
        results = {}
        for daw in self.daws:
            try:
                if operation == "play":
                    self.client.transport_control(playing=True, daw=daw)
                elif operation == "stop":
                    self.client.transport_control(playing=False, daw=daw)
                elif operation == "mute_all":
                    tracks = self.client.get_tracks(daw=daw)
                    for t in tracks:
                        # Assuming a set_mixer or custom command for mute
                        pass
                elif operation == "sync_bpm":
                    bpm = kwargs.get("bpm", 120.0)
                    self.client.transport_control(playing=None, bpm=bpm, daw=daw)
                results[daw] = "Success"
            except Exception as e:
                results[daw] = f"Error: {str(e)}"
        return results

def main():
    parser = argparse.ArgumentParser(description="SuperDAW Universal Studio CLI")
    parser.add_argument("--server", default="./bin/superdaw-mcp", help="Path to SuperDAW MCP server binary")
    parser.add_argument("command", choices=["play-all", "stop-all", "sync-bpm", "status"], help="Command to execute")
    parser.add_argument("--bpm", type=float, default=120.0, help="BPM for sync")

    args = parser.parse_args()
    adapter = UniversalStudioAdapter(args.server)

    try:
        if args.command == "play-all":
            res = adapter.run_studio_wide("play")
            print(json.dumps(res, indent=2))
        elif args.command == "stop-all":
            res = adapter.run_studio_wide("stop")
            print(json.dumps(res, indent=2))
        elif args.command == "sync-bpm":
            res = adapter.run_studio_wide("sync_bpm", bpm=args.bpm)
            print(json.dumps(res, indent=2))
        elif args.command == "status":
            # Gather status from all DAWs
            status = {}
            for daw in adapter.daws:
                try:
                    state = adapter.client.get_transport_state(daw=daw)
                    status[daw] = state
                except:
                    status[daw] = "Offline"
            print(json.dumps(status, indent=2))
    finally:
        adapter.close()

if __name__ == "__main__":
    main()
