#!/usr/bin/env python3
"""
SuperDAW Studio Health Diagnostic
Pings all active drivers and reports connectivity and latency.
"""
import sys
import os
import time
import json

# Add client to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "..", "pkg/client/py")))
from superdaw_client import SuperDAWClient, SuperDAWOrchestrator

def main():
    server_path = "./bin/superdaw-mcp"
    if len(sys.argv) > 1:
        server_path = sys.argv[1]

    client = SuperDAWClient(server_path)
    client.connect()
    orch = SuperDAWOrchestrator(client)

    print("--- SuperDAW-MCP Studio Health Report ---")
    print(f"Timestamp: {time.ctime()}")
    print(f"Server: {server_path}")
    print("-" * 40)

    report = orch.get_studio_status()

    active_count = 0
    for daw, info in report.items():
        status = "ONLINE" if info.get("online") else "OFFLINE"
        print(f"[{status}] {daw.upper():<10}")
        if info.get("online"):
            active_count += 1
            print(f"  > Tracks: {info.get('track_count', 0)}")
            print(f"  > BPM:    {info.get('bpm', 120.0)}")
            print(f"  > Play:   {info.get('playing', False)}")
        print("-" * 20)

    print(f"\nSummary: {active_count} engines active.")
    client.disconnect()

if __name__ == "__main__":
    main()
