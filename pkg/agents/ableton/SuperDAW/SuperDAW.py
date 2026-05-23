import Live
from _Framework.ControlSurface import ControlSurface
from .osc_server import OSCServer
from .song import SongHandler
from .track import TrackHandler
from .device import DeviceHandler
from .clip import ClipHandler
from .constants import OSC_LISTEN_PORT, OSC_RESPONSE_PORT
import logging

class SuperDAW(ControlSurface):
    def __init__(self, c_instance):
        super(SuperDAW, self).__init__(c_instance)
        self._osc_server = OSCServer(local_addr=('127.0.0.1', OSC_LISTEN_PORT), remote_addr=('127.0.0.1', OSC_RESPONSE_PORT))
        self.handlers = [SongHandler(self), TrackHandler(self), DeviceHandler(self), ClipHandler(self)]

        self._osc_server.add_handler("/superdaw/transport/play", self._handle_play)
        self._osc_server.add_handler("/superdaw/transport/tempo", self._handle_tempo)
        self._osc_server.add_handler("/superdaw/track/volume", self._handle_volume)
        self._osc_server.add_handler("/superdaw/track/pan", self._handle_pan)
        self._osc_server.add_handler("/superdaw/track/create", self._handle_create_track)
        self._osc_server.add_handler("/superdaw/clip/write", self._handle_clip_write)

    def _handle_play(self, params):
        if params[0]: self.song().start_playing()
        else: self.song().stop_playing()

    def _handle_tempo(self, params):
        self.song().tempo = float(params[0])

    def _handle_volume(self, params):
        idx = int(params[0]); vol = float(params[1])
        if idx < len(self.song().tracks): self.song().tracks[idx].mixer_device.volume.value = vol

    def _handle_pan(self, params):
        idx = int(params[0]); pan = float(params[1])
        if idx < len(self.song().tracks): self.song().tracks[idx].mixer_device.panning.value = pan

    def _handle_create_track(self, params):
        name = params[0]; t = params[1]
        if t == "audio": self.song().create_audio_track()
        else: self.song().create_midi_track()
        self.song().tracks[-1].name = name

    def _handle_clip_write(self, params):
        t_idx = int(params[0]); c_idx = int(params[1]); notes = params[2:]
        if t_idx < len(self.song().tracks):
            slot = self.song().tracks[t_idx].clip_slots[c_idx]
            if not slot.has_clip: slot.create_clip(4.0)
            clip = slot.clip; clip.remove_notes(0, 0, 127, 127)
            new_notes = []
            for i in range(0, len(notes), 4):
                new_notes.append(Live.Clip.MidiNoteSpecification(pitch=int(notes[i]), velocity=int(notes[i+1]), start_time=float(notes[i+2]), duration=float(notes[i+3]), mute=False))
            clip.add_new_notes(tuple(new_notes))

    def update_display(self):
        super(SuperDAW, self).update_display()
        self._osc_server.process()

    def disconnect(self):
        self._osc_server.shutdown()
        super(SuperDAW, self).disconnect()
