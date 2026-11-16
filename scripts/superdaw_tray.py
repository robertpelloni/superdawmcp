#!/usr/bin/env python3
"""
SuperDAW Desktop System Tray Controller
Provides background process monitoring, web console launcher,
studio status inspection, and graceful server termination.
"""
import sys
import os
import time
import socket
import urllib.request
import urllib.error
import subprocess
import argparse

class SuperDAWTrayController:
    def __init__(self, host="127.0.0.1", http_port=8081, osc_port=12001, tcp_port=12002):
        self.host = host
        self.http_port = http_port
        self.osc_port = osc_port
        self.tcp_port = tcp_port
        self.base_url = f"http://{self.host}:{self.http_port}"

    def is_port_open(self, port, timeout=1.0):
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        sock.settimeout(timeout)
        try:
            res = sock.connect_ex((self.host, port))
            sock.close()
            return res == 0
        except Exception:
            return False

    def get_status(self):
        http_ok = self.is_port_open(self.http_port)
        tcp_ok = self.is_port_open(self.tcp_port)

        status = {
            "http_online": http_ok,
            "tcp_gateway_online": tcp_ok,
            "server_url": self.base_url,
            "mobile_url": f"{self.base_url}/remote",
        }
        return status

    def print_status(self):
        s = self.get_status()
        print("=" * 45)
        print("  SuperDAW-MCP System Tray Status Monitor")
        print("=" * 45)
        http_badge = "[ONLINE]" if s["http_online"] else "[OFFLINE]"
        tcp_badge = "[ONLINE]" if s["tcp_gateway_online"] else "[OFFLINE]"
        print(f"Master Console (HTTP 8081):   {http_badge} {s['server_url']}")
        print(f"Mobile Remote (/remote):       {http_badge} {s['mobile_url']}")
        print(f"TCP MCP Gateway (Port 12002):  {tcp_badge}")
        print("=" * 45)

    def quit_server(self):
        """Sends graceful shutdown request to /api/quit"""
        url = f"{self.base_url}/api/quit"
        print(f"Sending graceful quit signal to {url}...")
        try:
            req = urllib.request.Request(url, data=b"{}", headers={"Content-Type": "application/json"})
            with urllib.request.urlopen(req, timeout=3.0) as resp:
                print(f"Server response: {resp.status} {resp.reason}")
                return True
        except urllib.error.URLError as e:
            print(f"Shutdown signal failed (server may already be terminated): {e}")
            return False

def main():
    parser = argparse.ArgumentParser(description="SuperDAW System Tray Controller")
    parser.add_argument("--check", action="store_true", help="Check status and port availability")
    parser.add_argument("--quit", action="store_true", help="Send graceful termination signal to daemon")
    args = parser.parse_args()

    controller = SuperDAWTrayController()

    if args.quit:
        controller.quit_server()
        sys.exit(0)

    controller.print_status()

if __name__ == "__main__":
    main()
