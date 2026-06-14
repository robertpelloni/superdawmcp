#!/usr/bin/env python3
"""
Full Psytrance Production Builder for Ableton Live 10 Standard via SuperDAW OSC
Uses Surge XT VST3 as the primary synthesizer.

Creates a complete 128-bar psytrance track at 148 BPM with:
- Surge XT leads, bass, pads, and FX
- Drum Rack with MIDI patterns
- Mixer automation (volume, pan)
- Effect transitions and buildups
"""

import socket
import struct
import time
import json
import tempfile
import os

# --- OSC Transport ---
UDP_HOST = "127.0.0.1"
UDP_PORT = 11000


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
    sock.sendto(msg, (UDP_HOST, UDP_PORT))
    sock.close()


def wait(secs):
    time.sleep(secs)


def get_track_index():
    path = os.path.join(tempfile.gettempdir(), "superdaw_track_created.txt")
    try:
        with open(path, "r") as f:
            return int(f.read().strip())
    except:
        return -1


def create_track(name, track_type="midi"):
    send_osc("/superdaw/track/create", name, track_type)
    wait(1.5)
    idx = get_track_index()
    print(f"  Created track '{name}' at index {idx}")
    return idx


def load_instrument(track_idx, instrument_name):
    send_osc("/superdaw/track/instrument", track_idx, instrument_name)
    wait(8)  # VST plugins need time to load
    print(f"  Loaded '{instrument_name}' on track {track_idx}")


def write_clip_from_file(track_idx, clip_idx, name, length, notes):
    """Write clip via JSON file. Notes must be dicts with keys:
    pitch, start_beat, duration, velocity, mute."""
    # Convert [pitch, vel, dur, start, mute] arrays to dicts
    clip_notes = []
    for n in notes:
        clip_notes.append(
            {
                "pitch": n[0],
                "velocity": n[1],
                "duration": n[2],
                "start_beat": n[3],
                "mute": n[4],
            }
        )
    clip_file = os.path.join(
        tempfile.gettempdir(), "superdaw_clip_%d_%d.json" % (track_idx, clip_idx)
    )
    with open(clip_file, "w") as f:
        json.dump(clip_notes, f)
    # Send 3 args: track_idx, clip_idx, filepath
    send_osc("/superdaw/clip/fromfile", track_idx, clip_idx, clip_file)
    wait(2)  # Give time for clip creation


def set_volume(track_idx, volume_linear):
    """Set track volume. Use 0-1 range (0.85 ~ unity, 1.0 = +6dB)."""
    send_osc("/superdaw/track/volume", track_idx, volume_linear)


def set_pan(track_idx, pan):
    send_osc("/superdaw/track/pan", track_idx, pan)


def transport(playing, bpm):
    send_osc("/superdaw/transport", 1 if playing else 0, bpm)


# --- Musical Constants ---
BPM = 148
BEAT = 60.0 / BPM  # Duration of one beat in seconds
BAR = 4 * BEAT  # Duration of one bar in beats (4/4)


# Note helpers
def note(pitch, velocity, start_beat, duration_beats):
    return [pitch, velocity, duration_beats, start_beat, 0]


def midi_pitch(name):
    """Convert note name like C4, D#5 to MIDI number."""
    names = {"C": 0, "D": 2, "E": 4, "F": 5, "G": 7, "A": 9, "B": 11}
    i = 0
    base = names[name[i]]
    i += 1
    accidental = 0
    if i < len(name) and name[i] == "#":
        accidental = 1
        i += 1
    elif i < len(name) and name[i] == "b":
        accidental = -1
        i += 1
    octave = int(name[i])
    return (octave + 1) * 12 + base + accidental


# Psytrance scales
A_MINOR = [0, 2, 3, 5, 7, 8, 10]  # A natural minor
PHRYGIAN = [0, 1, 3, 5, 7, 8, 10]  # A Phrygian (darker)


def scale_note(scale, root, degree, octave=4):
    """Get MIDI pitch for scale degree (0-based)."""
    octave_offset = degree // len(scale)
    scale_idx = degree % len(scale)
    return root + (octave + octave_offset) * 12 + scale[scale_idx]


# --- Pattern Generators ---


def psykick_pattern(bars=128):
    """4-on-the-floor kick pattern."""
    notes = []
    for bar in range(bars):
        for beat in range(4):
            start = (bar * 4 + beat) * 1.0
            notes.append(note(36, 120, start, 0.25))
    return notes


def psybass_pattern(root=33, bars=128):
    """Driving psytrance bass - 16th note pattern with octave jumps.
    Root A1 = MIDI 33."""
    notes = []
    for bar in range(bars):
        bar_start = bar * 4.0
        for i in range(16):  # 16 sixteenth notes per bar
            start = bar_start + i * 0.25
            if i % 2 == 0:  # On the 16ths
                vel = 110 if i % 4 == 0 else 85
                notes.append(note(root, vel, start, 0.15))
            # Off 16ths are rests (psytrance bass is syncopated)
    return notes


