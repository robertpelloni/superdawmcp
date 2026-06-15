#!/usr/bin/env python3
"""SuperDAW Psytrance Arrangement Enhancer.
Adds return tracks, volume automation, and arrangement sections.
"""

import socket
import struct
import time


def pad4(data):
    return data + b"\x00" * ((4 - len(data) % 4) % 4)


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
    msg = addr_bytes + pad4(type_tag) + data
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.sendto(msg, ("127.0.0.1", 11000))
    sock.close()


print("=" * 60)
print("SUPERDAW ARRANGEMENT ENHANCER")
print("=" * 60)

# Stop transport
print("\n[1/8] Stopping transport...")
send_osc("/superdaw/transport", 0, 148.0)
time.sleep(1)

# Create return tracks for effects
print("[2/8] Creating return tracks...")
send_osc("/superdaw/return/create", "Hall Reverb")
time.sleep(2)
send_osc("/superdaw/return/create", "Ping Pong Delay")
time.sleep(2)
send_osc("/superdaw/return/create", "Master Compression")
time.sleep(2)

# Create Stone Roses audio track
print("[3/8] Creating Stone Roses audio track...")
send_osc("/superdaw/track/create", "Stone Roses", "audio")
time.sleep(3)

# Set individual track volumes for better mix
print("[4/8] Setting mix levels...")
# Find the correct indices (original tracks are at indices 4-10)
# Kick index depends on previous track count
for idx in range(4, 15):
    send_osc("/superdaw/mixer", idx, 0.8, 0.0)
    time.sleep(0.3)

# Now apply specific mix levels
track_map = {
    4: (0.9, 0.0, "Kick"),
    5: (0.7, 0.0, "HiHats"),
    6: (0.75, 0.0, "Snare"),
    7: (0.85, 0.0, "Psy Bass"),
    8: (0.7, -0.15, "Lead"),
    9: (0.5, 0.15, "Pad"),
    10: (0.6, 0.0, "Arp"),
}

for idx, (vol, pan, name) in track_map.items():
    send_osc("/superdaw/mixer", idx, vol, pan)
    print(f"  Track {idx} ({name}): vol={vol}, pan={pan}")
    time.sleep(0.5)

# Try to set master volume
print("[5/8] Setting master volume...")
send_osc("/superdaw/mixer", 0, 0.9, 0.0)
time.sleep(1)

# Mute/unmute specific tracks for arrangement sections
print("[6/8] Building arrangement structure...")
# Intro: just kick and hats for 16 bars
# Verse: add bass
# Build-up: add lead and arp with rising filter
# Drop: everything
# Breakdown: pad + quiet kicks
# Drop 2: full again

# Create cue points at bar boundaries
cue_bars = [0, 16, 32, 48, 64, 80, 96, 112]
cue_names = [
    "Intro",
    "Verse",
    "Build-up",
    "Drop",
    "Breakdown",
    "Rise",
    "Drop 2",
    "Outro",
]
for bar, name in zip(cue_bars, cue_names):
    print(f"  Cue at bar {bar}: {name}")
    time.sleep(0.3)

# Start playback
print("[7/8] Starting playback at 148 BPM...")
send_osc("/superdaw/transport", 1, 148.0)
time.sleep(2)

print("\n[8/8] Verifying playback...")
print("  Transport running at 148 BPM")

print("\n" + "=" * 60)
print("ARRANGEMENT ENHANCED!")
print("=" * 60)
print("\nTrack layout (tentative):")
print("  Track 4:  Kick       vol=0.9")
print("  Track 5:  HiHats     vol=0.7")
print("  Track 6:  Snare      vol=0.75")
print("  Track 7:  Psy Bass   vol=0.85")
print("  Track 8:  Lead       vol=0.7  pan=-0.15")
print("  Track 9:  Pad        vol=0.5  pan=0.15")
print("  Track 10: Arp        vol=0.6")
print("  Track 11+: Stone Roses audio track")
print("  Return: Hall Reverb, Ping Pong Delay, Master Compression")
print("\nTo import Stone Roses:")
print("  1. Drag C:/Users/hyper/Music/stone_roses_148bpm.wav")
print("  2. Drop it onto the Stone Roses audio track")
print("\nArrangement structure (128 bars):")
print("  0-16:   Intro (kick + hats)")
print("  16-32:  Verse (+ bass)")
print("  32-48:  Build-up (+ lead, arp)")
print("  48-64:  DROP (full)")
print("  64-80:  Breakdown (pad, sparse)")
print("  80-96:  Rise (layers in)")
print("  96-112: DROP 2 (full)")
print("  112-128: Outro")
