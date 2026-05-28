from typing import Optional, List, Dict, Any, Callable
from ..client import SuperDAWClient

class DAWAdapter:
    def __init__(self, client: SuperDAWClient, daw_name: str):
        self.client = client
        self.daw_name = daw_name
        self._on_transport_update_cbs = []

        # Subscribe to telemetry from the core
        self.client.on_notification("superdaw/transport_update", self._handle_transport_update)

    def _handle_transport_update(self, params):
        if params.get("daw") == self.daw_name:
            for cb in self._on_transport_update_cbs:
                cb(params.get("playing"), params.get("tempo"))

    def on_transport_update(self, callback: Callable[[bool, float], None]):
        self._on_transport_update_cbs.append(callback)

    def play(self):
        return self.client.transport_control(playing=True, daw=self.daw_name)

    def stop(self):
        return self.client.transport_control(playing=False, daw=self.daw_name)

    def set_volume(self, track_id: str, volume: float):
        return self.client.set_mixer(track_id, volume, daw=self.daw_name)

class AbletonLive(DAWAdapter):
    def __init__(self, client: SuperDAWClient):
        super().__init__(client, "ableton")

    def fire_scene(self, scene_index: int):
        return self.client._call("tools/call", {
            "name": "superdaw_fire_scene",
            "arguments": {"daw": "ableton", "scene_index": scene_index}
        })

class Reaper(DAWAdapter):
    def __init__(self, client: SuperDAWClient):
        super().__init__(client, "reaper")

    def run_action(self, action_id: str):
        return self.client._call("tools/call", {
            "name": "superdaw_custom_command",
            "arguments": {"daw": "reaper", "command": "run_action", "args": {"action_id": action_id}}
        })

class LogicPro(DAWAdapter):
    def __init__(self, client: SuperDAWClient):
        super().__init__(client, "logic")

class Bitwig(DAWAdapter):
    def __init__(self, client: SuperDAWClient):
        super().__init__(client, "bitwig")

class FLStudio(DAWAdapter):
    def __init__(self, client: SuperDAWClient):
        super().__init__(client, "flstudio")

class Cubase(DAWAdapter):
    def __init__(self, client: SuperDAWClient):
        super().__init__(client, "cubase")

class Ardour(DAWAdapter):
    def __init__(self, client: SuperDAWClient):
        super().__init__(client, "ardour")

class ProTools(DAWAdapter):
    def __init__(self, client: SuperDAWClient):
        super().__init__(client, "protools")
