#!/usr/bin/env python3
import sys
import os
import cmd

# Add pkg/client/py to sys.path to import superdaw_client
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..', 'pkg', 'client', 'py')))

from superdaw_client.client import SuperDAWClient

class SuperDAWShell(cmd.Cmd):
    intro = 'Welcome to the SuperDAW-MCP Shell. Type help or ? to list commands.\n'
    prompt = '(superdaw) '

    def __init__(self, server_path):
        super().__init__()
        self.client = SuperDAWClient(server_path)
        self.client.connect()
        self.active_daw = "ableton"

    def do_daw(self, arg):
        """Switch active DAW: daw ableton"""
        self.active_daw = arg
        print(f"Active DAW: {self.active_daw}")

    def do_play(self, arg):
        """Start transport playback"""
        self.client.transport_control(playing=True, daw=self.active_daw)

    def do_stop(self, arg):
        """Stop transport playback"""
        self.client.transport_control(playing=False, daw=self.active_daw)

    def do_vol(self, arg):
        """Set track volume: vol <track_id> <0.0-1.0>"""
        args = arg.split()
        if len(args) < 2:
            print("Usage: vol <track_id> <volume>")
            return
        self.client.set_mixer(args[0], float(args[1]), daw=self.active_daw)

    def do_exit(self, arg):
        """Exit the shell"""
        self.client.disconnect()
        return True

    def do_EOF(self, arg):
        return self.do_exit(arg)

if __name__ == '__main__':
    server_path = "./bin/superdaw-mcp"
    if len(sys.argv) > 1:
        server_path = sys.argv[1]

    if not os.path.exists(server_path):
        print(f"Error: Server not found at {server_path}")
        sys.exit(1)

    SuperDAWShell(server_path).cmdloop()
