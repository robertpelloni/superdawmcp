#!/usr/bin/env python3
"""
Automate Surge XT parameters via MIDI CC for the SuperDAW psytrance production.

This writes MIDI CC automation envelopes to the existing clips.
Common Surge XT MIDI CC mappings:
  CC 74 = Filter Cutoff (0-127)
  CC 71 = Filter Resonance (0-127)
  CC 73 = Amp Attack (0-127)
  CC 75 = Amp Decay (0-127)
  CC 76 = Amp Sustain (0-127)
  CC 77 = Amp Release (0-127)
  CC 1  = Mod Wheel (0-127)
  CC 11 = Expression (0-127)
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


def write_cc_batch(track_idx, clip_idx, cc_events):
    """Write batch CC automation via JSON file.
    cc_events: list of {"cc": int, "value": int, "beat": float}
    """
    clip_file = os.path.join(
        tempfile.gettempdir(), "superdaw_cc_%d_%d.json" % (track_idx, clip_idx)
    )
    with open(clip_file, "w") as f:
        json.dump(cc_events, f)
    send_osc("/superdaw/clip/cc", track_idx, clip_idx, clip_file)
    wait(2)


def generate_cc_ramp(cc_num, start_val, end_val, start_beat, end_beat, steps=16):
    """Generate a CC ramp from start_val to end_val over the given beat range."""
    events = []
    for i in range(steps + 1):
        beat = start_beat + (end_beat - start_beat) * i / steps
        val = int(start_val + (end_val - start_val) * i / steps)
        val = max(0, min(127, val))
        events.append({"cc": cc_num, "value": val, "beat": round(beat, 2)})
    return events


def generate_cc_lfo(cc_num, min_val, max_val, start_beat, end_beat, rate_bars=1):
    """Generate CC LFO (sine wave) over the given beat range."""
    import math

    events = []
    beat = start_beat
    while beat <= end_beat:
        phase = (beat - start_beat) / (rate_bars * 4)  # 4 beats per bar
        val = int(
            (min_val + max_val) / 2
            + (max_val - min_val) / 2 * math.sin(2 * math.pi * phase)
        )
        val = max(0, min(127, val))
        events.append({"cc": cc_num, "value": val, "beat": round(beat, 2)})
        beat += 0.5  # Every half beat
    return events


def generate_cc_steps(cc_num, values, start_beat, beat_length=4.0):
    """Generate CC step sequence (one value per bar/beat_length)."""
    events = []
    for i, val in enumerate(values):
        beat = start_beat + i * beat_length
        events.append({"cc": cc_num, "value": int(val), "beat": round(beat, 2)})
    return events


def automate_bass(track_idx=8):
    """Automate Surge XT bass with filter sweeps."""
    print("  Automating bass filter (CC 74)...")

    events = []

    # Bass filter automation: low cutoff in intro, opens up in drops
    # Bars 0-8: Very closed filter (dark, building tension)
    events += generate_cc_ramp(74, 15, 25, 0, 32, steps=8)

    # Bars 8-16: Opens slightly
    events += generate_cc_ramp(74, 25, 45, 32, 64, steps=8)

    # Bars 16-32: Full open (main drop)
    events += generate_cc_ramp(74, 45, 80, 64, 128, steps=16)

    # Bars 32-64: Filter sweep (classic psytrance buildup)
    events += generate_cc_ramp(74, 80, 100, 128, 256, steps=32)

    # Bars 64-96: Filter opens and closes rhythmically
    events += generate_cc_lfo(74, 40, 100, 256, 384, rate_bars=0.5)

    # Bars 96-128: Final section - wide filter sweeps
    events += generate_cc_lfo(74, 30, 110, 384, 512, rate_bars=1)

    write_cc_batch(track_idx, 0, events)
    print(f"    Wrote {len(events)} CC events")


def automate_lead(track_idx=9):
    """Automate Surge XT lead with expression and filter."""
    print("  Automating lead expression (CC 1, 74)...")

    events = []

    # Lead comes in at bar 16 (beat 64)
    # Expression fade-in
    events += generate_cc_ramp(11, 0, 100, 64, 80, steps=8)

    # Filter opens with the melody
    events += generate_cc_ramp(74, 60, 90, 64, 128, steps=16)

    # Vibrato (mod wheel) increases in intensity
    events += generate_cc_ramp(1, 0, 40, 64, 192, steps=16)

    # Filter sweep in the breakdown
    events += generate_cc_ramp(74, 90, 50, 192, 256, steps=16)

    # Filter opens again in the drop
    events += generate_cc_ramp(74, 50, 100, 256, 320, steps=16)

    # Expression swells in the climax
    events += generate_cc_ramp(11, 80, 127, 320, 384, steps=16)

    # Final section - filter and expression automation
    events += generate_cc_lfo(74, 50, 100, 384, 512, rate_bars=2)
    events += generate_cc_ramp(11, 100, 60, 384, 512, steps=16)

    write_cc_batch(track_idx, 0, events)
    print(f"    Wrote {len(events)} CC events")


def automate_pad(track_idx=10):
    """Automate Surge XT pad with filter and expression."""
    print("  Automating pad filter and expression (CC 74, 11)...")

    events = []

    # Pad filter: slowly opens throughout the track
    events += generate_cc_ramp(74, 30, 70, 0, 256, steps=32)

    # Expression swells
    events += generate_cc_ramp(11, 40, 80, 0, 128, steps=16)
    events += generate_cc_ramp(11, 80, 50, 128, 256, steps=16)
    events += generate_cc_ramp(11, 50, 90, 256, 384, steps=16)
    events += generate_cc_ramp(11, 90, 60, 384, 512, steps=16)

    # Filter opens more in second half
    events += generate_cc_ramp(74, 70, 100, 256, 512, steps=32)

    write_cc_batch(track_idx, 0, events)
    print(f"    Wrote {len(events)} CC events")


def automate_arp(track_idx=11):
    """Automate Surge XT arp with filter and resonance."""
    print("  Automating arp filter and resonance (CC 74, 71)...")

    events = []

    # Arp filter: rhythmic filter sweeps (classic psytrance)
    # Use LFO pattern for filter cutoff
    events += generate_cc_lfo(74, 30, 90, 0, 128, rate_bars=0.25)
    events += generate_cc_lfo(74, 40, 100, 128, 256, rate_bars=0.5)
    events += generate_cc_lfo(74, 50, 110, 256, 384, rate_bars=0.25)
    events += generate_cc_lfo(74, 60, 120, 384, 512, rate_bars=0.125)

    # Resonance increases with the filter
    events += generate_cc_ramp(71, 40, 80, 0, 256, steps=32)
    events += generate_cc_ramp(71, 80, 60, 256, 512, steps=32)

    write_cc_batch(track_idx, 0, events)
    print(f"    Wrote {len(events)} CC events")


def main():
    print("=" * 60)
    print("SURGE XT MIDI CC AUTOMATION")
    print("=" * 60)
    print()
    print("Writing MIDI CC automation envelopes to clips...")
    print("This controls Surge XT's filter, resonance, and expression.")
    print()

    automate_bass(8)
    automate_lead(9)
    automate_pad(10)
    automate_arp(11)

    print()
    print("=" * 60)
    print("Automation complete! Check the clips in Ableton.")
    print("CC 74 = Filter Cutoff, CC 71 = Resonance, CC 11 = Expression")
    print("=" * 60)


if __name__ == "__main__":
    main()
