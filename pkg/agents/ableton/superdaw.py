import Live
from _Framework.ControlSurface import ControlSurface
from .superdawpkg.osc_server import OSCServer
import logging
import json

class SuperDAW(ControlSurface):
    def __init__(self, c_instance):
        super(SuperDAW, self).__init__(c_instance)
        self._logger = logging.getLogger("superdaw")
        self._osc_server = OSCServer(local_addr=('127.0.0.1', 11000), remote_addr=('127.0.0.1', 11001))
        self.init_handlers()
        self._logger.info("SuperDAW Ableton Adapter Initialized")

    def init_handlers(self):
        # Transport
        self._osc_server.add_handler("/superdaw/transport/play", self._handle_play)
        self._osc_server.add_handler("/superdaw/transport/tempo", self._handle_tempo)

        # Mixer
        self._osc_server.add_handler("/superdaw/track/volume", self._handle_volume)
        self._osc_server.add_handler("/superdaw/track/create", self._handle_create_track)

        # Clips
        self._osc_server.add_handler("/superdaw/clip/write", self._handle_clip_write)

    def _handle_play(self, params):
        play = int(params[0])
        if play:
            self.song().start_playing()
        else:
            self.song().stop_playing()

    def _handle_tempo(self, params):
        self.song().tempo = float(params[0])

    def _handle_volume(self, params):
        track_idx = int(params[0])
        volume = float(params[1])
        if track_idx < len(self.song().tracks):
            self.song().tracks[track_idx].mixer_device.volume.value = volume

    def _handle_create_track(self, params):
        name = params[0]
        track_type = params[1]
        if track_type == "audio":
            self.song().create_audio_track()
        else:
            self.song().create_midi_track()
        self.song().tracks[-1].name = name

    def _handle_clip_write(self, params):
        track_idx = int(params[0])
        clip_idx = int(params[1])
        # Flattened notes are in params[2:]
        notes_data = params[2:]
        if track_idx < len(self.song().tracks):
            track = self.song().tracks[track_idx]
            if clip_idx < len(track.clip_slots):
                clip_slot = track.clip_slots[clip_idx]
                if not clip_slot.has_clip:
                    clip_slot.create_clip(4.0)
                clip = clip_slot.clip
                clip.remove_notes(0, 0, 127, 127)

                new_notes = []
                for i in range(0, len(notes_data), 4):
                    pitch = int(notes_data[i])
                    velocity = int(notes_data[i+1])
                    start = float(notes_data[i+2])
                    duration = float(notes_data[i+3])
                    new_notes.append(Live.Clip.MidiNoteSpecification(
                        pitch=pitch,
                        start_time=start,
                        duration=duration,
                        velocity=velocity,
                        mute=False
                    ))
                clip.add_new_notes(tuple(new_notes))

    def update_display(self):
        super(SuperDAW, self).update_display()
        self._osc_server.process()

    def disconnect(self):
        self._osc_server.shutdown()
        super(SuperDAW, self).disconnect()
