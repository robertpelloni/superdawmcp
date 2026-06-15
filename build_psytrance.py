#!/usr/bin/env python3
"""
Full-On Psytrance Production Builder for Ableton Live 10 Standard
=================================================================

128 bars @ 148 BPM (~8 minutes)

Song Structure:
  1.  Ambient Intro       (bars 1-8)    - pad, riser, atmosphere
  2.  Groove Intro         (bars 9-24)   - kick+bass enter, hats build
  3.  First Breakdown      (bars 25-32)  - filter close, drums drop
  4.  First Drop           (bars 33-56)  - full power, lead arp, stabs
  5.  Second Breakdown     (bars 57-64)  - riser, filter sweep
  6.  Second Drop          (bars 65-88)  - full power, variation
  7.  Bridge               (bars 89-104) - pad solo, glitchy hats
  8.  Final Breakdown      (bars 105-112)- tension build
  9.  Final Drop           (bars 113-120)- maximum energy
  10. Outro                (bars 121-128)- decay, fade, end

Mixing: Volume hierarchy, stereo field, frequency carving
Effects: Reverb, Delay, Auto Filter, Compressor, Saturator, EQ Three
Mastering: Compressor on master concept, saturation, limiting via volume
Transitions: Filter sweeps, reverb tails, delay throws, riser FX, volume automation
"""

import socket
import struct
import time
import json
import os
import tempfile
import threading
from pathlib import Path
import sys

sys.path.append(str(Path(__file__).parent / "pkg" / "client" / "py"))


# ============================================================
# OSC TRANSPORT
# ============================================================


def pad4(data):
    data = data + b"\x00"
    extra = (4 - len(data) % 4) % 4
    return data + b"\x00" * extra


def send_osc(address, *args):
    addr_bytes = pad4(address.encode("utf-8"))
    type_tag = b","
    data = b""
    for arg in args:
        if isinstance(arg, int):
            type_tag += b"i"
            data += struct.pack(">i", arg)
        elif isinstance(arg, float):
            type_tag += b"f"
            data += struct.pack(">f", arg)
        elif isinstance(arg, str):
            type_tag += b"s"
            data += pad4(arg.encode("utf-8"))
    type_tag_bytes = pad4(type_tag)
    msg = addr_bytes + type_tag_bytes + data
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.sendto(msg, ("127.0.0.1", 11000))
    sock.close()


def write_clip(track_idx, clip_slot, notes):
    """Write clip data to temp file and send path via OSC to avoid UDP overflow."""
    clip_data = json.dumps(notes)
    tmp_path = os.path.join(tempfile.gettempdir(), "superdaw_clip_%d_%d.json" % (track_idx, clip_slot))
    with open(tmp_path, "w") as f:
        f.write(clip_data)
    send_osc("/superdaw/clip/fromfile", track_idx, clip_slot, tmp_path)
    time.sleep(1.0)


def read_temp(filename, timeout=3.0):
    path = os.path.join(tempfile.gettempdir(), filename)
    try:
        old_mtime = os.path.getmtime(path)
    except Exception:
        old_mtime = 0
    deadline = time.time() + timeout
    while time.time() < deadline:
        try:
            mtime = os.path.getmtime(path)
            if mtime > old_mtime:
                with open(path, "r") as f:
                    return f.read().strip()
        except Exception:
            pass
        time.sleep(0.2)
    try:
        with open(path, "r") as f:
            return f.read().strip()
    except Exception:
        return None


def log(msg):
    print("[>] " + msg)


# ============================================================
# SONG PARAMETERS
# ============================================================
BPM = 148
BARS = 128

# Section boundaries (in bars)
SECTIONS = {
    "ambient_intro": (0, 8),
    "groove_intro": (8, 24),
    "breakdown_1": (24, 32),
    "drop_1": (32, 56),
    "breakdown_2": (56, 64),
    "drop_2": (64, 88),
    "bridge": (88, 104),
    "breakdown_3": (104, 112),
    "drop_3": (112, 120),
    "outro": (120, 128),
}

# ============================================================
# TRACK LAYOUT
# ============================================================
# (track_name, instrument_preset, volume, pan, mute_intro_bars)
TRACK_LAYOUT = [
    # -- Rhythm Section (front of mix) --
    ("KICK", "808 Core Kit", 0.92, 0.0, 8),  # 4
    ("HAT CL", "808 Core Kit", 0.62, 0.35, 8),  # 5
    ("HAT OP", "808 Core Kit", 0.55, -0.25, 16),  # 6
    ("SNARE", "808 Core Kit", 0.78, -0.15, 8),  # 7
    ("PERC", "Percussion Core Kit", 0.58, 0.45, 24),  # 8
    # -- Bass Section (center, powerful) --
    ("BASS", "Analog Saw Bass", 0.88, 0.0, 8),  # 9
    # -- Synth Section (stereo field) --
    ("LEAD", "Miami Pluck", 0.68, 0.30, 32),  # 10
    ("ARP", "Zapp Bass", 0.55, -0.40, 32),  # 11
    ("STAB", "Drubb Stabs", 0.65, 0.0, 32),  # 12
    # -- Atmospheric Section (wide) --
    ("PAD", "All Alone Pad", 0.48, -0.30, 0),  # 13
    ("DARK PAD", "Deep In Dark", 0.35, 0.50, 0),  # 14
    ("RISER", "ENCOM Riser", 0.55, 0.0, 999),  # 15 - triggered manually
    ("SWEEP", "Chord Sweeper", 0.40, 0.0, 999),  # 16 - triggered manually
    ("FX", "Ominous Thud", 0.70, 0.0, 999),  # 17 - triggered manually
]

# Effect chains per track
# (track_name_key, [effect_search_names])
EFFECT_CHAINS = {
    "KICK": ["Compressor", "EQ Three"],
    "BASS": ["Auto Filter", "Saturator"],
    "LEAD": ["Reverb", "Delay"],
    "ARP": ["Reverb", "Auto Filter"],
    "PAD": ["Reverb"],
    "DARK": ["Reverb", "Delay"],
    "STAB": ["Reverb"],
    "HATCL": ["Reverb"],
    "HATOP": ["Delay"],
    "SNARE": ["Compressor", "Reverb"],
    "PERC": ["Reverb"],
    "RISER": ["Reverb"],
    "SWEEP": ["Reverb", "Delay"],
    "FX": ["Reverb"],
}


# ============================================================
# PARAMETER VALUE CONVERTERS
# ============================================================
def filter_hz(norm):
    """0-1 normalized -> Auto Filter Frequency (20-135)."""
    return 20.0 + norm * 115.0


def rack_wet(norm):
    """0-1 normalized -> Rack Dry/Wet (0-127)."""
    return norm * 127.0


def band_wet(norm):
    """0-1 normalized -> Reverb band Dry/Wet (0-127)."""
    return norm * 127.0


# ============================================================
# MIDI CLIP BUILDERS
# ============================================================