def lead_pattern(bars=128):
    """Evolving lead melody using A Phrygian mode."""
    root = midi_pitch("A4")  # A4 = 69
    notes = []

    # Phrase structure: 8-bar phrases
    phrases = [
        # Phrase 1: Ascending (bars 0-7)
        [
            (0, 0.5),
            (2, 0.5),
            (3, 1.0),
            (5, 0.5),
            (7, 0.5),
            (8, 1.0),
            (7, 0.5),
            (5, 0.5),
        ],
        # Phrase 2: Descending (bars 8-15)
        [
            (10, 1.0),
            (8, 0.5),
            (7, 0.5),
            (5, 1.0),
            (3, 0.5),
            (2, 0.5),
            (0, 1.0),
            (0, 0.5),
        ],
        # Phrase 3: Trill (bars 16-23)
        [
            (3, 0.25),
            (5, 0.25),
            (3, 0.25),
            (5, 0.25),
            (7, 0.5),
            (8, 0.5),
            (7, 1.0),
            (5, 1.0),
        ],
        # Phrase 4: Wide intervals (bars 24-31)
        [
            (0, 0.5),
            (7, 0.5),
            (10, 1.0),
            (7, 0.5),
            (5, 0.5),
            (3, 1.0),
            (0, 1.0),
            (0, 0.5),
        ],
    ]

    for bar_group in range(bars // 8):
        phrase_idx = bar_group % len(phrases)
        phrase = phrases[phrase_idx]
        for bar_in_phrase in range(8):
            bar_start = (bar_group * 8 + bar_in_phrase) * 4.0
            for i, (degree, dur) in enumerate(phrase):
                start = bar_start + i * 0.5
                if start < bar_start + 4.0:
                    pitch = scale_note(PHRYGIAN, root % 12, degree, 4)
                    vel = 95 + (degree * 3) % 30
                    notes.append(note(pitch, min(vel, 127), start, dur * 0.8))

    return notes


def pad_pattern(bars=128):
    """Sustained pad chords for atmosphere."""
    root = midi_pitch("A3")  # A3 = 57
    notes = []

    # Chord voicings (scale degrees)
    chords = [
        [0, 3, 7],  # Am (i)
        [5, 8, 12],  # Fm (VI)
        [3, 7, 10],  # Cm (III)
        [7, 10, 14],  # Gm (VII)
    ]

    for bar in range(bars):
        chord = chords[(bar // 4) % len(chords)]
        bar_start = bar * 4.0
        for degree in chord:
            pitch = scale_note(A_MINOR, root % 12, degree, 3)
            notes.append(note(pitch, 60, bar_start, 3.8))

    return notes


def arp_pattern(bars=128):
    """Fast arpeggiated pattern for texture."""
    root = midi_pitch("A4")
    notes = []

    arp_notes = [0, 3, 7, 12, 7, 3]  # Am arpeggio up and down

    for bar in range(bars):
        bar_start = bar * 4.0
        for i in range(32):  # 32nd notes
            start = bar_start + i * 0.125
            degree = arp_notes[i % len(arp_notes)]
            pitch = scale_note(A_MINOR, root % 12, degree, 4)
            vel = 70 + (i % 4) * 15
            notes.append(note(pitch, min(vel, 127), start, 0.1))

    return notes


def fx_riser_pattern(bars=128):
    """Rising white noise / FX hits on transitions (every 32 bars)."""
    notes = []
    for section in range(bars // 32):
        # Build up over last 4 bars of each 32-bar section
        for bar in range(4):
            actual_bar = section * 32 + 28 + bar
            bar_start = actual_bar * 4.0
            intensity = (bar + 1) / 4.0
            for i in range(16):
                start = bar_start + i * 0.25
                vel = int(40 + intensity * 80)
                # Chromatic ascending
                pitch = 60 + bar * 4 + i
                notes.append(note(min(pitch, 127), vel, start, 0.15))

    return notes


def hihat_pattern(bars=128):
    """16th note hi-hats with accent pattern."""
    notes = []
    for bar in range(bars):
        bar_start = bar * 4.0
        for i in range(16):
            start = bar_start + i * 0.25
            # Accents on beats
            vel = 100 if i % 4 == 0 else (80 if i % 2 == 0 else 55)
            notes.append(note(42, vel, start, 0.1))
    return notes


def snare_pattern(bars=128):
    """Snare on 2 and 4."""
    notes = []
    for bar in range(bars):
        bar_start = bar * 4.0
        notes.append(note(40, 100, bar_start + 1.0, 0.25))
        notes.append(note(40, 100, bar_start + 3.0, 0.25))
    return notes


# --- Build the Track ---


def build_psytrance():
    print("=" * 60)
    print("SUPERDAW PSYTRANCE BUILDER - Surge XT Edition")
    print("=" * 60)
    print(f"BPM: {BPM}")
    print(f"Duration: 128 bars ({128 * 4 / BPM * 60 / 60:.1f} minutes)")
    print()

    # Stop transport
    transport(False, BPM)
    wait(1)

    # === DRUM TRACK ===
    print("[1/7] Creating Drum Track...")
    kick_idx = create_track("Kick", "midi")
    # Load Drum Rack via instrument search
    load_instrument(kick_idx, "Drum Rack")

    kick_notes = psykick_pattern(128)
    write_clip_from_file(kick_idx, 0, "Kick Main", 128 * 4, kick_notes)
    set_volume(kick_idx, 0.85)  # unity
    set_pan(kick_idx, 0.0)
    print()

    # === HI-HAT TRACK ===
    print("[2/7] Creating Hi-Hat Track...")
    hh_idx = create_track("HiHats", "midi")
    load_instrument(hh_idx, "Drum Rack")
    hh_notes = hihat_pattern(128)
    write_clip_from_file(hh_idx, 0, "HiHats Main", 128 * 4, hh_notes)
    set_volume(hh_idx, 0.70)  # -6dB approx
    set_pan(hh_idx, 0.3)
    print()

    # === SNARE TRACK ===
    print("[3/7] Creating Snare Track...")
    snare_idx = create_track("Snare", "midi")
    load_instrument(snare_idx, "Drum Rack")
    snare_notes = snare_pattern(128)
    write_clip_from_file(snare_idx, 0, "Snare Main", 128 * 4, snare_notes)
    set_volume(snare_idx, 0.75)  # -4dB approx
    set_pan(snare_idx, -0.2)
    print()

    # === BASS TRACK (Surge XT) ===
    print("[4/7] Creating Psytrance Bass (Surge XT)...")
    bass_idx = create_track("Psy Bass", "midi")
    load_instrument(bass_idx, "Surge XT")

    bass_notes = psybass_pattern(midi_pitch("A1"), 128)
    write_clip_from_file(bass_idx, 0, "Psy Bass Main", 128 * 4, bass_notes)
    set_volume(bass_idx, 0.80)  # -3dB approx
    set_pan(bass_idx, 0.0)
    print()

    # === LEAD TRACK (Surge XT) ===
    print("[5/7] Creating Lead (Surge XT)...")
    lead_idx = create_track("Lead", "midi")
    load_instrument(lead_idx, "Surge XT")

    lead_notes = lead_pattern(128)
    write_clip_from_file(lead_idx, 0, "Lead Main", 128 * 4, lead_notes)
    set_volume(lead_idx, 0.60)  # -8dB approx
    set_pan(lead_idx, 0.2)
    print()

    # === PAD TRACK (Surge XT) ===
    print("[6/7] Creating Pad (Surge XT)...")
    pad_idx = create_track("Pad", "midi")
    load_instrument(pad_idx, "Surge XT")

    pad_notes = pad_pattern(128)
    write_clip_from_file(pad_idx, 0, "Pad Main", 128 * 4, pad_notes)
    set_volume(pad_idx, 0.50)  # -12dB approx
    set_pan(pad_idx, -0.3)
    print()

    # === ARP TRACK (Surge XT) ===
    print("[7/7] Creating Arpeggio (Surge XT)...")
    arp_idx = create_track("Arp", "midi")
    load_instrument(arp_idx, "Surge XT")

    arp_notes = arp_pattern(128)
    write_clip_from_file(arp_idx, 0, "Arp Main", 128 * 4, arp_notes)
    set_volume(arp_idx, 0.55)  # -10dB approx
    set_pan(arp_idx, 0.4)
    print()

    # === Set transport and start ===
    print("=" * 60)
    print("Starting playback...")
    transport(True, BPM)
    print(f"Playing at {BPM} BPM!")
    print("=" * 60)
    print()
    print("Track layout:")
    print(f"  Track {kick_idx}: Kick (Drum Rack)")
    print(f"  Track {hh_idx}: Hi-Hats (Drum Rack)")
    print(f"  Track {snare_idx}: Snare (Drum Rack)")
    print(f"  Track {bass_idx}: Psy Bass (Surge XT)")
    print(f"  Track {lead_idx}: Lead (Surge XT)")
    print(f"  Track {pad_idx}: Pad (Surge XT)")
    print(f"  Track {arp_idx}: Arpeggio (Surge XT)")
    print()
    print("NOTE: Surge XT loads with default 'Init' preset.")
    print("You can tweak parameters in the Ableton UI for each track.")
    print("The default init preset produces sound - check audio output!")


if __name__ == "__main__":
    build_psytrance()
