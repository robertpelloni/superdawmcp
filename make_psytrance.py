#!/usr/bin/env python3
# make_psytrance.py - Build a psy-trance demo in Ableton Live 10 Standard
# Uses preset drum kits and clip launching for actual audio playback

import sys
import time
import socket
import struct
import json
import os
import tempfile
from pathlib import Path

# Add client library to path
client_path = Path(__file__).parent / "pkg" / "client" / "py"
sys.path.append(str(client_path))
from superdaw_client import SuperDAWClient


def log(msg):
    print("[>] " + msg)


def send_osc(address, *args):
    """Send a raw OSC message directly to Ableton on port 11000."""

    def pad4(data):
        data = data + b"\x00"
        extra = (4 - len(data) % 4) % 4
        return data + b"\x00" * extra

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


def read_temp(filename, timeout=3.0):
    """Read a track count/index from a temp file written by SuperDAW.py."""
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
                    val = int(f.read().strip())
                return val
        except Exception:
            pass
        time.sleep(0.2)
    # Fallback: try reading whatever is there
    try:
        with open(path, "r") as f:
            return int(f.read().strip())
    except Exception:
        return None


def main():
    # 1. Query existing track count
    send_osc("/superdaw/track/count")
    time.sleep(1.0)
    base_idx = read_temp("superdaw_track_count.txt", timeout=5.0)
    if base_idx is None:
        base_idx = 0
        log("WARNING: Could not get track count, assuming 0")
    else:
        log("Existing track count: " + str(base_idx))

    # 2. Connect to Go MCP server
    client = SuperDAWClient("./bin/superdaw-mcp")
    client.connect()
    log("Connected to SuperDAW-MCP")

    # 3. Set tempo before creating tracks
    client.transport_control(playing=False, bpm=148.0)
    log("Tempo set to 148 BPM")
    time.sleep(0.5)

    # 4. Create MIDI tracks
    track_names = ["Psy Drums", "Psy Bass", "Psy Lead", "Psy Pad"]
    track_indices = []
    for name in track_names:
        send_osc("/superdaw/track/create", name, "midi")
        idx = read_temp("superdaw_track_created.txt", timeout=3.0)
        if idx is not None:
            track_indices.append(idx)
            log("Created '" + name + "' at index " + str(idx))
        else:
            fallback = base_idx + len(track_indices)
            track_indices.append(fallback)
            log("Created '" + name + "' (fallback index " + str(fallback) + ")")
        time.sleep(0.5)

    drums_idx = track_indices[0]
    bass_idx = track_indices[1]
    lead_idx = track_indices[2]
    pad_idx = track_indices[3]

    # 5. Load instruments with presets that have actual audio
    log("Loading instruments with presets...")

    # Drums: Load a preset drum kit (has samples!)
    send_osc("/superdaw/track/instrument", drums_idx, "808 Core Kit")
    time.sleep(5)

    # Bass: Simpler (empty, but will play MIDI notes)
    send_osc("/superdaw/track/instrument", bass_idx, "Simpler")
    time.sleep(5)

    # Lead: Impulse (has some built-in sounds)
    send_osc("/superdaw/track/instrument", lead_idx, "Impulse")
    time.sleep(5)

    # Pad: Simpler
    send_osc("/superdaw/track/instrument", pad_idx, "Simpler")
    time.sleep(5)

    log("Instruments loaded")

    # 6. Write MIDI clips
    # Drums: kick + hat pattern (4 bars)
    drum_notes = []
    for bar in range(4):
        offset = bar * 4
        # Kick on every beat
        drum_notes.append(
            {"pitch": 36, "start_beat": offset, "duration": 0.25, "velocity": 120}
        )
        drum_notes.append(
            {"pitch": 36, "start_beat": offset + 1, "duration": 0.25, "velocity": 120}
        )
        drum_notes.append(
            {"pitch": 36, "start_beat": offset + 2, "duration": 0.25, "velocity": 120}
        )
        drum_notes.append(
            {"pitch": 36, "start_beat": offset + 3, "duration": 0.25, "velocity": 120}
        )
        # Hi-hat on offbeats
        for beat in range(8):
            drum_notes.append(
                {
                    "pitch": 42,
                    "start_beat": offset + 0.5 + beat * 0.5,
                    "duration": 0.125,
                    "velocity": 80 if beat % 2 == 0 else 60,
                }
            )

    # Bass: driving 16th note psytrance bass on E1
    bass_notes = []
    for i in range(64):
        bass_notes.append(
            {
                "pitch": 40,  # E1
                "start_beat": i * 0.25,
                "duration": 0.2,
                "velocity": 110 if i % 4 == 0 else 85,
            }
        )

    # Lead: arpeggiated pattern
    lead_pitches = [69, 72, 76, 80, 76, 72]
    lead_notes = []
    for i in range(16):
        for j, p in enumerate(lead_pitches):
            lead_notes.append(
                {
                    "pitch": p,
                    "start_beat": i * 1.5 + j * 0.125,
                    "duration": 0.125,
                    "velocity": 100,
                }
            )

    # Pad: sustained Am chord
    pad_notes = []
    for beat in [0, 4, 8, 12]:
        pad_notes.append(
            {"pitch": 57, "start_beat": beat, "duration": 4.0, "velocity": 70}
        )
        pad_notes.append(
            {"pitch": 60, "start_beat": beat, "duration": 4.0, "velocity": 70}
        )
        pad_notes.append(
            {"pitch": 64, "start_beat": beat, "duration": 4.0, "velocity": 70}
        )

    log("Writing clips...")
    send_osc("/superdaw/clip/write", drums_idx, 0, json.dumps(drum_notes))
    time.sleep(1)
    send_osc("/superdaw/clip/write", bass_idx, 0, json.dumps(bass_notes))
    time.sleep(1)
    send_osc("/superdaw/clip/write", lead_idx, 0, json.dumps(lead_notes))
    time.sleep(1)
    send_osc("/superdaw/clip/write", pad_idx, 0, json.dumps(pad_notes))
    time.sleep(1)
    log("Clips written")

    # 7. Set mix levels
    send_osc("/superdaw/track/volume", drums_idx, 0.85)
    send_osc("/superdaw/track/volume", bass_idx, 0.80)
    send_osc("/superdaw/track/volume", lead_idx, 0.70)
    send_osc("/superdaw/track/volume", pad_idx, 0.55)
    time.sleep(0.5)

    # 8. Launch all clips in session view
    log("Launching clips...")
    for idx in track_indices:
        send_osc("/superdaw/clip/launch", idx, 0)
        time.sleep(0.2)
    time.sleep(1)

    # 9. Start transport
    send_osc("/superdaw/transport/play", 1)
    log("Playing psytrance for 16 seconds...")

    time.sleep(16)

    # 10. Stop
    send_osc("/superdaw/transport/play", 0)
    client.transport_control(playing=False)
    client.disconnect()
    log("Done!")


if __name__ == "__main__":
    try:
        main()
    except Exception as e:
        print("[!] Error: " + str(e))
        import traceback

        traceback.print_exc()
        sys.exit(1)