def build_kick():
    """4-on-the-floor with section dynamics."""
    notes = []
    for bar in range(BARS):
        b = bar * 4
        if bar < 8:
            # Ambient intro: no kick
            pass
        elif bar < 24:
            # Groove intro: solid 4/4
            for beat in range(4):
                notes.append(
                    {
                        "pitch": 36,
                        "start_beat": b + beat,
                        "duration": 0.25,
                        "velocity": 118,
                    }
                )
        elif bar < 32:
            # Breakdown 1: kick drops to every 2 bars
            if bar % 2 == 0:
                notes.append(
                    {"pitch": 36, "start_beat": b, "duration": 0.25, "velocity": 95}
                )
        elif bar < 56:
            # Drop 1: full power
            for beat in range(4):
                notes.append(
                    {
                        "pitch": 36,
                        "start_beat": b + beat,
                        "duration": 0.25,
                        "velocity": 127,
                    }
                )
        elif bar < 64:
            # Breakdown 2: sparse
            if bar % 4 == 0:
                notes.append(
                    {"pitch": 36, "start_beat": b, "duration": 0.25, "velocity": 90}
                )
        elif bar < 88:
            # Drop 2: full power with ghost kicks
            for beat in range(4):
                notes.append(
                    {
                        "pitch": 36,
                        "start_beat": b + beat,
                        "duration": 0.25,
                        "velocity": 127,
                    }
                )
            # Ghost kick on "&" of 2 for drive
            if bar % 2 == 0:
                notes.append(
                    {
                        "pitch": 36,
                        "start_beat": b + 2.5,
                        "duration": 0.15,
                        "velocity": 60,
                    }
                )
        elif bar < 104:
            # Bridge: kick every 2 beats
            for beat in range(0, 4, 2):
                notes.append(
                    {
                        "pitch": 36,
                        "start_beat": b + beat,
                        "duration": 0.25,
                        "velocity": 100,
                    }
                )
        elif bar < 112:
            # Breakdown 3: kick absent
            if bar == 108:
                notes.append(
                    {"pitch": 36, "start_beat": b, "duration": 0.25, "velocity": 80}
                )
        elif bar < 120:
            # Final drop: maximum energy
            for beat in range(4):
                notes.append(
                    {
                        "pitch": 36,
                        "start_beat": b + beat,
                        "duration": 0.25,
                        "velocity": 127,
                    }
                )
        else:
            # Outro: kick thins and fades
            for beat in range(0, 4, 2):
                vel = max(50, 120 - (bar - 120) * 9)
                notes.append(
                    {
                        "pitch": 36,
                        "start_beat": b + beat,
                        "duration": 0.25,
                        "velocity": vel,
                    }
                )
    return notes


def build_hat_closed():
    """Closed hat - 16th notes with varying patterns per section."""
    notes = []
    for bar in range(BARS):
        b = bar * 4
        if bar < 8:
            pass
        elif bar < 16:
            # Intro: 8th notes
            for beat in range(8):
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.5,
                        "duration": 0.1,
                        "velocity": 65,
                    }
                )
        elif bar < 24:
            # Groove: 16ths
            for beat in range(16):
                vel = 85 if beat % 4 == 0 else 55
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.1,
                        "velocity": vel,
                    }
                )
        elif bar < 32:
            # Breakdown: hats fade
            for beat in range(8):
                vel = max(20, 55 - (bar - 24) * 5)
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.5,
                        "duration": 0.1,
                        "velocity": vel,
                    }
                )
        elif bar < 56:
            # Drop 1: driving 16ths
            for beat in range(16):
                vel = 90 if beat % 2 == 0 else 62
                if beat % 4 == 0:
                    vel = 95
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.1,
                        "velocity": vel,
                    }
                )
        elif bar < 64:
            # Breakdown 2: sparse
            for beat in range(4):
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat,
                        "duration": 0.1,
                        "velocity": 35,
                    }
                )
        elif bar < 88:
            # Drop 2: same driving but with accent pattern
            for beat in range(16):
                if beat % 8 == 0:
                    vel = 100
                elif beat % 4 == 0:
                    vel = 85
                elif beat % 2 == 0:
                    vel = 70
                else:
                    vel = 55
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.1,
                        "velocity": vel,
                    }
                )
        elif bar < 104:
            # Bridge: half-time feel, 8ths
            for beat in range(8):
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.5,
                        "duration": 0.15,
                        "velocity": 50,
                    }
                )
        elif bar < 112:
            # Breakdown 3: building hats
            for beat in range(16):
                vel = min(100, 30 + (bar - 104) * 8)
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.1,
                        "velocity": vel,
                    }
                )
        elif bar < 120:
            # Final drop: maximum
            for beat in range(16):
                vel = 95 if beat % 2 == 0 else 65
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.1,
                        "velocity": vel,
                    }
                )
        else:
            # Outro: thin out
            for beat in range(8):
                vel = max(20, 60 - (bar - 120) * 5)
                notes.append(
                    {
                        "pitch": 42,
                        "start_beat": b + beat * 0.5,
                        "duration": 0.1,
                        "velocity": vel,
                    }
                )
    return notes


def build_hat_open():
    """Open hat on offbeats, accent patterns."""
    notes = []
    for bar in range(BARS):
        b = bar * 4
        if bar < 16:
            pass
        elif bar < 24:
            # Groove: offbeat open hats
            for beat in range(4):
                notes.append(
                    {
                        "pitch": 46,
                        "start_beat": b + beat + 0.5,
                        "duration": 0.2,
                        "velocity": 55,
                    }
                )
        elif bar < 32:
            # Breakdown: less
            if bar % 2 == 0:
                notes.append(
                    {
                        "pitch": 46,
                        "start_beat": b + 2.5,
                        "duration": 0.25,
                        "velocity": 40,
                    }
                )
        elif bar < 56:
            # Drop 1: every other offbeat
            for beat in range(4):
                if beat % 2 == 0:
                    notes.append(
                        {
                            "pitch": 46,
                            "start_beat": b + beat + 0.5,
                            "duration": 0.15,
                            "velocity": 60,
                        }
                    )
        elif bar < 64:
            pass
        elif bar < 88:
            # Drop 2: double-time offbeats
            for beat in range(8):
                notes.append(
                    {
                        "pitch": 46,
                        "start_beat": b + beat * 0.5 + 0.25,
                        "duration": 0.12,
                        "velocity": 50,
                    }
                )
        elif bar < 104:
            # Bridge: sparse
            notes.append(
                {"pitch": 46, "start_beat": b + 1, "duration": 0.3, "velocity": 40}
            )
            notes.append(
                {"pitch": 46, "start_beat": b + 3, "duration": 0.3, "velocity": 40}
            )
        elif bar < 112:
            # Building
            if bar >= 108:
                notes.append(
                    {"pitch": 46, "start_beat": b + 1, "duration": 0.2, "velocity": 55}
                )
                notes.append(
                    {"pitch": 46, "start_beat": b + 3, "duration": 0.2, "velocity": 55}
                )
        elif bar < 120:
            # Final drop
            for beat in range(4):
                if beat % 2 == 0:
                    notes.append(
                        {
                            "pitch": 46,
                            "start_beat": b + beat + 0.5,
                            "duration": 0.15,
                            "velocity": 65,
                        }
                    )
        else:
            # Outro
            if bar % 2 == 0:
                notes.append(
                    {"pitch": 46, "start_beat": b + 1, "duration": 0.25, "velocity": 35}
                )
    return notes


