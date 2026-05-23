#!/usr/bin/env python3
import sys
import os
import argparse

# Add pkg/client/py to sys.path to import superdaw_client
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'pkg', 'client', 'py')))

from superdaw_client.client import SuperDAWClient

def main():
    parser = argparse.ArgumentParser(description="SuperDAW Cubase Adapter CLI")
    parser.add_argument("--server", default="./bin/superdaw-mcp", help="Path to SuperDAW MCP server binary")
    subparsers = parser.add_subparsers(dest="command", help="Command to execute")

    # Transport
    tp = subparsers.add_parser("transport", help="Control transport")
    tp.add_argument("action", choices=["play", "stop"], help="Play or stop")

    # Mixer
    mx = subparsers.add_parser("mixer", help="Control mixer")
    mx.add_argument("track_id", help="Track ID")
    mx.add_argument("--vol", type=float, help="Volume (0.0-1.0)")
    mx.add_argument("--pan", type=float, help="Panning (-1.0 to 1.0)")

    args = parser.parse_args()

    client = SuperDAWClient(args.server)
    client.connect()

    try:
        if args.command == "transport":
            playing = (args.action == "play")
            client.transport_control(playing, daw="cubase")
            print(f"Cubase Transport: {args.action}")

        elif args.command == "mixer":
            client.set_mixer(args.track_id, args.vol if args.vol is not None else 0.8, args.pan if args.pan is not None else 0.0, daw="cubase")
            print(f"Cubase Mixer: Track {args.track_id} (Vol: {args.vol}, Pan: {args.pan})")

        else:
            parser.print_help()
    finally:
        client.disconnect()

if __name__ == "__main__":
    main()
