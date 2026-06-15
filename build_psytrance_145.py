#!/usr/bin/env python3
"""
Full Psytrance Track Builder - 145 BPM
Comprehensive psytrance production for Ableton Live 10 Standard via SuperDAW OSC.

Creates a complete 128-bar psytrance track at 145 BPM with:
- Drum Rack: Kick, HiHats, Snare/Clap, Cymbals
- Surge XT: Psy Bass, Lead, Pad, Arp, FX Riser
- Song structure: Intro -> Build-up -> Drop -> Breakdown -> Rise -> Drop -> Outro
- Automation: Filter sweeps, volume, reverb, delay builds
- Return tracks: Hall Reverb, Ping Pong Delay
- Master chain: EQ8 + Limiter
"""

import socket
import struct
import time
import json
import tempfile
import os
import math

# -- OSC Transport --

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
        if isinstance(arg, bool):
            type_tag += b"i"
            data += struct.pack(">i", 1 if arg else 0)
        elif isinstance(arg, int):
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
    """Read the most recent Created-track index from the Ableton Live log.
    The SuperDAW ``_cmd_track_create`` method logs ``"Created track '...' at index N"``
    each time a track is successfully added.  We parse the log file (which is
    guaranteed to be up-to-date because Ableton flushes it synchronously) rather
    than relying on the temp-file approach that often returns a stale value.
    """
    import tempfile as _tf

    # Log file location (standard Live 10.1.43 path on Windows)
    _log = os.path.join(
        os.environ.get("APPDATA", ""),
        "Ableton",
        "Live 10.1.43",
        "Preferences",
        "Log.txt",
    )
    try:
        with open(_log, "r") as f:
            # Read the last ~500 lines (more than enough for any session)
            lines = f.readlines()[-500:]
        # Search backwards for the most recent "Created track" message
        for line in reversed(lines):
            if "Created track" in line and "index" in line:
                # Extract index from: Created track '...' at index N
                return int(line.strip().split("index ")[-1])
    except Exception:
        pass
    # Fallback to old temp-file approach
    try:
        with open(
            os.path.join(_tf.gettempdir(), "superdaw_track_created.txt"), "r"
        ) as f:
            return int(f.read().strip())
    except:
        return -1


# Global counter for sequential track indices after the first creation
next_track_idx = None


def create_track(name, track_type="midi"):
    """Create a new track and return its index.
    The SuperDAW ``track/create`` OSC command writes the new track index to a
    temporary file, but that file only contains the *last* created index.  When
    creating several tracks in a row the file is not updated for each track,
    causing stale values (as observed in earlier runs).  To work around this
    we keep a local ``next_track_idx`` counter.  The first call reads the true
    index from the temp file; subsequent calls infer the next index by
    incrementing the counter.
    """
    global next_track_idx
    send_osc("/superdaw/track/create", name, track_type)
    wait(1.5)
    if next_track_idx is None:
        idx = get_track_index()
        # If the index is valid, start counting from the next slot.
        next_track_idx = idx + 1 if idx >= 0 else 0
    else:
        idx = next_track_idx
        next_track_idx += 1
    print(f"  Track '{name}' index = {idx}")

    return idx


def load_instrument(track_idx, instrument_name):
    send_osc("/superdaw/track/instrument", track_idx, instrument_name)
    wait(8)
    print(f"  Loaded '{instrument_name}' on track {track_idx}")


def load_device(track_idx, device_name):
    send_osc("/superdaw/device/load", track_idx, device_name)
    wait(2)


def write_clip_from_file(track_idx, clip_idx, notes, length_beats=None):
    """Write clip via JSON file. Notes are dicts with pitch, velocity, duration, start_beat, mute."""
    clip_file = os.path.join(
        tempfile.gettempdir(), f"sd_clip_{track_idx}_{clip_idx}.json"
    )
    with open(clip_file, "w") as f:
        json.dump(notes, f)
    send_osc("/superdaw/clip/fromfile", track_idx, clip_idx, clip_file)
    if length_beats:
        send_osc("/superdaw/clip/length", track_idx, clip_idx, length_beats)