def build_snare():
    """Backbeat with ghost notes and fills."""
    notes = []
    for bar in range(BARS):
        b = bar * 4
        if bar < 8:
            pass
        elif bar < 24:
            # Groove: 2 and 4
            notes.append(
                {"pitch": 40, "start_beat": b + 1, "duration": 0.2, "velocity": 95}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 3, "duration": 0.2, "velocity": 95}
            )
        elif bar < 32:
            # Breakdown: snare roll building at end
            if bar >= 30:
                for i in range(16):
                    vel = 60 + i * 4
                    notes.append(
                        {
                            "pitch": 40,
                            "start_beat": b + i * 0.25,
                            "duration": 0.1,
                            "velocity": min(vel, 127),
                        }
                    )
            else:
                notes.append(
                    {"pitch": 40, "start_beat": b + 1, "duration": 0.2, "velocity": 80}
                )
        elif bar < 56:
            # Drop 1: backbeat + ghost notes
            notes.append(
                {"pitch": 40, "start_beat": b + 1, "duration": 0.2, "velocity": 110}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 3, "duration": 0.2, "velocity": 110}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 1.5, "duration": 0.1, "velocity": 48}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 3.5, "duration": 0.1, "velocity": 48}
            )
        elif bar < 64:
            # Breakdown 2: snare roll
            if bar >= 62:
                for i in range(16):
                    vel = 50 + i * 5
                    notes.append(
                        {
                            "pitch": 40,
                            "start_beat": b + i * 0.25,
                            "duration": 0.1,
                            "velocity": min(vel, 127),
                        }
                    )
        elif bar < 88:
            # Drop 2: similar + rim shot variation
            notes.append(
                {"pitch": 40, "start_beat": b + 1, "duration": 0.2, "velocity": 112}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 3, "duration": 0.2, "velocity": 112}
            )
            # Ghost on "e" and "ah"
            notes.append(
                {"pitch": 40, "start_beat": b + 0.75, "duration": 0.08, "velocity": 42}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 2.75, "duration": 0.08, "velocity": 42}
            )
            # Extra rim on bar endings
            if bar % 4 == 3:
                notes.append(
                    {
                        "pitch": 37,
                        "start_beat": b + 3.75,
                        "duration": 0.1,
                        "velocity": 70,
                    }
                )
        elif bar < 104:
            # Bridge: sparse
            notes.append(
                {"pitch": 40, "start_beat": b + 3, "duration": 0.2, "velocity": 70}
            )
        elif bar < 112:
            # Breakdown 3: building snare roll
            for i in range(min(8 + (bar - 104) * 2, 32)):
                vel = 40 + i * 3
                notes.append(
                    {
                        "pitch": 40,
                        "start_beat": b + i * 0.125,
                        "duration": 0.08,
                        "velocity": min(vel, 127),
                    }
                )
        elif bar < 120:
            # Final drop: aggressive
            notes.append(
                {"pitch": 40, "start_beat": b + 1, "duration": 0.2, "velocity": 115}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 3, "duration": 0.2, "velocity": 115}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 1.5, "duration": 0.1, "velocity": 55}
            )
            notes.append(
                {"pitch": 40, "start_beat": b + 3.5, "duration": 0.1, "velocity": 55}
            )
        else:
            # Outro: simple backbeat fading
            vel = max(40, 100 - (bar - 120) * 8)
            notes.append(
                {"pitch": 40, "start_beat": b + 1, "duration": 0.2, "velocity": vel}
            )
    return notes


def build_perc():
    """Percussion - congas, toms, shakers for groove."""
    notes = []
    for bar in range(BARS):
        b = bar * 4
        if bar < 24:
            pass
        elif bar < 32:
            # Light percussion in breakdown
            notes.append(
                {"pitch": 48, "start_beat": b + 2, "duration": 0.15, "velocity": 55}
            )
        elif bar < 56:
            # Drop 1: conga pattern
            notes.append(
                {"pitch": 48, "start_beat": b + 0.5, "duration": 0.15, "velocity": 65}
            )
            notes.append(
                {"pitch": 50, "start_beat": b + 2.5, "duration": 0.15, "velocity": 60}
            )
            if bar % 2 == 0:
                notes.append(
                    {
                        "pitch": 45,
                        "start_beat": b + 1.5,
                        "duration": 0.15,
                        "velocity": 55,
                    }
                )
        elif bar < 64:
            pass
        elif bar < 88:
            # Drop 2: busier percussion
            notes.append(
                {"pitch": 48, "start_beat": b + 0.5, "duration": 0.12, "velocity": 70}
            )
            notes.append(
                {"pitch": 50, "start_beat": b + 1.5, "duration": 0.12, "velocity": 60}
            )
            notes.append(
                {"pitch": 48, "start_beat": b + 2.5, "duration": 0.12, "velocity": 65}
            )
            notes.append(
                {"pitch": 50, "start_beat": b + 3.5, "duration": 0.12, "velocity": 55}
            )
            # Shaker on 16ths (pitch 70)
            for beat in range(16):
                if beat % 4 != 0:
                    notes.append(
                        {
                            "pitch": 70,
                            "start_beat": b + beat * 0.25,
                            "duration": 0.08,
                            "velocity": 35,
                        }
                    )
        elif bar < 104:
            # Bridge: light
            notes.append(
                {"pitch": 48, "start_beat": b + 2, "duration": 0.2, "velocity": 45}
            )
        elif bar < 112:
            pass
        elif bar < 120:
            # Final drop: full percussion
            notes.append(
                {"pitch": 48, "start_beat": b + 0.5, "duration": 0.12, "velocity": 72}
            )
            notes.append(
                {"pitch": 50, "start_beat": b + 2.5, "duration": 0.12, "velocity": 65}
            )
            if bar % 2 == 1:
                notes.append(
                    {
                        "pitch": 45,
                        "start_beat": b + 1.75,
                        "duration": 0.1,
                        "velocity": 55,
                    }
                )
                notes.append(
                    {
                        "pitch": 45,
                        "start_beat": b + 3.75,
                        "duration": 0.1,
                        "velocity": 50,
                    }
                )
        else:
            if bar % 4 == 0:
                notes.append(
                    {"pitch": 48, "start_beat": b + 2, "duration": 0.15, "velocity": 40}
                )
    return notes


