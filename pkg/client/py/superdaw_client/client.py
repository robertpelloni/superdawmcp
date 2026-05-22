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
        self.process = subprocess.Popen(
            [self.server_path] + self.args,
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True,
            bufsize=1
        )

    def disconnect(self):
        if self.process:
            self.process.stdin.close()
            self.process.wait()

    def _call(self, method: str, params: Dict[str, Any]) -> Any:
        with self._lock:
            req = {
                "jsonrpc": "2.0",
                "method": method,
                "params": params,
                "id": self.request_id
            }
            self.request_id += 1

            self.process.stdin.write(json.dumps(req) + "\n")
            line = self.process.stdout.readline()
            if not line:
                raise Exception("Server disconnected unexpectedly")

            res = json.loads(line)
            if "error" in res:
                raise Exception(f"RPC Error: {res['error']['message']} (code {res['error']['code']})")

            return res.get("result")

    def set_mixer(self, track_id: str, volume: float, pan: float = 0.0, daw: Optional[str] = None):
        args = {
            "track_id": track_id,
            "volume": volume,
            "pan": pan
        }
        if daw:
            args["daw"] = daw

        params = {
            "name": "superdaw_set_mixer",
            "arguments": args
        }
        return self._call("tools/call", params)

    def write_midi(self, track_id: str, notes: List[Dict[str, Any]], daw: Optional[str] = None):
        args = {
            "track_id": track_id,
            "notes": notes
        }
        if daw:
            args["daw"] = daw

        params = {
            "name": "superdaw_write_midi",
            "arguments": args
        }
        return self._call("tools/call", params)
