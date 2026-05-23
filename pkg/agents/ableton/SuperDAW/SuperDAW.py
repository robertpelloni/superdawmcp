import Live
import logging
from _Framework.ControlSurface import ControlSurface
from .pythonosc.osc_server import BlockingOSCUDPServer
from .pythonosc.dispatcher import Dispatcher
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
    Thread-safe implementation using a task queue.
    """
    def __init__(self, c_instance):
        super(SuperDAW, self).__init__(c_instance)
        self.log = logging.getLogger("SuperDAW")
        self._task_queue = queue.Queue()

        self._dispatcher = Dispatcher()
        self._dispatcher.map("/superdaw/*", self._enqueue_task)

        self._server = BlockingOSCUDPServer(('127.0.0.1', OSC_LISTEN_PORT), self._dispatcher)

        self._server_thread = threading.Thread(target=self._server.serve_forever)
        self._server_thread.daemon = True
        self._server_thread.start()

    def _enqueue_task(self, address, *args):
        """Add OSC message to queue to be processed on main thread."""
        self._task_queue.put((address, args))

    def update_display(self):
        """
        Called by Live every 100ms on the main thread.
        Process all queued OSC tasks here.
        """
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
        elif addr == "/superdaw/track/create":
            name = args[0]; t = args[1]
            if t == "audio": self.song().create_audio_track()
            else: self.song().create_midi_track()
            self.song().tracks[-1].name = name
        elif addr == "/superdaw/clip/write":
            t_idx = int(args[0]); c_idx = int(args[1]); notes = args[2:]
            if t_idx < len(self.song().tracks):
                track = self.song().tracks[t_idx]
                if not track.has_midi_input: return
                slot = track.clip_slots[c_idx]
                if not slot.has_clip: slot.create_clip(4.0)
                clip = slot.clip; clip.remove_notes(0, 0, 127, 127)
                new_notes = []
                for i in range(0, len(notes), 4):
                    new_notes.append(Live.Clip.MidiNoteSpecification(
                        pitch=int(notes[i]),
                        velocity=int(notes[i+1]),
                        start_time=float(notes[i+2]),
                        duration=float(notes[i+3]),
                        mute=False
                    ))
                clip.add_new_notes(tuple(new_notes))

    def disconnect(self):
        self._server.shutdown()
        self._server.server_close()
        super(SuperDAW, self).disconnect()