def write_cc_automation(track_idx, clip_idx, cc_number, events):
    """Write CC automation to a clip. events is list of {cc, beat, value}."""
    tmp = os.path.join(tempfile.gettempdir(), f"cc_auto_{track_idx}_{cc_number}.json")
    with open(tmp, "w") as f:
        json.dump(events, f)
    send_osc("/superdaw/clip/cc", track_idx, clip_idx, tmp)


def write_clip_notes(track_idx, clip_idx, notes):
    """Write individual MIDI notes to a clip."""
    for n in notes:
        send_osc(
            "/superdaw/clip/note",
            track_idx,
            clip_idx,
            n["pitch"],
            n["start_beat"],
            n["duration"],
            n["velocity"],
            0,
        )  # mute=0


def set_volume(track_idx, volume):
    send_osc("/superdaw/track/volume", track_idx, volume)


def set_pan(track_idx, pan):
    send_osc("/superdaw/track/pan", track_idx, pan)


def transport(playing, bpm):
    send_osc("/superdaw/transport", 1 if playing else 0, bpm)


# -- Musical Constants --

BPM = 145
BEAT = 60.0 / BPM
BAR = 4.0  # beats per bar
TOTAL_BARS = 128
TOTAL_BEATS = TOTAL_BARS * 4.0

# Psytrance scales
A_PHRYGIAN = [0, 1, 3, 5, 7, 8, 10]  # Darker - A Phrygian
A_MINOR = [0, 2, 3, 5, 7, 8, 10]  # A Natural Minor


def midi_pitch(name):
    names = {"C": 0, "D": 2, "E": 4, "F": 5, "G": 7, "A": 9, "B": 11}
    i = 0
    base = names.get(name[i].upper(), 0)
    i += 1
    accidental = 0
    while i < len(name) and name[i] in "#b":
        if name[i] == "#":
            accidental += 1
        elif name[i] == "b":
            accidental -= 1
        i += 1
    octave_str = name[i] if i < len(name) else "4"
    octave = int(octave_str)
    return (octave + 1) * 12 + base + accidental


def scale_note(scale, key, degree, octave=4):
    octave_off = degree // len(scale)
    scale_idx = degree % len(scale)
    return key + (octave + octave_off) * 12 + scale[scale_idx]


def note(pitch, velocity, start_beat, duration):
    # Clamp pitch to valid MIDI range (0-127)
    return {
        "pitch": max(0, min(pitch, 127)),
        "velocity": min(velocity, 127),
        "start_beat": start_beat,
        "duration": duration,
        "mute": 0,
    }


# -- Pattern Generators --


def kick_pattern(total_bars=TOTAL_BARS):
    """4-on-the-floor kick drum."""
    notes = []
    for bar in range(total_bars):
        bar_start = bar * 4.0
        for beat in range(4):
            start = bar_start + beat
            vel = 127 if beat == 0 else 115
            notes.append(note(36, vel, start, 0.22))
        # Extra kick on beat 4.5 of every 8th bar for variation
        if bar % 8 == 7:
            notes.append(note(36, 105, bar_start + 4.5, 0.18))
    return notes


def hihat_pattern(total_bars=TOTAL_BARS):
    """16th note hi-hats with open hat accents."""
    notes = []
    for bar in range(total_bars):
        bar_start = bar * 4.0
        for i in range(16):
            start = bar_start + i * 0.25
            if i % 4 == 0:
                vel = 110  # Beat accent
            elif i % 2 == 0:
                vel = 85  # 8th note
            else:
                vel = 60  # Off 16th
            notes.append(note(42, vel, start, 0.08))
        # Open hi-hat on last 16th of every 4th bar
        if bar % 4 == 3:
            notes.append(note(46, 90, bar_start + 3.75, 0.20))
    return notes


