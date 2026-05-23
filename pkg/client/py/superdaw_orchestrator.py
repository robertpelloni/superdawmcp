import os
import sys
import time
from typing import List, Dict, Optional

# Add local client to path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))
from superdaw_client.client import SuperDAWClient

class SuperDAWOrchestrator:
    """
    High-level orchestrator for managing multi-DAW sessions.
    Handles automatic synchronization and cross-DAW routing logic.
    """
    def __init__(self, server_path: str):
        self.client = SuperDAWClient(server_path)
        self.active_daws: List[str] = []

    def start_session(self, daws: List[str], bpm: float = 120.0):
        self.client.connect()
        self.active_daws = daws
        print(f"Starting Multi-DAW Session: {', '.join(daws)} @ {bpm} BPM")

        for daw in daws:
            self.client.transport_control(playing=False, bpm=bpm, daw=daw)

    def sync_transport(self, playing: bool):
        for daw in self.active_daws:
            self.client.transport_control(playing=playing, daw=daw)

    def create_unified_track(self, name: str, track_type: str = "midi"):
        """Creates the same track across all active DAWs for parallel processing."""
        for daw in self.active_daws:
            self.client.create_track(name=name, track_type=track_type, daw=daw)

    def distribute_notes(self, track_id: str, notes: List[Dict]):
        """Sends different patterns to different DAWs for a layered arrangement."""
        # Simple distribution: send to first DAW for now
        if self.active_daws:
            self.client.write_midi(track_id=track_id, notes=notes, daw=self.active_daws[0])

    def close_session(self):
        self.sync_transport(False)
        self.client.disconnect()
        print("Session closed.")
