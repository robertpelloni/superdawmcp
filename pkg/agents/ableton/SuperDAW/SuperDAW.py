import Live
import logging
from _Framework.ControlSurface import ControlSurface
from .pythonosc.osc_server import BlockingOSCUDPServer
from .pythonosc.dispatcher import Dispatcher
from .pythonosc.udp_client import SimpleUDPClient
import threading
import json
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

        self._active_listeners = []

        # Track listeners for volume/pan updates
        self._setup_listeners()

    def _setup_listeners(self):
        # Initial state send
        self._send_full_state()

        # Listen to transport
        self.song().add_is_playing_listener(self._on_playing_changed)
        self.song().add_tempo_listener(self._on_tempo_changed)
        self.song().add_tracks_listener(self._on_arrangement_changed)

        # Monitor selected device for parameter changes
        self.song().view.add_selected_track_listener(self._on_selected_track_changed)
        self._on_selected_track_changed()

    def _on_selected_track_changed(self):
        self.song().view.selected_track.view.add_selected_device_listener(self._on_selected_device_changed)
        self._on_selected_device_changed()

    def _on_selected_device_changed(self):
        # Clear existing parameter listeners
        for param in self._active_listeners:
            if param.value_has_listener(self._on_param_value_changed):
                param.remove_value_listener(self._on_param_value_changed)
        self._active_listeners = []

        device = self.song().view.selected_track.view.selected_device
        if device:
            for param in device.parameters:
                param.add_value_listener(self._on_param_value_changed)
                self._active_listeners.append(param)

    def _on_param_value_changed(self):
        # Broadcast all parameter values for the selected device
        # Address schema: /superdaw/state/plugin/param {track_idx} {device_name} {param_name} {value}
        track = self.song().view.selected_track
        track_idx = list(self.song().tracks).index(track)
        device = track.view.selected_device
        if device:
            for param in device.parameters:
                # We only send the one that actually changed if we had a reference,
                # but for simplicity we can send all or try to find the match.
                # Here we just broadcast the changed value if we can identify it.
                pass
            # More efficient: find which param changed. But Ableton doesn't pass the param to the callback.
            # So we broadcast current state of selected device.
            params_state = []
            for p in device.parameters:
                params_state.append({"n": p.name, "v": p.value})
            self._client.send_message("/superdaw/state/plugin/params", [track_idx, device.name, json.dumps(params_state)])

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
        elif addr == "/superdaw/plugin/param":
            track_idx = int(args[0]); device_name = str(args[1]); param_idx = int(args[2]); val = float(args[3])
            if track_idx < len(self.song().tracks):
                track = self.song().tracks[track_idx]
                for device in track.devices:
                    if device.name == device_name:
                        if param_idx < len(device.parameters):
                            device.parameters[param_idx].value = val
                        break
        elif addr == "/superdaw/midi/cc":
            track_idx = int(args[0]); ctrl = int(args[1]); val = int(args[2])
            # Ableton MIDI Remote Script doesn't have a direct "send_cc_to_track"
            # but we can use the track's MIDI input if it was possible or control parameters.
            # Usually SendCC is used for external hardware or automation.
            pass
        elif addr == "/superdaw/clip/write":
            track_idx = int(args[0])
            clip_idx = int(args[1])
            if track_idx < len(self.song().tracks):
                track = self.song().tracks[track_idx]
                if clip_idx < len(track.clip_slots):
                    slot = track.clip_slots[clip_idx]
                    if not slot.has_clip:
                        slot.create_clip(4.0) # 1 bar default
                    clip = slot.clip
                    try:
                        notes_data = json.loads(args[2])
                        notes = []
                        for n in notes_data:
                            notes.append(Live.Clip.MidiNoteSpecification(
                                pitch=int(n['pitch']),
                                start_time=float(n['start_beat']),
                                duration=float(n['duration']),
                                velocity=int(n['velocity']),
                                mute=False
                            ))
                        clip.add_new_notes(tuple(notes))
                    except:
                        pass

    def disconnect(self):
        self.song().remove_is_playing_listener(self._on_playing_changed)
        self.song().remove_tempo_listener(self._on_tempo_changed)
        self._server.shutdown()
        self._server.server_close()
        super(SuperDAW, self).disconnect()
