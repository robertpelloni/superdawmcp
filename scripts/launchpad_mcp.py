#!/usr/bin/env python3
"""
Launchpad-to-SuperDAW Bridge
Maps Novation Launchpad grid and side buttons to SuperDAW-MCP commands.
Requires mido: pip install mido python-rtmidi
"""
import sys
import os
import argparse
import time

try:
    import mido
except ImportError:
    print("Error: mido not found. Run 'pip install mido python-rtmidi'")
    sys.exit(1)

# Add the local client library to the path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "pkg", "client", "py")))
from superdaw_client.client import SuperDAWClient

class LaunchpadMCP:
    def __init__(self, server_path, port_name, target_daw):
        self.client = SuperDAWClient(server_path)
        self.port_name = port_name
        self.target_daw = target_daw
        self.out_port = None

    def start(self):
        self.client.connect()
        print(f"Launchpad Bridge Active: '{self.port_name}' -> SuperDAW-MCP ({self.target_daw})")

        try:
            self.out_port = mido.open_output(self.port_name)
            with mido.open_input(self.port_name) as inport:
                self._reset_lights()
                for msg in inport:
                    self._handle_midi(msg)
        except KeyboardInterrupt:
            pass
        finally:
            self.client.disconnect()
            if self.out_port:
                self.out_port.close()

    def _reset_lights(self):
        # Clear grid
        for i in range(128):
            self.out_port.send(mido.Message('note_on', note=i, velocity=0))

    def _handle_midi(self, msg):
        if msg.type == 'note_on' and msg.velocity > 0:
            # Map top row to Transport
            if msg.note == 104: # Up/Play
                self.client.transport_control(playing=True, daw=self.target_daw)
            elif msg.note == 105: # Down/Stop
                self.client.transport_control(playing=False, daw=self.target_daw)

            # Map grid columns to track volumes
            col = msg.note % 16
            row = msg.note // 16
            if col < 8 and row < 8:
                vol = (7 - row) / 7.0
                self.client.set_mixer(track_id=str(col), volume=vol, daw=self.target_daw)
                # Light up the column to feedback volume level
                for r in range(8):
                    vel = 12 if (7-r) <= (7-row) else 0 # Green for active
                    self.out_port.send(mido.Message('note_on', note=r*16 + col, velocity=vel))

if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument("--server", default="./bin/superdaw-mcp")
    parser.add_argument("--port", required=True, help="Launchpad MIDI port name")
    parser.add_argument("--daw", default="ableton")
    args = parser.parse_args()

    LaunchpadMCP(args.server, args.port, args.daw).start()
