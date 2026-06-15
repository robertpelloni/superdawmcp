#!/usr/bin/env python3
"""SuperDAW – add transitions, effects, mastering chain and render the mix."""

import socket
import time
import json
import tempfile
import os


def pad4(b):
    return b + b"\x00" * ((4 - len(b) % 4) % 4)


def send_osc(address, *args):
    a = pad4(address.encode("utf-8"))
    type_tag = b","
    data = b""
    for arg_val in args:
        type_tag += b"s"
        data += pad4(str(arg_val).encode("utf-8"))
    msg = a + pad4(type_tag) + data
    socket.socket(socket.AF_INET, socket.SOCK_DGRAM).sendto(msg, ("127.0.0.1", 11000))


def write_cc_envelope(track_idx, start_beat, end_beat, cc, start_val, end_val):
    env = [
        {"cc": cc, "beat": start_beat, "value": start_val},
        {"cc": cc, "beat": end_beat, "value": end_val},
    ]
    tmp = os.path.join(tempfile.gettempdir(), f"cc_env_{track_idx}_{cc}.json")
    with open(tmp, "w") as f:
        json.dump(env, f)
    send_osc("/superdaw/clip/cc", track_idx, 0, tmp)


def fade_master_volume(start_vol, end_vol, steps=8, step_delay=0.25):
    print(f"Master-volume fade {start_vol:.2f} -> {end_vol:.2f}")
    for i in range(steps + 1):
        vol = start_vol + (end_vol - start_vol) * i / steps
        send_osc("/superdaw/track/volume", 0, vol)
        time.sleep(step_delay)


def ramp_return_param(
    return_idx, dev_idx, param_idx, start_val, end_val, steps=6, step_delay=0.3
):
    print(f"Ramping return {return_idx} dev {dev_idx} param {param_idx}")
    for i in range(steps + 1):
        val = start_val + (end_val - start_val) * i / steps
        send_osc("/superdaw/device/param", return_idx, dev_idx, param_idx, val)
        time.sleep(step_delay)


def set_master_eq_and_limiter():
    print("[Master Bus] Loading EQ8")
    send_osc("/superdaw/device/load", 0, "EQ8")
    time.sleep(2)
    print("[Master Bus] Loading Limiter")
    send_osc("/superdaw/device/load", 0, "Limiter")
    time.sleep(2)
    print("[EQ] Setting initial EQ values")
    send_osc("/superdaw/device/param", 0, 0, 0, 0.3)
    send_osc("/superdaw/device/param", 0, 0, 9, -0.5)
    time.sleep(1)
    print("[Limiter] Setting ceiling to -0.2 dB")
    send_osc("/superdaw/device/param", 0, 1, 0, 0.98)
    time.sleep(1)


def render_final_mix():
    out_path = "C:/Users/hyper/Music/psytrance_final_148bpm.wav"
    sr = 48000
    print(f"[Render] Export -> {out_path} @ {sr}Hz")
    send_osc("/superdaw/render", out_path, sr)


def main():
    print("=" * 60)
    print("ADDING TRANSITIONS AND EFFECTS")
    print("=" * 60)

    print("\n[1] Ensuring transport at 148 BPM...")
    send_osc("/superdaw/transport", 1, 148.0)
    time.sleep(1)

    print("\n[2] Master-volume fade-in for drop (bars 48-64)...")
    fade_master_volume(0.80, 0.95, steps=8, step_delay=0.25)

    print("\n[3] Hall Reverb dry/wet ramp...")
    ramp_return_param(12, 0, 0, 0.30, 0.70, steps=6, step_delay=0.3)

    print("\n[4] Ping-Pong Delay feedback ramp...")
    ramp_return_param(13, 0, 1, 0.20, 0.60, steps=8, step_delay=0.3)

    print("\n[5] Filter-cutoff sweeps on all synths (CC 74)...")
    synth_tracks = [7, 8, 9, 10, 12]
    start_drop = 48 * 4
    end_drop = 64 * 4
    for t in synth_tracks:
        write_cc_envelope(t, start_drop, end_drop, 74, 0.2, 0.9)

    print("\n[6] Inverse filter sweep for breakdown (bars 64-80)...")
    start_break = 64 * 4
    end_break = 80 * 4
    for t in synth_tracks:
        write_cc_envelope(t, start_break, end_break, 74, 0.9, 0.2)

    print("\n[7] Adding mastering chain (EQ8 + Limiter)...")
    set_master_eq_and_limiter()

    print("\n[8] Master-volume fade-out for outro...")
    fade_master_volume(0.95, 0.70, steps=8, step_delay=0.25)

    print("\n[9] Sending render command...")
    render_final_mix()

    print("\n" + "=" * 60)
    print("ALL TRANSITIONS AND EFFECTS APPLIED!")
    print("=" * 60)
    print("\nWhat happened:")
    print("  - Master volume fades in at drop, fades out at outro")
    print("  - Hall Reverb dry/wet ramps before the drop")
    print("  - Ping-Pong Delay feedback ramps during the rise")
    print(
        "  - CC 74 filter sweeps on all synths (low -> high on drop, high -> low on breakdown)"
    )
    print("  - EQ8 + Limiter added to Master track")
    print("  - Render queued to C:/Users/hyper/Music/psytrance_final_148bpm.wav")
    print("\nManual step: Drag stone_roses_148bpm.wav onto Track 11 (Stone Roses)")


if __name__ == "__main__":
    main()