def build_bass():
    """Driving psytrance bassline with section variation."""
    notes = []
    for bar in range(BARS):
        b = bar * 4
        if bar < 8:
            # Ambient: no bass
            pass
        elif bar < 24:
            # Groove intro: quarter notes on E1 (40)
            for beat in range(4):
                notes.append(
                    {
                        "pitch": 40,
                        "start_beat": b + beat,
                        "duration": 0.38,
                        "velocity": 105,
                    }
                )
        elif bar < 32:
            # Breakdown 1: bass drops out
            if bar == 24:
                notes.append(
                    {"pitch": 40, "start_beat": b, "duration": 2.0, "velocity": 80}
                )
        elif bar < 56:
            # Drop 1: full 16th note drive
            for beat in range(16):
                vel = 120 if beat % 4 == 0 else 88
                if beat % 2 == 0:
                    vel += 5
                notes.append(
                    {
                        "pitch": 40,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.2,
                        "velocity": vel,
                    }
                )
        elif bar < 64:
            # Breakdown 2: single sustained note
            if bar == 56:
                notes.append(
                    {"pitch": 40, "start_beat": b, "duration": 4.0, "velocity": 75}
                )
        elif bar < 88:
            # Drop 2: 16ths with octave variation
            for beat in range(16):
                pitch = (
                    40 if (bar // 2) % 2 == 0 or beat < 12 else 52
                )  # occasional octave up
                vel = 120 if beat % 4 == 0 else 85
                notes.append(
                    {
                        "pitch": pitch,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.2,
                        "velocity": vel,
                    }
                )
        elif bar < 104:
            # Bridge: half-time bass
            for beat in range(0, 16, 2):
                notes.append(
                    {
                        "pitch": 40,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.35,
                        "velocity": 80,
                    }
                )
        elif bar < 112:
            # Breakdown 3: bass absent
            pass
        elif bar < 120:
            # Final drop: maximum drive
            for beat in range(16):
                vel = 127 if beat % 4 == 0 else 92
                notes.append(
                    {
                        "pitch": 40,
                        "start_beat": b + beat * 0.25,
                        "duration": 0.18,
                        "velocity": vel,
                    }
                )
        else:
            # Outro: bass thins to 8ths then out
            if bar < 124:
                for beat in range(8):
                    vel = max(40, 90 - (bar - 120) * 10)
                    notes.append(
                        {
                            "pitch": 40,
                            "start_beat": b + beat * 0.5,
                            "duration": 0.3,
                            "velocity": vel,
                        }
                    )
    return notes


def build_lead():
    """Arpeggiated psytrance lead with evolving patterns."""
    notes = []
    # Am pentatonic: A C D E G across octaves
    arp_intro = [69, 72, 76]  # sparse
    arp_drop1 = [69, 72, 76, 80, 76, 72]  # full arp
    arp_drop2 = [69, 72, 74, 76, 80, 76, 74, 72]  # pentatonic variation
    arp_final = [76, 80, 84, 88, 84, 80]  # high octave

    for bar in range(BARS):
        b = bar * 4
        if bar < 32:
            pass  # No lead until first drop
        elif bar < 56:
            # Drop 1: fast arpeggio
            for cycle in range(4):
                for j, p in enumerate(arp_drop1):
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + cycle * 1.0 + j * 0.125,
                            "duration": 0.125,
                            "velocity": 100,
                        }
                    )
        elif bar < 64:
            pass  # Breakdown
        elif bar < 88:
            # Drop 2: pentatonic variation with wider range
            for cycle in range(2):
                for j, p in enumerate(arp_drop2):
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + cycle * 2.0 + j * 0.2,
                            "duration": 0.15,
                            "velocity": 95,
                        }
                    )
        elif bar < 104:
            # Bridge: sparse lead
            if bar % 2 == 0:
                for i, p in enumerate(arp_intro):
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + i * 1.33,
                            "duration": 0.6,
                            "velocity": 65,
                        }
                    )
        elif bar < 112:
            pass  # Breakdown
        elif bar < 120:
            # Final drop: high energy arp
            for cycle in range(8):
                for j, p in enumerate(arp_final):
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + cycle * 0.5 + j * 0.0625,
                            "duration": 0.0625,
                            "velocity": 105,
                        }
                    )
        else:
            # Outro: fading lead
            for i, p in enumerate(arp_intro):
                vel = max(25, 60 - (bar - 120) * 5)
                notes.append(
                    {
                        "pitch": p,
                        "start_beat": b + i * 1.33,
                        "duration": 0.5,
                        "velocity": vel,
                    }
                )
    return notes


def build_arp():
    """Zappy counter-arp for stereo width."""
    notes = []
    # Different pattern from lead, offset rhythm
    pattern = [52, 55, 57, 60, 57, 55]  # E A C E(up)
    for bar in range(BARS):
        b = bar * 4
        if bar < 32:
            pass
        elif bar < 56:
            # Drop 1: offbeat arp (plays between lead notes)
            for cycle in range(2):
                for j, p in enumerate(pattern):
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + cycle * 2.0 + j * 0.25 + 0.125,
                            "duration": 0.15,
                            "velocity": 75,
                        }
                    )
        elif bar < 64:
            pass
        elif bar < 88:
            # Drop 2: different rhythm
            for j, p in enumerate(pattern[:4]):
                notes.append(
                    {
                        "pitch": p,
                        "start_beat": b + j * 0.75,
                        "duration": 0.25,
                        "velocity": 80,
                    }
                )
            if bar % 2 == 1:
                for j, p in enumerate(pattern[2:]):
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + 2 + j * 0.5,
                            "duration": 0.2,
                            "velocity": 70,
                        }
                    )
        elif bar < 104:
            # Bridge: subtle
            if bar % 4 == 0:
                for i, p in enumerate(pattern[:3]):
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + i * 1.0,
                            "duration": 0.5,
                            "velocity": 50,
                        }
                    )
        elif bar < 112:
            pass
        elif bar < 120:
            # Final drop: tight with lead
            for cycle in range(4):
                for j, p in enumerate(pattern):
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + cycle * 1.0 + j * 0.125 + 0.0625,
                            "duration": 0.1,
                            "velocity": 85,
                        }
                    )
    return notes


def build_stab():
    """Chord stabs on transitions."""
    notes = []
    # Am chord stabs: A3+C4+E4
    stab_chord = [57, 60, 64]
    # Dm chord stabs: D3+F3+A3
    stab_chord_2 = [50, 53, 57]

    for bar in range(BARS):
        b = bar * 4
        if bar < 32:
            pass
        elif bar < 56:
            # Drop 1: stabs every 4 bars on 1
            if bar % 4 == 0:
                for p in stab_chord:
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b,
                            "duration": 0.125,
                            "velocity": 110,
                        }
                    )
            # Stab on "4-and" before every 8-bar phrase
            if bar % 8 == 7:
                for p in stab_chord_2:
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + 3.75,
                            "duration": 0.125,
                            "velocity": 95,
                        }
                    )
        elif bar < 64:
            pass
        elif bar < 88:
            # Drop 2: stabs every 2 bars, alternating chords
            if bar % 2 == 0:
                chord = stab_chord if (bar // 2) % 2 == 0 else stab_chord_2
                for p in chord:
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b,
                            "duration": 0.125,
                            "velocity": 115,
                        }
                    )
        elif bar < 104:
            # Bridge: occasional stab
            if bar % 4 == 2:
                for p in stab_chord:
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + 2,
                            "duration": 0.2,
                            "velocity": 70,
                        }
                    )
        elif bar < 112:
            # Building stabs
            if bar >= 108:
                for p in stab_chord:
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + 3.75,
                            "duration": 0.1,
                            "velocity": 90,
                        }
                    )
        elif bar < 120:
            # Final drop: aggressive stabs every bar
            if bar % 2 == 0:
                for p in stab_chord:
                    notes.append(
                        {"pitch": p, "start_beat": b, "duration": 0.1, "velocity": 120}
                    )
            else:
                for p in stab_chord_2:
                    notes.append(
                        {
                            "pitch": p,
                            "start_beat": b + 2,
                            "duration": 0.1,
                            "velocity": 110,
                        }
                    )
    return notes


