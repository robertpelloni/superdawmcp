import json
import subprocess
import threading
from typing import List, Dict, Any, Optional

class SuperDAWClient:
    def __init__(self, server_path: str, args: Optional[List[str]] = None):
        self.server_path = server_path
        self.args = args or []
        self.process = None
        self.request_id = 1
        self._lock = threading.Lock()
    def connect(self):
        self.process = subprocess.Popen([self.server_path] + self.args, stdin=subprocess.PIPE, stdout=subprocess.PIPE, text=True, bufsize=1)
    def disconnect(self):
        if self.process:
            self.process.stdin.close()
            self.process.wait()
    def _call(self, method: str, params: Dict[str, Any]) -> Any:
        with self._lock:
            req = {"jsonrpc": "2.0", "method": method, "params": params, "id": self.request_id}
            self.request_id += 1
            self.process.stdin.write(json.dumps(req) + "\n")
            line = self.process.stdout.readline()
            return json.loads(line).get("result")
    def set_mixer(self, track_id: str, volume: float, pan: float = 0.0, daw: Optional[str] = None):
        args = {"track_id": track_id, "volume": volume, "pan": pan}
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_set_mixer", "arguments": args})
    def create_track(self, name: str, track_type: str, daw: Optional[str] = None):
        args = {"name": name, "type": track_type}
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_create_track", "arguments": args})
    def transport_control(self, playing: bool, bpm: Optional[float] = None, daw: Optional[str] = None):
        args = {"playing": playing}
        if bpm: args["bpm"] = bpm
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_transport_control", "arguments": args})
