from typing import List, Dict, Any, Optional
from .client import SuperDAWClient
from .adapters import AbletonLive, Reaper, LogicPro, Bitwig, FLStudio, Cubase, Ardour, ProTools

class SuperDAWOrchestrator:
    """
    High-level orchestrator for managing a multi-DAW studio environment.
    Provides synchronized control across multiple heterogeneous engines.
    """
    def __init__(self, client: SuperDAWClient):
        self.client = client
        self.ableton = AbletonLive(client)
        self.reaper = Reaper(client)
        self.logic = LogicPro(client)
        self.bitwig = Bitwig(client)
        self.flstudio = FLStudio(client)
        self.cubase = Cubase(client)
        self.ardour = Ardour(client)
        self.protools = ProTools(client)

        self.adapters = {
            "ableton": self.ableton,
            "reaper": self.reaper,
            "logic": self.logic,
            "bitwig": self.bitwig,
            "flstudio": self.flstudio,
            "cubase": self.cubase,
            "ardour": self.ardour,
            "protools": self.protools
        }

    def play_all(self):
        """Starts transport on all registered and active DAW instances."""
        for name, adapter in self.adapters.items():
            try:
                adapter.play()
            except Exception as e:
                print(f"Warning: Failed to start {name}: {e}")

    def stop_all(self):
        """Stops transport on all registered DAW instances."""
        for name, adapter in self.adapters.items():
            try:
                adapter.stop()
            except Exception as e:
                print(f"Warning: Failed to stop {name}: {e}")

    def sync_tempo(self, bpm: float):
        """Synchronizes tempo across all active DAW engines."""
        for name, adapter in self.adapters.items():
            try:
                self.client.transport_control(playing=None, bpm=bpm, daw=name)
            except Exception as e:
                print(f"Warning: Failed to sync tempo for {name}: {e}")

    def get_studio_status(self) -> Dict[str, Any]:
        """Returns a diagnostic report of the entire studio state."""
        status = {}
        for name in self.adapters:
            try:
                state = self.client.get_transport_state(daw=name)
                tracks = self.client.get_tracks(daw=name)
                status[name] = {
                    "online": True,
                    "playing": state.get("playing", False),
                    "bpm": state.get("bpm", 120.0),
                    "track_count": len(tracks)
                }
            except:
                status[name] = {"online": False}
        return status
