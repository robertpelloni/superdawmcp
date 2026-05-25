import Live
import logging
from _Framework.ControlSurface import ControlSurface
from .pythonosc.osc_server import BlockingOSCUDPServer
from .pythonosc.dispatcher import Dispatcher
from .pythonosc.udp_client import SimpleUDPClient
import threading
try:
    import Queue as queue
except ImportError:
    import queue

# Constants
OSC_LISTEN_PORT = 11000
OSC_RESPONSE_PORT = 11001

class SuperDAW(ControlSurface):
    """
    SuperDAW-MCP Ableton Live Agent
    Bridges Ableton Live API to SuperDAW Go Core via OSC.
    Supports Bidirectional Communication for real-time state sync.
    """
    def __init__(self, c_instance):
        super(SuperDAW, self).__init__(c_instance)
        self.log = logging.getLogger("SuperDAW")
        self._task_queue = queue.Queue()

        # OSC Server for incoming commands
        self._dispatcher = Dispatcher()
        self._dispatcher.map("/superdaw/*", self._enqueue_task)
        self._server = BlockingOSCUDPServer(('127.0.0.1', OSC_LISTEN_PORT), self._dispatcher)

        # OSC Client for outgoing state updates
        self._client = SimpleUDPClient('127.0.0.1', OSC_RESPONSE_PORT)

        self._server_thread = threading.Thread(target=self._server.serve_forever)
        self._server_thread.daemon = True
        self._server_thread.start()

        # Track listeners for volume/pan updates
        self._setup_listeners()

    def _setup_listeners(self):
        # Initial state send
        self._send_full_state()

        # Listen to transport
        self.song().add_is_playing_listener(self._on_playing_changed)
        self.song().add_tempo_listener(self._on_tempo_changed)
        self.song().add_tracks_listener(self._on_arrangement_changed)

    def _on_arrangement_changed(self):
        self._send_arrangement_state()

    def _send_arrangement_state(self):
        arrangement = []
        for track in self.song().tracks:
            clips = []
            if hasattr(track, 'arrangement_clips'):
                for clip in track.arrangement_clips:
                    clips.append({
                        "name": clip.name,
                        "start": clip.start_time,
                        "end": clip.end_time
                    })
            arrangement.append({
                "track": track.name,
                "clips": clips
            })
        import json
        self._client.send_message("/superdaw/state/arrangement", json.dumps(arrangement))

    def _on_playing_changed(self):
        self._client.send_message("/superdaw/state/playing", self.song().is_playing)

    def _on_tempo_changed(self):
        self._client.send_message("/superdaw/state/tempo", self.song().tempo)

    def _send_full_state(self):
        self._on_playing_changed()
        self._on_tempo_changed()

    def _enqueue_task(self, address, *args):
        self._task_queue.put((address, args))

    def update_display(self):
        super(SuperDAW, self).update_display()
        while not self._task_queue.empty():
            try:
                addr, args = self._task_queue.get_nowait()
                self._handle_task(addr, args)
            except queue.Empty:
                break
            except Exception as e:
                self.log.error("Error processing task %s: %s", addr, str(e))

    def _handle_task(self, addr, args):
        if addr == "/superdaw/transport/play":
            if args[0]: self.song().start_playing()
            else: self.song().stop_playing()
        elif addr == "/superdaw/transport/tempo":
            self.song().tempo = float(args[0])
        elif addr == "/superdaw/track/volume":
            idx = int(args[0]); vol = float(args[1])
            if idx < len(self.song().tracks):
                self.song().tracks[idx].mixer_device.volume.value = vol
        elif addr == "/superdaw/track/pan":
            idx = int(args[0]); pan = float(args[1])
            if idx < len(self.song().tracks):
                self.song().tracks[idx].mixer_device.panning.value = pan

    def disconnect(self):
        self.song().remove_is_playing_listener(self._on_playing_changed)
        self.song().remove_tempo_listener(self._on_tempo_changed)
        self._server.shutdown()
        self._server.server_close()
        super(SuperDAW, self).disconnect()
