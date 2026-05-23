#!/usr/bin/env python3
"""
MIDI-to-MCP Bridge Connector
Translates MIDI CC and Note messages into SuperDAW-MCP commands.
Requires mido and python-rtmidi: pip install mido python-rtmidi
"""
import sys
import os
import argparse

try:
    import mido
except ImportError:
    print("Error: mido not found. Run 'pip install mido python-rtmidi'")
    sys.exit(1)

# Add the local client library to the path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))
from superdaw_client.client import SuperDAWClient

def main():
    parser = argparse.ArgumentParser(description="SuperDAW MIDI-to-MCP Bridge")
    parser.add_argument("--server", default="./bin/superdaw-mcp", help="Path to SuperDAW MCP server binary")
    parser.add_argument("--port", help="MIDI Input Port Name")
    parser.add_argument("--daw", default="ableton", choices=["ableton", "reaper", "ardour", "bitwig"], help="Target DAW")
    args = parser.parse_args()

    if not args.port:
        print("Available MIDI Input Ports:")
        for name in mido.get_input_names():
            print(f" - {name}")
        sys.exit(0)

    client = SuperDAWClient(args.server)
    client.connect()

    print(f"Bridge Active: MIDI Port '{args.port}' -> SuperDAW-MCP ({args.daw})")
    print("Mapping: CC 7 -> Volume, CC 10 -> Pan, Note 60 -> Play, Note 62 -> Stop")

    try:
        with mido.open_input(args.port) as inport:
            for msg in inport:
                if msg.type == 'control_change':
                    if msg.control == 7: # Volume
                        vol = msg.value / 127.0
                        client.set_mixer(track_id="0", volume=vol, daw=args.daw)
                        print(f"CC 7: Volume -> {vol:.2f}")
                    elif msg.control == 10: # Pan
                        pan = (msg.value / 63.5) - 1.0
                        client.set_mixer(track_id="0", volume=0.8, pan=pan, daw=args.daw)
                        print(f"CC 10: Pan -> {pan:.2f}")

                elif msg.type == 'note_on' and msg.velocity > 0:
                    if msg.note == 60: # C3 -> Play
                        client.transport_control(playing=True, daw=args.daw)
                        print("Note 60: Play")
                    elif msg.note == 62: # D3 -> Stop
                        client.transport_control(playing=False, daw=args.daw)
                        print("Note 62: Stop")

    except KeyboardInterrupt:
        print("\nStopping bridge...")
    finally:
        client.disconnect()

if __name__ == "__main__":
    main()
