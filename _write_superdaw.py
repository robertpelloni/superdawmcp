#!/usr/bin/env python3
"""Write cleaned SuperDAW.py to Ableton remote scripts directory."""

import os
import shutil

content = r'''# -*- coding: utf-8 -*-
import struct
import threading
import json
import socket
import sys
import traceback

try:
    import Queue as queue
except ImportError:
    import queue

from _Framework.ControlSurface import ControlSurface

# ========================================================================== #
#  Minimal OSC decoder (pure Python, no dependencies)
# ========================================================================== #

def _pad4(data):
    """Return data padded to 4-byte boundary with null bytes."""
    extra = (4 - len(data) % 4) % 4
    return data + b'\x00' * extra


def _decode_osc_packet(data):
    """Decode an OSC bundle or message.  Yields (address, args) tuples."""
    try:
        if data.startswith(b'#bundle'):
            pos = 8
            pos += 8
            while pos < len(data):
                size = struct.unpack('>I', data[pos:pos + 4])[0]
                pos += 4
                for msg in _decode_osc_packet(data[pos:pos + size]):
                    yield msg
                pos += size
        else:
            addr_end = data.find(b'\x00')
            if addr_end < 0:
                return
            address = data[:addr_end].decode('utf-8', errors='replace')
            pos = addr_end + 1
            pos = ((pos + 3) // 4) * 4

            if pos >= len(data) or data[pos:pos + 1] != b',':
                yield (address, [])
                return
            tag_end = data.find(b'\x00', pos)
            if tag_end < 0:
                yield (address, [])
                return
            tags = data[pos + 1:tag_end]
            pos = tag_end + 1
            pos = ((pos + 3) // 4) * 4

            args = []
            for t in tags:
                if t == ord('i'):
                    if pos + 4 > len(data):
                        args.append(None)
                        continue
                    val = struct.unpack('>i', data[pos:pos + 4])[0]
                    args.append(val)
                    pos += 4
                elif t == ord('f'):
                    if pos + 4 > len(data):
                        args.append(None)
                        continue
                    val = struct.unpack('>f', data[pos:pos + 4])[0]
                    args.append(val)
                    pos += 4
                elif t == ord('s'):
                    s_end = data.find(b'\x00', pos)
                    if s_end < 0:
                        args.append(None)
                        continue
                    val = data[pos:s_end].decode('utf-8', errors='replace')
                    pos = ((s_end + 4) // 4) * 4
                    args.append(val)
                elif t == ord('b'):
                    if pos + 4 > len(data):
                        args.append(None)
                        continue
                    size = struct.unpack('>I', data[pos:pos + 4])[0]
                    pos += 4
                    val = data[pos:pos + size]
                    pos = ((pos + size + 3) // 4) * 4
                    args.append(val)
                elif t == ord('T'):
                    args.append(True)
                elif t == ord('F'):
                    args.append(False)
                else:
                    args.append(None)
            yield (address, args)
    except Exception:
        return


# ========================================================================== #
#  SuperDAW Agent
# ========================================================================== #

OSC_LISTEN_PORT = 11000
OSC_RESPONSE_PORT = 11001


def _track_id(raw):
    """Parse track ID (string or int) to valid track index."""
    try:
        if raw is None:
            return 0
        return int(raw)
    except (ValueError, TypeError):
        return 0


class SuperDAW(ControlSurface):
    """SuperDAW-MCP Ableton Live 10 Agent via raw OSC/UDP."""

    def __init__(self, c_instance):
        super(SuperDAW, self).__init__(c_instance)
        self._log = self.log_message
        self._task_queue = queue.Queue()
        self._running = True

        self._client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

        self._server = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        self._server.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
        self._server.bind(('127.0.0.1', OSC_LISTEN_PORT))
        self._server.settimeout(0.1)

        self._server_thread = threading.Thread(target=self._server_loop)
        self._server_thread.daemon = True
        self._server_thread.start()

        self._setup_listeners()
        self.show_message("SuperDAW Agent ready (port %d)" % OSC_LISTEN_PORT)

    def _send_osc(self, address, value):
        """Build and send a minimal OSC message."""
        try:
            addr_bytes = _pad4(address.encode('utf-8'))
            type_tag = b','
            data = b''
            if isinstance(value, bool):
                type_tag += b'T' if value else b'F'
            elif isinstance(value, (int, long)):
                type_tag += b'i'
                data = struct.pack('>i', value)
            elif isinstance(value, float):
                type_tag += b'f'
                data = struct.pack('>f', value)
            elif isinstance(value, basestring):
                type_tag += b's'
                data = _pad4(value.encode('utf-8'))
            elif isinstance(value, bytes):
                type_tag += b'b'
                data = struct.pack('>I', len(value)) + _pad4(value)
            else:
                return
            type_tag_bytes = _pad4(type_tag)
            msg = addr_bytes + type_tag_bytes + data
            self._client.sendto(msg, ('127.0.0.1', OSC_RESPONSE_PORT))
        except Exception:
            pass

    def _server_loop(self):
        while self._running:
            try:
                data, addr = self._server.recvfrom(65536)
                for osc_addr, osc_args in _decode_osc_packet(data):
                    self._task_queue.put((osc_addr, osc_args))
            except socket.timeout:
                pass
            except Exception:
                self._log("SuperDAW server error: " + traceback.format_exc())

    def _setup_listeners(self):
        self.song().add_is_playing_listener(self._on_playing_changed)
        self.song().add_tempo_listener(self._on_tempo_changed)

    def _on_playing_changed(self):
        self._send_osc("/superdaw/state/playing", self.song().is_playing)

    def _on_tempo_changed(self):
        self._send_osc("/superdaw/state/tempo", float(self.song().tempo))

    def update_display(self):
        super(SuperDAW, self).update_display()
        while not self._task_queue.empty():
            try:
                addr, args = self._task_queue.get_nowait()
                self._handle(addr, args)
            except queue.Empty:
                break
            except Exception as e:
                self._log("SuperDAW error: %s" % str(e))

    def _handle(self, addr, args):
        self._log("OSC handle: %s args=%s" % (addr, args))
        if not isinstance(args, (list, tuple)):
            args = []
        args = [a for a in args if a is not None]
        new_args = []
        for arg in args:
            if isinstance(arg, bytes):
                new_args.append(arg.decode('utf-8', errors='replace'))
            else:
                new_args.append(arg)
        args = new_args

        if addr == "/superdaw/transport/play":
            if args and args[0]:
                self.song().start_playing()
            else:
                self.song().stop_playing()
        elif addr == "/superdaw/transport/tempo":
            if args and args[0] is not None:
                try:
                    self.song().tempo = float(args[0])
                except Exception as e:
                    self._log("SuperDAW tempo error: %s" % str(e))
        elif addr == "/superdaw/track/create":
            self._cmd_track_create(args)
        elif addr == "/superdaw/track/volume":
            self._cmd_track_volume(args)
        elif addr == "/superdaw/track/pan":
            self._cmd_track_pan(args)
        elif addr == "/superdaw/clip/write":
            self._cmd_clip_write(args)
        elif addr == "/superdaw/clip/delete":
            self._cmd_clip_delete(args)
        elif addr == "/superdaw/plugin/param":
            self._cmd_plugin_param(args)
        elif addr == "/superdaw/midi/cc":
            self._cmd_midi_cc(args)
        elif addr == "/live/scene/fire":
            self._cmd_scene_fire(args)

    def _cmd_track_create(self, args):
        name = str(args[0]) if len(args) > 0 else "Track"
        tracks = list(self.song().tracks)
        if tracks:
            try:
                self.song().duplicate_track(tracks[-1])
                new = list(self.song().tracks)
                if len(new) > len(tracks):
                    new[-1].name = name
                    track_index = len(new) - 1
                    self._send_osc("/superdaw/state/track_created", track_index)
            except Exception:
                pass

    def _cmd_track_volume(self, args):
        if len(args) < 2:
            return
        if args[0] is None or args[1] is None:
            return
        idx = _track_id(args[0])
        try:
            vol = float(args[1])
        except (ValueError, TypeError):
            return
        if 0 <= idx < len(self.song().tracks):
            self.song().tracks[idx].mixer_device.volume.value = vol

    def _cmd_track_pan(self, args):
        if len(args) < 2:
            return
        if args[0] is None or args[1] is None:
            return
        idx = _track_id(args[0])
        try:
            pan = float(args[1])
        except (ValueError, TypeError):
            return
        if 0 <= idx < len(self.song().tracks):
            self.song().tracks[idx].mixer_device.panning.value = pan

    def _cmd_clip_write(self, args):
        if len(args) < 3:
            return
        track_idx = _track_id(args[0])
        try:
            clip_idx = int(args[1])
        except (ValueError, TypeError):
            return
        tracks = self.song().tracks
        if track_idx >= len(tracks):
            return
        try:
            slot = tracks[track_idx].clip_slots[clip_idx]
        except IndexError:
            return
        try:
            if not slot.has_clip:
                slot.create_clip(4.0)
            clip = slot.clip
            raw = args[2]
            if isinstance(raw, bytes):
                raw = raw.decode('utf-8', errors='replace')
            notes_data = json.loads(raw)
            notes = []
            for n in notes_data:
                pitch = int(n.get('pitch', 60))
                start = float(n.get('start_beat', 0.0))
                dur = float(n.get('duration', 0.25))
                vel = int(n.get('velocity', 100))
                notes.append((pitch, start, dur, vel))
            clip.add_new_notes(tuple(notes))
        except Exception:
            self._log("SuperDAW clip_write error: " + traceback.format_exc())

    def _cmd_clip_delete(self, args):
        if len(args) < 2:
            return
        track_idx = _track_id(args[0])
        try:
            clip_idx = int(args[1])
        except (ValueError, TypeError):
            return
        tracks = self.song().tracks
        if track_idx >= len(tracks):
            return
        try:
            slot = tracks[track_idx].clip_slots[clip_idx]
        except IndexError:
            return
        if slot.has_clip:
            slot.delete_clip()

    def _cmd_plugin_param(self, args):
        if len(args) < 4:
            return
        if args[0] is None or args[2] is None or args[3] is None:
            return
        track_idx = _track_id(args[0])
        try:
            param_index = int(args[2])
            value = float(args[3])
        except (ValueError, TypeError):
            return
        tracks = self.song().tracks
        if track_idx >= len(tracks):
            return
        all_params = []
        for dev in tracks[track_idx].devices:
            if hasattr(dev, 'parameters'):
                all_params.extend(list(dev.parameters))
        if 0 <= param_index < len(all_params):
            all_params[param_index].value = value

    def _cmd_midi_cc(self, args):
        if len(args) < 3:
            return
        if args[0] is None or args[1] is None or args[2] is None:
            return
        track_idx = _track_id(args[0])
        try:
            controller = int(args[1])
            value = int(args[2])
        except (ValueError, TypeError):
            return
        if 0 <= track_idx < len(self.song().tracks):
            try:
                from _Framework import MidiMap
                MidiMap.send_midi(self.song().tracks[track_idx], (0xB0, controller, value))
            except Exception:
                pass

    def _cmd_scene_fire(self, args):
        if not args:
            return
        if args[0] is None:
            return
        try:
            idx = int(args[0])
        except (ValueError, TypeError):
            return
        if 0 <= idx < len(self.song().scenes):
            self.song().scenes[idx].fire()

    def disconnect(self):
        self._running = False
        self.song().remove_is_playing_listener(self._on_playing_changed)
        self.song().remove_tempo_listener(self._on_tempo_changed)
        try:
            self._server.close()
        except Exception:
            pass
        super(SuperDAW, self).disconnect()
'''

target = r"C:\ProgramData\Ableton\Live 10 Standard\Resources\MIDI Remote Scripts\SuperDAW\SuperDAW.py"

# Delete __pycache__ first
pycache = os.path.join(os.path.dirname(target), "__pycache__")
if os.path.isdir(pycache):
    shutil.rmtree(pycache)

with open(target, "wb") as f:
    f.write(content.encode("utf-8"))

print("Written to " + target)
print("Also delete stale __pycache__")
