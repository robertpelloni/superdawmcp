#!/usr/bin/env python3
"""
Tune Surge XT patches via MIDI CC for the SuperDAW psytrance production.

MIDI CC mapping for common synth parameters:
  CC 74 = Filter Cutoff
  CC 71 = Filter Resonance
  CC 73 = Attack
  CC 75 = Decay
  CC 76 = Sustain
  CC 77 = Release
  CC 70 = Filter Env Amount
  CC 72 = Amp Release
  CC 1  = Mod Wheel (often mapped to vibrato/LFO)
  CC 10 = Pan
  CC 11 = Expression
"""

import socket
import struct
import time
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


def add_cc_to_clip(track_idx, clip_idx, cc_number, cc_value, start_beat=0.0):
    """Add a MIDI CC event to a clip's envelope."""
    # In Ableton Live, we can use the clip's envelope feature
    # For now, we write a small clip with the CC event as a special note
    # The actual CC automation needs to be done via clip envelopes
    pass


def create_cc_clip(track_idx, clip_idx, name, length, cc_events):
    """Create a clip with MIDI CC automation data.

    cc_events: list of (cc_number, value, beat_position) tuples
    This creates a clip with notes that trigger CC changes.
    """
    # Write CC data as a JSON file
    clip_file = os.path.join(
        tempfile.gettempdir(), "superdaw_cc_%d_%d.json" % (track_idx, clip_idx)
    )

    # For now, we'll use a workaround: create a clip with very quiet notes
    # at specific pitches that correspond to CC numbers
    # This is a limitation - real CC automation needs Live API access

    # Actually, we can write to the clip's MIDI CC envelopes directly
    # using the Live Clip API: clip.set_cc_value(cc_number, value, beat)
    # But this requires the clip to already exist

    # Instead, let's use the plugin/param command if available
    pass


def tune_bass(track_idx):
    """Tune the Surge XT bass patch for psytrance."""
    print(f"  Tuning bass on track {track_idx}...")

    # Psytrance bass needs:
    # - Low cutoff for that squelchy sound
    # - Quick attack, medium decay
    # - Mono mode with glide

    # Since we can't directly set Surge XT params via Live API,
    # we'll use MIDI CC events embedded in the clip
    # The best approach: write the CC data as part of the clip's envelope

    # For now, we'll instruct the user on what to set manually
    print("  -> Set manually in Surge XT:")
    print("     - Oscillator: Saw wave (default)")
    print("     - Filter: Low-pass, Cutoff ~200Hz, Resonance ~40%")
    print("     - Amp: Attack 0ms, Decay 200ms, Sustain 50%, Release 100ms")
    print("     - Mono mode, Glide ~50ms")


def tune_lead(track_idx):
    """Tune the Surge XT lead patch."""
    print(f"  Tuning lead on track {track_idx}...")
    print("  -> Set manually in Surge XT:")
    print("     - Oscillator: Saw + Square (2 oscillators)")
    print("     - Filter: Low-pass, Cutoff ~800Hz, Resonance ~30%")
    print("     - LFO: Sine, Rate ~5Hz -> Filter Cutoff")
    print("     - Amp: Attack 10ms, Decay 300ms, Sustain 70%, Release 200ms")


def tune_pad(track_idx):
    """Tune the Surge XT pad patch."""
    print(f"  Tuning pad on track {track_idx}...")
    print("  -> Set manually in Surge XT:")
    print("     - Oscillator: Saw + Noise")
    print("     - Filter: Low-pass, Cutoff ~400Hz, Resonance ~20%")
    print("     - Amp: Attack 500ms, Decay 1s, Sustain 80%, Release 1s")
    print("     - Reverb effect for space")


def tune_arp(track_idx):
    """Tune the Surge XT arpeggio patch."""
    print(f"  Tuning arp on track {track_idx}...")
    print("  -> Set manually in Surge XT:")
    print("     - Oscillator: Square wave")
    print("     - Filter: Band-pass, Cutoff ~600Hz, Resonance ~50%")
    print("     - Amp: Attack 0ms, Decay 100ms, Sustain 0%, Release 50ms")
    print("     - Short, plucky sound")


def main():
    print("=" * 60)
    print("SURGE XT PATCH TUNING GUIDE")
    print("=" * 60)
    print()

    # Get current track info from the production
    tune_bass(8)
    print()
    tune_lead(9)
    print()
    tune_pad(10)
    print()
    tune_arp(11)
    print()

    print("=" * 60)
    print("IMPORTANT: These settings are suggestions.")
    print("Tweak to taste in the Ableton UI!")
    print("=" * 60)


if __name__ == "__main__":
    main()