def snare_pattern(total_bars=TOTAL_BARS):
    """Snare/clap on beats 2 and 4 with fills."""
    notes = []
    for bar in range(total_bars):
        bar_start = bar * 4.0
        notes.append(note(40, 110, bar_start + 1.0, 0.18))
        notes.append(note(40, 108, bar_start + 3.0, 0.18))

        # Snare fill on bars divisible by 8
        if bar % 8 == 7:
            for off in [3.5, 3.75]:
                notes.append(note(40, 80, bar_start + off, 0.10))
    return notes


def psybass_pattern(key=midi_pitch("A1"), total_bars=TOTAL_BARS):
    """Driving psytrance bassline - syncopated 16th note pattern."""
    notes = []

    # Bass patterns (1 = note, 0 = rest)
    patterns = {
        0: [1, 0, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0],  # Main
        1: [1, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0],  # Variant A
        2: [1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 0, 1, 0, 0, 0],  # Variant B
        3: [1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 0],  # Double hit
    }

    for bar in range(total_bars):
        bar_start = bar * 4.0
        pat = patterns[(bar // 8) % len(patterns)]

        # Every 32 bars, add an octave-jump bass run
        octave_jump = 0
        if bar % 32 >= 28:
            octave_jump = 12  # Jump up an octave for fill

        for i, play in enumerate(pat):
            if play:
                start = bar_start + i * 0.25
                vel = 115 if i % 4 == 0 else 90
                pitch = key + octave_jump

                # Add chromatic passing tone on some off-beats
                if i == 9 and bar % 4 == 2:
                    pitch = key + 1  # Minor 2nd walk-up

                notes.append(note(pitch, vel, start, 0.16))
    return notes


def lead_melody(total_bars=TOTAL_BARS):
    """Main lead melody - A Phrygian mode."""
    # Use pitch class (0-11) for scale_note
    key_pc = midi_pitch("A") % 12
    notes = []

    # 4 different melodic phrases
    phrases = [
        # Phrase 0 - Catchy hook (bars 0-3)
        [
            (0, 0.5),
            (3, 0.25),
            (5, 0.25),
            (7, 1.0),
            (5, 0.5),
            (3, 0.5),
            (1, 0.5),
            (0, 1.0),
            (0, 0.5),
        ],
        # Phrase 1 - Ascending (bars 4-7)
        [
            (0, 0.25),
            (1, 0.25),
            (3, 0.25),
            (5, 0.25),
            (7, 0.5),
            (8, 0.5),
            (10, 1.0),
            (8, 0.5),
            (7, 0.5),
            (5, 1.0),
        ],
        # Phrase 2 - Trill play (bars 8-11)
        [
            (5, 0.25),
            (7, 0.25),
            (5, 0.25),
            (7, 0.25),
            (8, 0.5),
            (7, 0.5),
            (5, 1.0),
            (3, 0.5),
            (0, 0.5),
        ],
        # Phrase 3 - Wide intervals (bars 12-15)
        [
            (0, 0.75),
            (7, 0.25),
            (10, 0.5),
            (8, 0.5),
            (7, 0.5),
            (5, 0.5),
            (3, 0.75),
            (0, 0.25),
            (0, 0.5),
        ],
    ]

    for bar in range(total_bars):
        bar_start = bar * 4.0
        phrase_idx = (bar // 4) % len(phrases)
        phrase = phrases[phrase_idx]
        note_start = 0.0
        for degree, dur in phrase:
            start = bar_start + note_start
            if start >= bar_start + 4.0:
                break
            pitch = scale_note(A_PHRYGIAN, key_pc, degree, 4)
            vel = 100 + abs(degree - 5) * 3
            notes.append(note(pitch, min(vel, 127), start, dur * 0.85))
            note_start += dur

    return notes


def pad_chords(total_bars=TOTAL_BARS):
    """Sustained atmospheric pad chords."""
    key = midi_pitch("A")

    # Chord progression: Am - Fm - Cm - Gm
    chord_prog = [
        [0, 3, 7],  # Am (i)
        [5, 8, 12],  # Fm (VI)
        [3, 7, 10],  # Cm (III)
        [7, 10, 14],  # Gm (VII)
    ]

    notes = []
    for bar in range(total_bars):
        bar_start = bar * 4.0
        chord = chord_prog[(bar // 2) % len(chord_prog)]
        for degree in chord:
            pitch = scale_note(A_MINOR, key, degree, 3)
            vel = 55 + (bar % 16) * 2
            notes.append(note(pitch, vel, bar_start, 3.85))
    return notes


def arp_pattern(total_bars=TOTAL_BARS):
    """Fast arpeggiated 32nd-note pattern.

    The original implementation used the full MIDI pitch returned by
    ``midi_pitch("A")`` as the key for ``scale_note`` which caused pitches
    to exceed the 0-127 MIDI range (e.g. 129).  ``scale_note`` expects a pitch
    class (0-11), so we take ``% 12`` to obtain the correct key.  We also
    lower the octave from 5 to 4 to keep the notes in a comfortable range.
    Finally we clamp any generated pitch to the valid MIDI range as a safety
    net.
    """
    # Pitch class for A (0-11) - used by scale_note
    key_pc = midi_pitch("A") % 12

    arp_sequence = [
        0,
        3,
        7,
        12,
        7,
        3,
        0,
        -3,  # Am arpeggio
        0,
        3,
        7,
        10,
        7,
        3,
        0,
        -2,
    ]  # Cm variant

    notes = []
    for bar in range(total_bars):
        bar_start = bar * 4.0
        seq = arp_sequence
        if (bar // 16) % 2 == 1:
            seq = [
                0,
                5,
                8,
                12,
                8,
                5,
                0,
                -3,  # Fm variant
                0,
                5,
                8,
                10,
                8,
                5,
                0,
                -2,
            ]  # Gm variant

        for i in range(32):
            start = bar_start + i * 0.125
            degree = seq[i % len(seq)]
            pitch = scale_note(A_MINOR, key_pc, max(0, degree), 4)
            pitch = max(0, min(pitch, 127))
            vel = 65 + (i % 8) * 8
            notes.append(note(pitch, min(vel, 127), start, 0.08))
    return notes


def fx_riser(total_bars=TOTAL_BARS):
    """FX riser - chromatic ascending for build-ups every 32 bars."""
    notes = []
    for section in range(total_bars // 32):
        for bar in range(4):
            actual_bar = section * 32 + 28 + bar
            if actual_bar >= total_bars:
                break
            bar_start = actual_bar * 4.0
            intensity = (bar + 1) / 4.0
            for i in range(16):
                start = bar_start + i * 0.25
                vel = int(50 + intensity * 77)
                pitch = 48 + bar * 6 + i
                notes.append(note(min(pitch, 120), min(vel, 127), start, 0.12))
    return notes


# -- Transition Effects --


def make_filter_sweep(start_bar, end_bar, cc=74, start_val=0.2, end_val=0.9):
    """Create a filter sweep automation curve."""
    events = []
    steps = int((end_bar - start_bar) * 16)
    for i in range(steps + 1):
        beat = (start_bar + (end_bar - start_bar) * i / steps) * 4.0
        ratio = i / steps
        # Smooth ease-in/out curve
        curved = ratio * ratio * (3 - 2 * ratio)
        val = start_val + (end_val - start_val) * curved
        events.append({"cc": cc, "beat": beat, "value": int(val * 127)})
    return events


def make_volume_ramp(start_bar, end_bar, start_vol, end_vol):
    """Create a volume automation ramp - step by step."""
    steps = max(int((end_bar - start_bar) * 4), 4)
    events = []
    for i in range(steps + 1):
        beat = (start_bar + (end_bar - start_bar) * i / steps) * 4.0
        ratio = i / steps
        vol = start_vol + (end_vol - start_vol) * ratio
        events.append({"beat": beat, "volume": vol})
    return events


def apply_filter_sweeps(tracks, drop_bars, synth_track_indices):
    """Apply filter sweeps to synth tracks at transition points."""
    print("  Applying filter sweeps...")

    # Build-up HP filter rise (bars drop_start-32 to drop_start)
    for section_idx, (drop_start, drop_end) in enumerate(drop_bars):
        build_start = drop_start - 16
        if build_start < 0:
            build_start = 0

        # High-pass rise before drop
        hp_rise = make_filter_sweep(
            build_start, drop_start, cc=74, start_val=0.1, end_val=0.7
        )
        tmp = os.path.join(tempfile.gettempdir(), f"hp_rise_{section_idx}.json")
        with open(tmp, "w") as f:
            json.dump(hp_rise, f)
        for t in synth_track_indices:
            send_osc("/superdaw/clip/cc", t, 0, tmp)

        # Drop filter sweep open
        drop_sweep = make_filter_sweep(
            drop_start, drop_end, cc=74, start_val=0.3, end_val=0.95
        )
        tmp2 = os.path.join(tempfile.gettempdir(), f"drop_sweep_{section_idx}.json")
        with open(tmp2, "w") as f:
            json.dump(drop_sweep, f)
        for t in synth_track_indices:
            send_osc("/superdaw/clip/cc", t, 0, tmp2)

        wait(0.5)


# -- Main Builder --


def build_song():
    print("=" * 70)
    print("  NEW PSYTRANCE SONG at 145 BPM")
    print("  Complete 128-bar production")
    print("=" * 70)
    print()

    # Stop transport
    transport(False, BPM)
    wait(1)

    # Set tempo
    transport(False, 145)
    wait(0.5)
    print(f"Tempo set to {BPM} BPM\n")

    # ---- Step 1: Create drum tracks ----
    print("[1/9] Creating drum tracks...")

    kick_idx = create_track("Kick", "midi")
    load_instrument(kick_idx, "Drum Rack")
    set_volume(kick_idx, 0.88)
    set_pan(kick_idx, 0.0)

    hh_idx = create_track("HiHats", "midi")
    load_instrument(hh_idx, "Drum Rack")
    set_volume(hh_idx, 0.72)
    set_pan(hh_idx, 0.25)

    snare_idx = create_track("Snare", "midi")
    load_instrument(snare_idx, "Drum Rack")
    set_volume(snare_idx, 0.78)
    set_pan(snare_idx, -0.15)

    print()

    # ---- Step 2: Write drum clips ----
    print("[2/9] Writing drum MIDI clips...")

    kick_notes = kick_pattern(TOTAL_BARS)
    write_clip_from_file(kick_idx, 0, kick_notes)
    print(f"  Kick: {len(kick_notes)} notes")

    hh_notes = hihat_pattern(TOTAL_BARS)
    write_clip_from_file(hh_idx, 0, hh_notes)
    print(f"  HiHats: {len(hh_notes)} notes")

    snare_notes = snare_pattern(TOTAL_BARS)
    write_clip_from_file(snare_idx, 0, snare_notes)
    print(f"  Snare: {len(snare_notes)} notes")

    wait(2)
    print()

    # ---- Step 3: Create synth tracks ----
    print("[3/9] Creating synth tracks with Surge XT...")

    bass_idx = create_track("Psy Bass", "midi")
    load_instrument(bass_idx, "Surge XT")
    set_volume(bass_idx, 0.82)
    set_pan(bass_idx, 0.0)

    # Ensure Lead gets a unique index (bass_idx + 1) since the temp-file
    # approach sometimes returns a stale value.
    lead_idx = create_track("Lead", "midi")
    load_instrument(lead_idx, "Surge XT")
    set_volume(lead_idx, 0.76)
    set_pan(lead_idx, -0.15)

    pad_idx = create_track("Pad", "midi")
    load_instrument(pad_idx, "Surge XT")
    set_volume(pad_idx, 0.65)
    set_pan(pad_idx, 0.35)

    arp_idx = create_track("Arp", "midi")
    load_instrument(arp_idx, "Surge XT")
    set_volume(arp_idx, 0.58)
    set_pan(arp_idx, 0.20)

    fx_idx = create_track("FX Riser", "midi")
    load_instrument(fx_idx, "Surge XT")
    set_volume(fx_idx, 0.55)
    set_pan(fx_idx, 0.0)

    synth_tracks = [bass_idx, lead_idx, pad_idx, arp_idx, fx_idx]
    print()

    # ---- Step 4: Write synth MIDI clips ----
    print("[4/9] Writing synth MIDI clips...")

    bass_notes = psybass_pattern(midi_pitch("A1"), TOTAL_BARS)
    write_clip_from_file(bass_idx, 0, bass_notes)
    print(f"  Bass: {len(bass_notes)} notes")

    lead_notes = lead_melody(TOTAL_BARS)
    write_clip_from_file(lead_idx, 0, lead_notes)
    print(f"  Lead: {len(lead_notes)} notes")

    pad_notes = pad_chords(TOTAL_BARS)
    write_clip_from_file(pad_idx, 0, pad_notes)
    print(f"  Pad: {len(pad_notes)} notes")

    arp_notes = arp_pattern(TOTAL_BARS)
    write_clip_from_file(arp_idx, 0, arp_notes)
    print(f"  Arp: {len(arp_notes)} notes")

    fx_notes = fx_riser(TOTAL_BARS)
    write_clip_from_file(fx_idx, 0, fx_notes)
    print(f"  FX Riser: {len(fx_notes)} notes")

    wait(2)
    print()

    # ---- Step 5: Create return tracks with effects ----
    print("[5/9] Setting up return tracks...")

    return_reverb_idx = create_track("Hall Reverb", "audio")
    load_device(return_reverb_idx, "Hall Reverb")
    wait(1)
    send_osc("/superdaw/device/param", return_reverb_idx, 0, 0, 0.40)

    return_delay_idx = create_track("Ping Pong Delay", "audio")
    load_device(return_delay_idx, "Ping Pong Delay")
    wait(1)
    send_osc("/superdaw/device/param", return_delay_idx, 0, 0, 0.35)
    send_osc("/superdaw/device/param", return_delay_idx, 0, 1, 0.45)

    print()

    # ---- Step 6: Set up sends ----
    print("[6/9] Setting up sends...")
    send_osc("/superdaw/mixer", kick_idx, 0.05, -1.0)  # Kick -> reverb (minimal)
    send_osc("/superdaw/mixer", snare_idx, 0.35, -1.0)  # Snare -> reverb
    send_osc("/superdaw/mixer", lead_idx, 0.25, -1.0)  # Lead -> reverb
    send_osc("/superdaw/mixer", pad_idx, 0.30, -1.0)  # Pad -> reverb
    send_osc("/superdaw/mixer", snare_idx, 0.20, 0.25)  # Snare -> delay
    send_osc("/superdaw/mixer", lead_idx, 0.15, 0.25)  # Lead -> delay
    wait(1)
    print()

    # ---- Step 7: Add filter sweeps and automation ----
    print("[7/9] Adding filter sweeps and volume automation...")

    # Define section boundaries in bars
    drop_sections = [
        (32, 48),  # Drop 1
        (80, 96),  # Drop 2
    ]

    apply_filter_sweeps(synth_tracks, drop_sections, synth_tracks)

    # Volume automation for intro/outro and section transitions
    print("  Adding volume automation...")

    # Intro fade-in (bars 0-4)
    for i in range(9):
        vol = 0.0 + 0.88 * (i / 8.0)
        send_osc("/superdaw/track/volume", kick_idx, vol)
        wait(0.2)

    # Build-up volume rise (bars 28-32 before drop 1)
    for bar in range(28, 33):
        for t in synth_tracks:
            vol = 0.40 + 0.40 * ((bar - 28) / 4.0)
            set_volume(t, vol)

    # Drop 1 volume boost (bars 32-48)
    for t in synth_tracks:
        set_volume(t, 0.85)
    set_volume(kick_idx, 0.92)
    set_volume(snare_idx, 0.82)
    wait(1)

    # Breakdown volume dip (bars 48-64)
    for bar in range(48, 65):
        ratio = (bar - 48) / 16.0
        dip = 1.0 - 0.3 * math.sin(ratio * math.pi)
        for t in synth_tracks:
            orig = {
                bass_idx: 0.82,
                lead_idx: 0.76,
                pad_idx: 0.65,
                arp_idx: 0.58,
                fx_idx: 0.55,
            }
            set_volume(t, min(0.95, orig.get(t, 0.7) * dip))

    # Reverb build during breakdown (bars 56-64)
    for i in range(9):
        wet = 0.35 + 0.35 * (i / 8.0)
        send_osc("/superdaw/device/param", return_reverb_idx, 0, 0, wet)
        wait(0.3)

    # Build-up 2 volume rise (bars 68-80)
    for bar in range(68, 81):
        for t in synth_tracks:
            vol = 0.50 + 0.35 * ((bar - 68) / 12.0)
            set_volume(t, vol)

    # Drop 2 full power (bars 80-96)
    for t in synth_tracks:
        set_volume(t, 0.90)
    set_volume(kick_idx, 0.95)
    set_volume(snare_idx, 0.85)
    wait(1)

    # Outro fade (bars 120-128)
    for i in range(9):
        ratio = i / 8.0
        vol = 0.95 - 0.95 * ratio
        send_osc("/superdaw/track/volume", 0, vol)
        wait(0.3)

    print()

    # ---- Step 8: Add mastering chain ----
    print("[8/9] Adding mastering chain to Master track...")
    load_device(0, "EQ8")
    load_device(0, "Limiter")
    send_osc("/superdaw/device/param", 0, 0, 0, 0.25)  # EQ low gain
    send_osc("/superdaw/device/param", 0, 0, 9, -0.3)  # EQ high gain
    send_osc("/superdaw/device/param", 0, 1, 0, 0.96)  # Limiter ceiling
    print()

    # ---- Step 9: Start playback ----
    print("[9/9] Starting playback at 145 BPM...")
    transport(True, 145)

    print()
    print("=" * 70)
    print("  [OK] NEW PSYTRANCE SONG COMPLETE!")
    print("=" * 70)
    print()
    print("  Track Layout:")
    print(f"    Track {kick_idx}: Kick (Drum Rack)")
    print(f"    Track {hh_idx}: HiHats (Drum Rack)")
    print(f"    Track {snare_idx}: Snare/Clap (Drum Rack)")
    print(f"    Track {bass_idx}: Psy Bass (Surge XT)")
    print(f"    Track {lead_idx}: Lead (Surge XT)")
    print(f"    Track {pad_idx}: Pad (Surge XT)")
    print(f"    Track {arp_idx}: Arp (Surge XT)")
    print(f"    Track {fx_idx}: FX Riser (Surge XT)")
    print(f"    Return {return_reverb_idx}: Hall Reverb")
    print(f"    Return {return_delay_idx}: Ping Pong Delay")
    print("    Track 0: Master (EQ8 + Limiter)")
    print()
    print("  Song Structure (128 bars @ 145 BPM):")
    print("    Bars 0-4:   Intro (drums fade in)")
    print("    Bars 4-16:  Groove established")
    print("    Bars 16-28: Elements enter")
    print("    Bars 28-32: Build-up (filter rise)")
    print("    Bars 32-48: DROP 1 (full power)")
    print("    Bars 48-64: Breakdown (pads, reverb build)")
    print("    Bars 64-80: Build-up 2 (elements return)")
    print("    Bars 80-96: DROP 2 (maximum energy)")
    print("    Bars 96-120: Groove sustain")
    print("    Bars 120-128: Outro fade")
    print()
    print("  Effects & Automation:")
    print("    - Filter sweeps (CC 74) on all synths")
    print("    - Volume automation for sections")
    print("    - Reverb/delay sends on snare, lead, pad")
    print("    - Hall Reverb wet increase during breakdown")
    print("    - EQ8 + Limiter on Master")
    print("    - 145 BPM tempo")


if __name__ == "__main__":
    build_song()
