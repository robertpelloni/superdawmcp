import time
from typing import List, Dict, Any, Optional
from .client import SuperDAWClient

class SuperDAWOrchestrator:
    """
    High-level session management for multi-DAW environments.
    """
    def __init__(self, client: SuperDAWClient):
        self.client = client
        self.daws = ["ableton", "reaper", "bitwig", "ardour", "flstudio", "logic", "cubase"]

    def studio_reset(self):
        """Reset mixer volumes and stop transport on all engines."""
        for daw in self.daws:
            try:
                self.client.transport_control(playing=False, daw=daw)
                self.client.set_mixer(track_id="0", volume=0.8, pan=0.0, daw=daw)
            except:
                continue

    def studio_play(self, bpm: float = 120.0):
        """Start synchronized playback across all engines."""
        for daw in self.daws:
            try:
                self.client.transport_control(playing=True, bpm=bpm, daw=daw)
            except:
                continue

    def studio_sync_tempo(self, bpm: float):
        """Force a global tempo sync."""
        for daw in self.daws:
            try:
                self.client.transport_control(bpm=bpm, daw=daw)
            except:
                continue

    def studio_report(self) -> Dict[str, Any]:
        """Query state from all active drivers."""
        report = {}
        for daw in self.daws:
            try:
                tracks = self.client.get_tracks(daw=daw)
                transport = self.client.get_transport_state(daw=daw)
                report[daw] = {
                    "active": True,
                    "track_count": len(tracks),
                    "transport": transport
                }
            except:
                report[daw] = {"active": False}
        return report

    def wait_for_daw(self, daw: str, timeout: int = 10):
        """Block until a specific DAW engine reports as active."""
        start = time.time()
        while time.time() - start < timeout:
            try:
                self.client.get_transport_state(daw=daw)
                return True
            except:
                time.sleep(1)
        return False