def build_pad():
    """Sustained Am chord pad with section dynamics."""
    notes = []
    for bar in range(BARS):
        b = bar * 4
        # Velocity by section
        if bar < 8:
            vel = 45  # ambient
        elif bar < 24:
            vel = 40  # under groove
        elif bar < 32:
            vel = 55  # breakdown swell
        elif bar < 56:
            vel = 35  # under drop (subtle)
        elif bar < 64:
            vel = 60  # breakdown swell
        elif bar < 88:
            vel = 35  # under drop
        elif bar < 104:
            vel = 55  # bridge feature
        elif bar < 112:
            vel = 65  # breakdown max
        elif bar < 120:
            vel = 30  # under final drop
        else:
            vel = max(15, 40 - (bar - 120) * 3)  # fade

        notes.append({"pitch": 57, "start_beat": b, "duration": 4.0, "velocity": vel})
        notes.append({"pitch": 60, "start_beat": b, "duration": 4.0, "velocity": vel})
        notes.append({"pitch": 64, "start_beat": b, "duration": 4.0, "velocity": vel})
    return notes


def build_dark_pad():
    """Dark atmospheric pad - lower register, different voicing."""
    notes = []
    for bar in range(BARS):
        b = bar * 4
        if bar < 8:
            vel = 40
        elif bar < 24:
            vel = 30
        elif bar < 32:
            vel = 50
        elif bar < 56:
            vel = 25
        elif bar < 64:
            vel = 50
        elif bar < 88:
            vel = 25
        elif bar < 104:
            vel = 50
        elif bar < 112:
            vel = 60
        elif bar < 120:
            vel = 20
        else:
            vel = max(10, 35 - (bar - 120) * 3)

        # Dark Am: A2 + C3 + E3 (lower voicing)
        notes.append({"pitch": 45, "start_beat": b, "duration": 4.0, "velocity": vel})
        notes.append({"pitch": 48, "start_beat": b, "duration": 4.0, "velocity": vel})
        notes.append({"pitch": 52, "start_beat": b, "duration": 4.0, "velocity": vel})
    return notes


def build_riser():
    """Riser FX - triggered at breakdown transitions."""
    notes = []
    # Breakdown 1 entry (bars 24-25): rising pitch
    for i in range(16):
        notes.append(
            {
                "pitch": 60 + i,
                "start_beat": 24 * 4 + i * 0.25,
                "duration": 0.25,
                "velocity": 60 + i * 4,
            }
        )
    # Breakdown 2 entry (bars 56-57): rising pitch
    for i in range(16):
        notes.append(
            {
                "pitch": 55 + i,
                "start_beat": 56 * 4 + i * 0.25,
                "duration": 0.25,
                "velocity": 55 + i * 4,
            }
        )
    # Final breakdown (bars 104-107): extended rise
    for i in range(48):
        vel = min(127, 40 + i * 2)
        notes.append(
            {
                "pitch": 48 + i // 4,
                "start_beat": 104 * 4 + i * 0.25,
                "duration": 0.2,
                "velocity": vel,
            }
        )
    return notes


def build_sweep():
    """Chord sweep FX at key moments."""
    notes = []
    # Intro sweep (bars 6-8)
    for i in range(8):
        notes.append(
            {
                "pitch": 60,
                "start_beat": 24 + i * 0.5,
                "duration": 0.5,
                "velocity": 40 + i * 8,
            }
        )
    # Drop 1 re-entry
    notes.append(
        {"pitch": 64, "start_beat": 31 * 4 + 2, "duration": 2.0, "velocity": 90}
    )
    # Drop 2 re-entry
    notes.append(
        {"pitch": 67, "start_beat": 63 * 4 + 2, "duration": 2.0, "velocity": 95}
    )
    # Final drop re-entry
    notes.append(
        {"pitch": 72, "start_beat": 111 * 4 + 2, "duration": 2.0, "velocity": 100}
    )
    # Outro end
    notes.append({"pitch": 55, "start_beat": 126 * 4, "duration": 4.0, "velocity": 60})
    return notes


def build_fx():
    """Impact hits and transition FX."""
    notes = []
    # Intro opening boom
    notes.append({"pitch": 36, "start_beat": 0, "duration": 2.0, "velocity": 100})
    # Groove entry hit (bar 8)
    notes.append({"pitch": 36, "start_beat": 32, "duration": 1.0, "velocity": 110})
    # Drop 1 impact (bar 32)
    notes.append({"pitch": 36, "start_beat": 128, "duration": 1.5, "velocity": 127})
    # Downlifter (bar 39-40)
    for i in range(16):
        notes.append(
            {
                "pitch": 46,
                "start_beat": 156 + i * 0.25,
                "duration": 0.25,
                "velocity": 90 - i * 4,
            }
        )
    # Drop 2 impact (bar 64)
    notes.append({"pitch": 36, "start_beat": 256, "duration": 1.5, "velocity": 127})
    # Mid-drop accent (bar 76)
    notes.append({"pitch": 36, "start_beat": 304, "duration": 0.5, "velocity": 110})
    # Bridge transition (bar 88)
    notes.append({"pitch": 36, "start_beat": 352, "duration": 1.0, "velocity": 95})
    # Final drop impact (bar 112)
    notes.append({"pitch": 36, "start_beat": 448, "duration": 2.0, "velocity": 127})
    # Outro final hit (bar 126)
    notes.append({"pitch": 36, "start_beat": 504, "duration": 2.0, "velocity": 80})
    # White noise riser textures (pitch 80 = shaker zone)
    for i in range(8):
        notes.append(
            {
                "pitch": 80,
                "start_beat": 120 + i * 2,
                "duration": 0.5,
                "velocity": 30 + i * 5,
            }
        )
    return notes


# ============================================================
# EFFECT TRANSITION ENGINE
# ============================================================


