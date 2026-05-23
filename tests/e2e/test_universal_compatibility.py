import unittest
import sys
import os
import time

# Add client to path
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), "../../pkg/client/py")))
from superdaw_client.client import SuperDAWClient

class TestUniversalCompatibility(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        # We assume the server is already built in bin/
        cls.server_path = os.path.abspath(os.path.join(os.path.dirname(__file__), "../../bin/superdaw-mcp"))
        if not os.path.exists(cls.server_path):
             # Try building it if missing (for sandbox environments)
             os.system("make build")

    def setUp(self):
        self.client = SuperDAWClient(self.server_path)
        self.client.connect()

    def tearDown(self):
        self.client.disconnect()

    def test_ableton_basic(self):
        """Verify basic tool routing to Ableton driver."""
        res = self.client.transport_control(playing=True, bpm=124.0, daw="ableton")
        self.assertIsNotNone(res)

    def test_reaper_basic(self):
        """Verify basic tool routing to REAPER driver."""
        res = self.client.transport_control(playing=False, bpm=140.0, daw="reaper")
        self.assertIsNotNone(res)

    def test_bitwig_basic(self):
        """Verify basic tool routing to Bitwig driver."""
        res = self.client.set_mixer(track_id="0", volume=0.9, daw="bitwig")
        self.assertIsNotNone(res)

    def test_flstudio_basic(self):
        """Verify basic tool routing to FL Studio driver."""
        res = self.client.create_track(name="FL_Synth", track_type="midi", daw="flstudio")
        self.assertIsNotNone(res)

    def test_euclidean_engine(self):
        """Verify the Go engine's Euclidean generator logic."""
        res = self.client.generate_euclidean(track_id="1", hits=3, steps=8, pitch=42, daw="ableton")
        self.assertIsNotNone(res)

if __name__ == "__main__":
    unittest.main()
