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
        self.show_message("SuperDAW Initializing...")

        self.osc_server = OSCServer(
            local_addr=('127.0.0.1', OSC_LISTEN_PORT),
            remote_addr=('127.0.0.1', OSC_RESPONSE_PORT)
        )

        self.handlers = [
            SongHandler(self),
            TrackHandler(self),
            DeviceHandler(self),
            ClipHandler(self)
        ]

        # Register unified protocol handlers
        self.osc_server.add_handler("/superdaw/transport/play", self._handle_play)
        self.osc_server.add_handler("/superdaw/transport/tempo", self._handle_tempo)
        self.osc_server.add_handler("/superdaw/track/volume", self._handle_volume)
        self.osc_server.add_handler("/superdaw/track/pan", self._handle_pan)
        self.osc_server.add_handler("/superdaw/track/create", self._handle_create_track)
        self.osc_server.add_handler("/superdaw/clip/write", self._handle_clip_write)

    def _handle_play(self, params):
        if params[0]: self.song().start_playing()
        else: self.song().stop_playing()

    def _handle_tempo(self, params):
        self.song().tempo = params[0]

    def _handle_volume(self, params):
        # Implementation logic here
        pass

    def _handle_pan(self, params):
        pass

    def _handle_create_track(self, params):
        pass

    def _handle_clip_write(self, params):
        pass

    def update_display(self):
        super(SuperDAW, self).update_display()
        self.osc_server.process()

    def disconnect(self):
        self.show_message("SuperDAW Disconnecting...")
        for handler in self.handlers:
            handler.clear_api()
        self.osc_server.shutdown()
        super(SuperDAW, self).disconnect()
