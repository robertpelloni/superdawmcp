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
    def write_midi(self, track_id: str, notes: List[Dict[str, Any]], daw: Optional[str] = None):
        args = {"track_id": track_id, "notes": notes}
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_write_midi", "arguments": args})
    def create_track(self, name: str, track_type: str, daw: Optional[str] = None):
        args = {"name": name, "type": track_type}
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_create_track", "arguments": args})
    def transport_control(self, playing: bool, bpm: Optional[float] = None, daw: Optional[str] = None):
        args = {"playing": playing}
        if bpm: args["bpm"] = bpm
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_transport_control", "arguments": args})
    def generate_euclidean(self, track_id: str, hits: int, steps: int, pitch: int, daw: Optional[str] = None):
        args = {"track_id": track_id, "hits": hits, "steps": steps, "pitch": pitch}
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_generate_euclidean", "arguments": args})
    def separate_stems(self, input_path: str, output_dir: str, stems: int = 4):
        args = {"input_path": input_path, "output_dir": output_dir, "stems": stems}
        return self._call("tools/call", {"name": "superdaw_separate_stems", "arguments": args})
    def list_clips(self, track_id: str, daw: Optional[str] = None):
        args = {"track_id": track_id}
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_list_clips", "arguments": args})
    def delete_clip(self, track_id: str, clip_idx: int, daw: Optional[str] = None):
        args = {"track_id": track_id, "clip_idx": clip_idx}
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_delete_clip", "arguments": args})
    def list_plugins(self):
        return self._call("tools/call", {"name": "superdaw_list_plugins", "arguments": {}})
    def get_plugin_params(self, plugin_name: str):
        return self._call("tools/call", {"name": "superdaw_get_plugin_params", "arguments": {"plugin_name": plugin_name}})
    def patch_audio(self, source_daw: str, source_track: str, dest_daw: str, dest_track: str):
        args = {"source_daw": source_daw, "source_track": source_track, "dest_daw": dest_daw, "dest_track": dest_track}
        return self._call("tools/call", {"name": "superdaw_patch_audio", "arguments": args})
    def import_generative(self, prompt: str, target_daw: str):
        args = {"prompt": prompt, "target_daw": target_daw}
        return self._call("tools/call", {"name": "superdaw_import_generative", "arguments": args})
    def generate_music(self, style: str, bars: int = 4, track_id: str = "0", daw: Optional[str] = None):
        args = {"style": style, "bars": bars, "track_id": track_id}
        if daw: args["daw"] = daw
        return self._call("tools/call", {"name": "superdaw_generate_music", "arguments": args})
    def get_tracks(self, daw: Optional[str] = None) -> List[Dict[str, Any]]:
        args = {}
        if daw: args["daw"] = daw
        res = self._call("tools/call", {"name": "superdaw_get_tracks", "arguments": args})
        if isinstance(res, list): return res
        return []
    def list_tools(self) -> List[Dict[str, Any]]:
        res = self._call("tools/list", {})
        if res and "tools" in res: return res["tools"]
        return []