def run_effect_transitions(indices):
    """Complex multi-track effect automation synchronized to song structure."""
    sec = 60.0 / BPM  # seconds per beat

    kick_idx, hatcl_idx, hatop_idx, snare_idx, perc_idx = indices[0:5]
    bass_idx, lead_idx, arp_idx, stab_idx = indices[5:9]
    pad_idx, dark_idx, riser_idx, sweep_idx, fx_idx = indices[9:14]

    # Section times in seconds
    t = {}
    for name, (start_bar, end_bar) in SECTIONS.items():
        t[name + "_start"] = start_bar * 4 * sec
        t[name + "_end"] = end_bar * 4 * sec

    schedule = []

    # ===== BASS: Double filter sweeps (Auto Filter + instrument cutoff) =====
    # Open during drops, closed during breakdowns
    # Ambient/Groove: open
    schedule.append(
        (
            0,
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(1.0),
        )
    )
    schedule.append(
        (
            0,
            "/superdaw/plugin/param",
            bass_idx,
            "Analog Saw Bass",
            "Filter Cutoff",
            107.0,
        )
    )
    # Breakdown 1: close
    schedule.append(
        (
            t["breakdown_1_start"] + 4 * sec,
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(0.5),
        )
    )
    schedule.append(
        (
            t["breakdown_1_start"] + 12 * sec,
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(0.1),
        )
    )
    schedule.append(
        (
            t["breakdown_1_start"] + 12 * sec,
            "/superdaw/plugin/param",
            bass_idx,
            "Analog Saw Bass",
            "Filter Cutoff",
            25.0,
        )
    )
    # Drop 1: open
    schedule.append(
        (
            t["drop_1_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(1.0),
        )
    )
    schedule.append(
        (
            t["drop_1_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Analog Saw Bass",
            "Filter Cutoff",
            127.0,
        )
    )
    # Breakdown 2: close
    schedule.append(
        (
            t["breakdown_2_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(0.1),
        )
    )
    schedule.append(
        (
            t["breakdown_2_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Analog Saw Bass",
            "Filter Cutoff",
            20.0,
        )
    )
    # Drop 2: open
    schedule.append(
        (
            t["drop_2_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(1.0),
        )
    )
    schedule.append(
        (
            t["drop_2_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Analog Saw Bass",
            "Filter Cutoff",
            127.0,
        )
    )
    # Bridge: mid
    schedule.append(
        (
            t["bridge_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(0.5),
        )
    )
    schedule.append(
        (
            t["bridge_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Analog Saw Bass",
            "Filter Cutoff",
            60.0,
        )
    )
    # Breakdown 3: close
    schedule.append(
        (
            t["breakdown_3_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(0.08),
        )
    )
    schedule.append(
        (
            t["breakdown_3_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Analog Saw Bass",
            "Filter Cutoff",
            15.0,
        )
    )
    # Final drop: open
    schedule.append(
        (
            t["drop_3_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(1.0),
        )
    )
    schedule.append(
        (
            t["drop_3_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Analog Saw Bass",
            "Filter Cutoff",
            127.0,
        )
    )
    # Outro: gradual close
    for i in range(4):
        tt = t["outro_start"] + i * (8 * sec)
        val = 1.0 - i * 0.25
        schedule.append(
            (
                tt,
                "/superdaw/plugin/param",
                bass_idx,
                "Auto Filter",
                "Frequency",
                filter_hz(max(val, 0.05)),
            )
        )
        schedule.append(
            (
                tt,
                "/superdaw/plugin/param",
                bass_idx,
                "Analog Saw Bass",
                "Filter Cutoff",
                max(10.0, val * 127.0),
            )
        )

    # ===== BASS: Saturator drive automation =====
    schedule.append((0, "/superdaw/plugin/param", bass_idx, "Saturator", "Drive", 5.0))
    schedule.append(
        (
            t["drop_1_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Saturator",
            "Drive",
            15.0,
        )
    )
    schedule.append(
        (
            t["drop_2_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Saturator",
            "Drive",
            18.0,
        )
    )
    schedule.append(
        (
            t["bridge_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Saturator",
            "Drive",
            3.0,
        )
    )
    schedule.append(
        (
            t["drop_3_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Saturator",
            "Drive",
            22.0,
        )
    )
    schedule.append(
        (
            t["outro_start"],
            "/superdaw/plugin/param",
            bass_idx,
            "Saturator",
            "Drive",
            2.0,
        )
    )

    # ===== LEAD: Reverb + Delay automation =====
    for band in ["Mid Dry/Wet", "High Dry/Wet"]:
        # Intro: medium
        schedule.append(
            (0, "/superdaw/plugin/param", lead_idx, "Reverb", band, band_wet(0.4))
        )
        # Breakdown 1: swell
        schedule.append(
            (
                t["breakdown_1_start"] + 8 * sec,
                "/superdaw/plugin/param",
                lead_idx,
                "Reverb",
                band,
                band_wet(0.9),
            )
        )
        # Drop 1: dry
        schedule.append(
            (
                t["drop_1_start"],
                "/superdaw/plugin/param",
                lead_idx,
                "Reverb",
                band,
                band_wet(0.08),
            )
        )
        # Breakdown 2: max
        schedule.append(
            (
                t["breakdown_2_start"],
                "/superdaw/plugin/param",
                lead_idx,
                "Reverb",
                band,
                band_wet(1.0),
            )
        )
        # Drop 2: dry
        schedule.append(
            (
                t["drop_2_start"],
                "/superdaw/plugin/param",
                lead_idx,
                "Reverb",
                band,
                band_wet(0.1),
            )
        )
        # Bridge: wet
        schedule.append(
            (
                t["bridge_start"],
                "/superdaw/plugin/param",
                lead_idx,
                "Reverb",
                band,
                band_wet(0.7),
            )
        )
        # Final drop: dry
        schedule.append(
            (
                t["drop_3_start"],
                "/superdaw/plugin/param",
                lead_idx,
                "Reverb",
                band,
                band_wet(0.05),
            )
        )
        # Outro: max
        schedule.append(
            (
                t["outro_start"] + 8 * sec,
                "/superdaw/plugin/param",
                lead_idx,
                "Reverb",
                band,
                band_wet(1.0),
            )
        )

    # Lead delay throws
    schedule.append(
        (0, "/superdaw/plugin/param", lead_idx, "Delay", "Rack Dry/Wet", rack_wet(0.0))
    )
    # Delay throw at end of groove intro
    schedule.append(
        (
            t["groove_intro_end"] - 4 * sec,
            "/superdaw/plugin/param",
            lead_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.7),
        )
    )
    schedule.append(
        (
            t["groove_intro_end"],
            "/superdaw/plugin/param",
            lead_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.0),
        )
    )
    # Drop 1: subtle delay
    schedule.append(
        (
            t["drop_1_start"],
            "/superdaw/plugin/param",
            lead_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.1),
        )
    )
    # Delay throw before drop 2
    schedule.append(
        (
            t["breakdown_2_end"] - 4 * sec,
            "/superdaw/plugin/param",
            lead_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.6),
        )
    )
    schedule.append(
        (
            t["drop_2_start"],
            "/superdaw/plugin/param",
            lead_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.1),
        )
    )
    # Outro: big delay
    schedule.append(
        (
            t["outro_start"],
            "/superdaw/plugin/param",
            lead_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.5),
        )
    )
    schedule.append(
        (
            t["outro_start"] + 16 * sec,
            "/superdaw/plugin/param",
            lead_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.9),
        )
    )

    # ===== ARP: Reverb + Auto Filter =====
    schedule.append(
        (
            t["drop_1_start"],
            "/superdaw/plugin/param",
            arp_idx,
            "Reverb",
            "Mid Dry/Wet",
            band_wet(0.3),
        )
    )
    schedule.append(
        (
            t["breakdown_2_start"],
            "/superdaw/plugin/param",
            arp_idx,
            "Reverb",
            "Mid Dry/Wet",
            band_wet(0.8),
        )
    )
    schedule.append(
        (
            t["drop_2_start"],
            "/superdaw/plugin/param",
            arp_idx,
            "Reverb",
            "Mid Dry/Wet",
            band_wet(0.25),
        )
    )
    # ARP filter sweeps
    schedule.append(
        (
            t["drop_1_start"],
            "/superdaw/plugin/param",
            arp_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(1.0),
        )
    )
    schedule.append(
        (
            t["breakdown_2_start"],
            "/superdaw/plugin/param",
            arp_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(0.2),
        )
    )
    schedule.append(
        (
            t["drop_2_start"],
            "/superdaw/plugin/param",
            arp_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(1.0),
        )
    )
    schedule.append(
        (
            t["bridge_start"],
            "/superdaw/plugin/param",
            arp_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(0.4),
        )
    )
    schedule.append(
        (
            t["drop_3_start"],
            "/superdaw/plugin/param",
            arp_idx,
            "Auto Filter",
            "Frequency",
            filter_hz(1.0),
        )
    )

    # ===== PAD + DARK PAD: Reverb swells =====
    for band in ["Mid Dry/Wet", "High Dry/Wet"]:
        schedule.append(
            (0, "/superdaw/plugin/param", pad_idx, "Reverb", band, band_wet(0.8))
        )
        schedule.append(
            (
                t["drop_1_start"],
                "/superdaw/plugin/param",
                pad_idx,
                "Reverb",
                band,
                band_wet(0.3),
            )
        )
        schedule.append(
            (
                t["breakdown_2_start"],
                "/superdaw/plugin/param",
                pad_idx,
                "Reverb",
                band,
                band_wet(0.95),
            )
        )
        schedule.append(
            (
                t["drop_2_start"],
                "/superdaw/plugin/param",
                pad_idx,
                "Reverb",
                band,
                band_wet(0.3),
            )
        )
        schedule.append(
            (
                t["bridge_start"],
                "/superdaw/plugin/param",
                pad_idx,
                "Reverb",
                band,
                band_wet(0.85),
            )
        )
        schedule.append(
            (
                t["outro_start"],
                "/superdaw/plugin/param",
                pad_idx,
                "Reverb",
                band,
                band_wet(1.0),
            )
        )

        schedule.append(
            (0, "/superdaw/plugin/param", dark_idx, "Reverb", band, band_wet(0.7))
        )
        schedule.append(
            (
                t["drop_1_start"],
                "/superdaw/plugin/param",
                dark_idx,
                "Reverb",
                band,
                band_wet(0.2),
            )
        )
        schedule.append(
            (
                t["breakdown_3_start"],
                "/superdaw/plugin/param",
                dark_idx,
                "Reverb",
                band,
                band_wet(1.0),
            )
        )
        schedule.append(
            (
                t["outro_start"],
                "/superdaw/plugin/param",
                dark_idx,
                "Reverb",
                band,
                band_wet(1.0),
            )
        )

    # DARK PAD delay for atmosphere
    schedule.append(
        (0, "/superdaw/plugin/param", dark_idx, "Delay", "Rack Dry/Wet", rack_wet(0.3))
    )
    schedule.append(
        (
            t["drop_1_start"],
            "/superdaw/plugin/param",
            dark_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.05),
        )
    )
    schedule.append(
        (
            t["bridge_start"],
            "/superdaw/plugin/param",
            dark_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.5),
        )
    )
    schedule.append(
        (
            t["outro_start"],
            "/superdaw/plugin/param",
            dark_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.7),
        )
    )

    # ===== DRUM VOLUME FADES during breakdowns =====
    # Breakdown 1: drums fade
    schedule.append(
        (t["breakdown_1_start"] + 8 * sec, "/superdaw/track/volume", kick_idx, 0.4)
    )
    schedule.append(
        (t["breakdown_1_start"] + 16 * sec, "/superdaw/track/volume", kick_idx, 0.15)
    )
    schedule.append(
        (t["breakdown_1_start"] + 8 * sec, "/superdaw/track/volume", hatcl_idx, 0.2)
    )
    schedule.append(
        (t["breakdown_1_start"] + 16 * sec, "/superdaw/track/volume", hatcl_idx, 0.05)
    )
    schedule.append(
        (t["breakdown_1_start"] + 8 * sec, "/superdaw/track/volume", snare_idx, 0.3)
    )
    # Pre-drop: snap back
    schedule.append(
        (t["drop_1_start"] - 2 * sec, "/superdaw/track/volume", kick_idx, 0.92)
    )
    schedule.append(
        (t["drop_1_start"] - 2 * sec, "/superdaw/track/volume", hatcl_idx, 0.62)
    )
    schedule.append(
        (t["drop_1_start"] - 2 * sec, "/superdaw/track/volume", snare_idx, 0.78)
    )

    # Breakdown 2: all drums drop
    schedule.append((t["breakdown_2_start"], "/superdaw/track/volume", kick_idx, 0.1))
    schedule.append((t["breakdown_2_start"], "/superdaw/track/volume", hatcl_idx, 0.05))
    schedule.append((t["breakdown_2_start"], "/superdaw/track/volume", snare_idx, 0.1))
    schedule.append((t["breakdown_2_start"], "/superdaw/track/volume", perc_idx, 0.0))
    schedule.append(
        (t["drop_2_start"] - 2 * sec, "/superdaw/track/volume", kick_idx, 0.92)
    )
    schedule.append(
        (t["drop_2_start"] - 2 * sec, "/superdaw/track/volume", hatcl_idx, 0.62)
    )
    schedule.append(
        (t["drop_2_start"] - 2 * sec, "/superdaw/track/volume", snare_idx, 0.78)
    )
    schedule.append(
        (t["drop_2_start"] - 2 * sec, "/superdaw/track/volume", perc_idx, 0.58)
    )

    # Bridge: drums at half
    schedule.append((t["bridge_start"], "/superdaw/track/volume", kick_idx, 0.55))
    schedule.append((t["bridge_start"], "/superdaw/track/volume", hatcl_idx, 0.35))

    # Breakdown 3: drums gone
    schedule.append((t["breakdown_3_start"], "/superdaw/track/volume", kick_idx, 0.0))
    schedule.append((t["breakdown_3_start"], "/superdaw/track/volume", hatcl_idx, 0.1))
    schedule.append((t["breakdown_3_start"], "/superdaw/track/volume", snare_idx, 0.1))
    # Final drop: everything back
    schedule.append(
        (t["drop_3_start"] - 1 * sec, "/superdaw/track/volume", kick_idx, 0.92)
    )
    schedule.append(
        (t["drop_3_start"] - 1 * sec, "/superdaw/track/volume", hatcl_idx, 0.62)
    )
    schedule.append(
        (t["drop_3_start"] - 1 * sec, "/superdaw/track/volume", snare_idx, 0.78)
    )
    schedule.append(
        (t["drop_3_start"] - 1 * sec, "/superdaw/track/volume", perc_idx, 0.58)
    )

    # Outro: everything fades
    for i in range(4):
        tt = t["outro_start"] + i * (8 * sec)
        vol = 0.92 - i * 0.2
        schedule.append((tt, "/superdaw/track/volume", kick_idx, max(vol, 0.0)))
        schedule.append((tt, "/superdaw/track/volume", hatcl_idx, max(vol * 0.7, 0.0)))
        schedule.append((tt, "/superdaw/track/volume", snare_idx, max(vol * 0.85, 0.0)))
        schedule.append((tt, "/superdaw/track/volume", perc_idx, max(vol * 0.6, 0.0)))

    # ===== HAT OPEN: Delay automation =====
    schedule.append(
        (0, "/superdaw/plugin/param", hatop_idx, "Delay", "Rack Dry/Wet", rack_wet(0.0))
    )
    schedule.append(
        (
            t["drop_1_start"] + 16 * sec,
            "/superdaw/plugin/param",
            hatop_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.3),
        )
    )
    schedule.append(
        (
            t["drop_2_start"] + 16 * sec,
            "/superdaw/plugin/param",
            hatop_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.4),
        )
    )
    schedule.append(
        (
            t["outro_start"],
            "/superdaw/plugin/param",
            hatop_idx,
            "Delay",
            "Rack Dry/Wet",
            rack_wet(0.6),
        )
    )

    # ===== STAB: Reverb automation =====
    for band in ["Mid Dry/Wet", "High Dry/Wet"]:
        schedule.append(
            (
                t["drop_1_start"],
                "/superdaw/plugin/param",
                stab_idx,
                "Reverb",
                band,
                band_wet(0.15),
            )
        )
        schedule.append(
            (
                t["breakdown_2_start"],
                "/superdaw/plugin/param",
                stab_idx,
                "Reverb",
                band,
                band_wet(0.5),
            )
        )
        schedule.append(
            (
                t["drop_2_start"],
                "/superdaw/plugin/param",
                stab_idx,
                "Reverb",
                band,
                band_wet(0.1),
            )
        )
        schedule.append(
            (
                t["drop_3_start"],
                "/superdaw/plugin/param",
                stab_idx,
                "Reverb",
                band,
                band_wet(0.08),
            )
        )

    # ===== RISER + SWEEP: Reverb always high =====
    for band in ["Mid Dry/Wet", "High Dry/Wet"]:
        schedule.append(
            (0, "/superdaw/plugin/param", riser_idx, "Reverb", band, band_wet(0.7))
        )
        schedule.append(
            (0, "/superdaw/plugin/param", sweep_idx, "Reverb", band, band_wet(0.6))
        )

    # SWEEP delay for trailing effect
    schedule.append(
        (0, "/superdaw/plugin/param", sweep_idx, "Delay", "Rack Dry/Wet", rack_wet(0.4))
    )

    # Sort by time
    schedule.sort(key=lambda x: x[0])

    start_time = time.time()

    def run_schedule():
        for entry in schedule:
            sched_time = entry[0]
            osc_addr = entry[1]
            osc_args = entry[2:]
            wait = sched_time - (time.time() - start_time)
            if wait > 0:
                time.sleep(wait)
            send_osc(osc_addr, *osc_args)

    thread = threading.Thread(target=run_schedule)
    thread.daemon = True
    thread.start()
    return thread


# ============================================================
# MAIN
# ============================================================


def main():
    log("=" * 60)
    log("  FULL-ON PSYTRANCE PRODUCTION BUILDER")
    log("  128 bars @ 148 BPM (~8 min)")
    log("=" * 60)

    # STEP 1: Create all tracks
    log("Creating %d tracks..." % len(TRACK_LAYOUT))
    track_indices = []
    base_idx = 4  # assume 4 default tracks
    for name, preset, vol, pan, mute_bars in TRACK_LAYOUT:
        send_osc("/superdaw/track/create", name, "midi")
        time.sleep(0.5)
        idx_str = read_temp("superdaw_track_created.txt", timeout=2.0)
        idx = int(idx_str) if idx_str else (base_idx + len(track_indices))
        track_indices.append(idx)
        log("  Track %d: %s (preset: %s)" % (idx, name, preset))

    # STEP 2: Set tempo
    log("Setting tempo to %d BPM..." % BPM)
    send_osc("/superdaw/transport/tempo", float(BPM))
    time.sleep(0.5)

    # STEP 3: Load instrument presets
    log("Loading instrument presets...")
    for i, (name, preset, vol, pan, mute_bars) in enumerate(TRACK_LAYOUT):
        log("  Loading '%s' on %s (track %d)..." % (preset, name, track_indices[i]))
        send_osc("/superdaw/track/instrument", track_indices[i], preset)
        time.sleep(5)

    # STEP 4: Load audio effects
    log("Loading audio effects...")
    # Map track names to indices
    name_to_idx = {}
    short_names = [
        "KICK",
        "HATCL",
        "HATOP",
        "SNARE",
        "PERC",
        "BASS",
        "LEAD",
        "ARP",
        "STAB",
        "PAD",
        "DARK",
        "RISER",
        "SWEEP",
        "FX",
    ]
    for i, sn in enumerate(short_names):
        name_to_idx[sn] = track_indices[i]

    for track_key, effects in EFFECT_CHAINS.items():
        tidx = name_to_idx[track_key]
        for fx_name in effects:
            log("  Loading '%s' on %s (track %d)..." % (fx_name, track_key, tidx))
            send_osc("/superdaw/track/effect", tidx, fx_name)
            time.sleep(5)
    log("All effects loaded")

    # STEP 5: Write MIDI clips
    log("Writing MIDI clips...")
    clip_builders = [
        (0, build_kick, "Kick"),
        (1, build_hat_closed, "Hat Closed"),
        (2, build_hat_open, "Hat Open"),
        (3, build_snare, "Snare"),
        (4, build_perc, "Percussion"),
        (5, build_bass, "Bass"),
        (6, build_lead, "Lead"),
        (7, build_arp, "Arp"),
        (8, build_stab, "Stab"),
        (9, build_pad, "Pad"),
        (10, build_dark_pad, "Dark Pad"),
        (11, build_riser, "Riser"),
        (12, build_sweep, "Sweep"),
        (13, build_fx, "FX"),
    ]
    for track_offset, builder, name in clip_builders:
        notes = builder()
        tidx = track_indices[track_offset]
        write_clip(tidx, 0, notes)
        log("  %s: %d notes on track %d" % (name, len(notes), tidx))

    # STEP 6: Set mix levels + panning
    log("Setting mix levels and panning...")
    for i, (name, preset, vol, pan, mute_bars) in enumerate(TRACK_LAYOUT):
        send_osc("/superdaw/track/volume", track_indices[i], vol)
        send_osc("/superdaw/track/pan", track_indices[i], pan)
    time.sleep(0.5)

    # STEP 7: Launch all clips
    log("Launching all clips...")
    for idx in track_indices:
        send_osc("/superdaw/clip/launch", idx, 0)
        time.sleep(0.3)

    # STEP 8: Start playback + effect automation
    send_osc("/superdaw/transport/play", 1)
    log("Starting effect automation thread...")
    fx_thread = run_effect_transitions(track_indices)

    log("=" * 60)
    log("  PLAYING FULL-ON PSYTRANCE PRODUCTION!")
    log("=" * 60)
    log("  Bars 1-8:    Ambient Intro")
    log("  Bars 9-24:   Groove Intro")
    log("  Bars 25-32:  Breakdown 1")
    log("  Bars 33-56:  DROP 1")
    log("  Bars 57-64:  Breakdown 2")
    log("  Bars 65-88:  DROP 2")
    log("  Bars 89-104: Bridge")
    log("  Bars 105-112:Breakdown 3")
    log("  Bars 113-120:FINAL DROP")
    log("  Bars 121-128:Outro")
    log("=" * 60)

    # Play for 60 seconds to hear major transitions
    play_time = 60
    log("Listening for %d seconds..." % play_time)
    time.sleep(play_time)

    send_osc("/superdaw/transport/play", 0)
    log("Stopped. Full production remains in Ableton!")


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print("[!] Error: " + str(e))
        import traceback

        traceback.print_exc()
        sys.exit(1)
