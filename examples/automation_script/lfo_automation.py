#!/usr/bin/env python3
"""
SuperDAW High-Frequency LFO Automation
Demonstrates smooth parameter modulation (Volume/Pan) using the MCP server.
Usage: python3 lfo_automation.py --daw ableton --interval 0.1
"""
import sys
import os
import time
import math
import argparse

# Add client to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))
from superdaw_client import SuperDAWClient

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--server", default="../../bin/superdaw-mcp")
    parser.add_argument("--daw", default="ableton")
    parser.add_argument("--interval", type=float, default=0.1)
    parser.add_argument("--duration", type=int, default=10)
    args = parser.parse_args()

    client = SuperDAWClient(args.server)
    client.connect()

    print(f"--- LFO Automation Active ({args.daw}) ---")
    print(f"Interval: {args.interval}s | Target: Track 0 Volume")

    start_time = time.time()
    try:
        while time.time() - start_time < args.duration:
            # Sine wave LFO: 0.2 to 0.8 volume
            t = time.time()
            val = 0.5 + 0.3 * math.sin(2 * math.pi * 0.5 * t) # 0.5 Hz LFO

            client.set_mixer(track_id="0", volume=val, daw=args.daw)

            time.sleep(args.interval)

    except KeyboardInterrupt:
        pass
    finally:
        print("Automation complete.")
        client.disconnect()

if __name__ == "__main__":
    main()
