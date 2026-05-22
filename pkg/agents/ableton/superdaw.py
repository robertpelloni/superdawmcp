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

    def update_display(self):
        super(SuperDAW, self).update_display()
        self._osc_server.process()

    def disconnect(self):
        self._osc_server.shutdown()
        super(SuperDAW, self).disconnect()
