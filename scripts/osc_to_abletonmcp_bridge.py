"""
OSC-to-AbletonMCP Bridge
Listens for SuperDAW OSC commands on UDP 11000 and forwards
them as JSON/TCP commands to AbletonMCP on port 9877.

Fixes:
- Creates clip before writing notes (AbletonMCP requires existing clip)
- Translates start_beat -> start_time in note data
- Properly sets track name after creation
"""

import socket
import json
import struct

ABLETON_MCP_HOST = "127.0.0.1"
ABLETON_MCP_PORT = 9877
OSC_LISTEN_PORT = 11000


# ---------- minimal OSC decode ----------
def decode_osc(data):
    """Yield (address, args) tuples from an OSC packet."""
    if data.startswith(b"#bundle"):
        pos = 8
        pos += 8
        while pos < len(data):
            size = struct.unpack(">I", data[pos : pos + 4])[0]
            pos += 4
            for msg in decode_osc(data[pos : pos + size]):
                yield msg
            pos += size
    else:
        addr_end = data.find(b"\x00")
        if addr_end < 0:
            return
        address = data[:addr_end].decode("utf-8", errors="replace")
        pos = addr_end + 1
        pos = ((pos + 3) // 4) * 4
        if pos >= len(data) or data[pos : pos + 1] != b",":
            yield (address, [])
            return
        tag_end = data.find(b"\x00", pos)
        tags = data[pos + 1 : tag_end]
        pos = tag_end + 1
        pos = ((pos + 3) // 4) * 4
        args = []
        for t in tags:
            if t == ord("i"):
                args.append(struct.unpack(">i", data[pos : pos + 4])[0])
                pos += 4
            elif t == ord("f"):
                args.append(struct.unpack(">f", data[pos : pos + 4])[0])
                pos += 4
            elif t == ord("s"):
                s_end = data.find(b"\x00", pos)
                args.append(data[pos:s_end].decode("utf-8", errors="replace"))
                pos = ((s_end + 4) // 4) * 4
            elif t == ord("b"):
                size = struct.unpack(">I", data[pos : pos + 4])[0]
                pos += 4
                args.append(data[pos : pos + size])
                pos = ((pos + size + 3) // 4) * 4
            elif t == ord("T"):
                args.append(True)
            elif t == ord("F"):
                args.append(False)
            else:
                args.append(None)
        yield (address, args)


# ---------- Bridge ----------
_last_track_index = 0  # track what index the last created track got


def send_ableton_cmd(command_type, **params):
    """Send a JSON command to the AbletonMCP server and return parsed response."""
    global _last_track_index
    payload = json.dumps(dict(type=command_type, **params))
    try:
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.settimeout(3.0)
        s.connect((ABLETON_MCP_HOST, ABLETON_MCP_PORT))
        s.sendall((payload + "\n").encode("utf-8"))
        buf = b""
        while True:
            try:
                chunk = s.recv(4096)
            except socket.timeout:
                break
            if not chunk:
                break
            buf += chunk
            if b"\n" in buf:
                break
        s.close()
        if buf:
            resp = json.loads(buf.decode("utf-8").split("\n")[0])
            # Track the index of newly created tracks
            if command_type == "create_midi_track":
                result = resp.get("result", {})
                if isinstance(result, dict) and "index" in result:
                    _last_track_index = int(result["index"])
            return resp
    except Exception as e:
        print(f"[Bridge] Error sending to AbletonMCP: {e}")
    return None


def handle_osc(addr, args):
    """Translate an OSC command to an AbletonMCP JSON command."""
    print(f"[Bridge] OSC {addr} {args}")

    if addr == "/superdaw/transport/play":
        if args and args[0]:
            resp = send_ableton_cmd("start_playback")
        else:
            resp = send_ableton_cmd("stop_playback")

    elif addr == "/superdaw/transport/tempo":
        if args:
            try:
                bpm = float(args[0])
                resp = send_ableton_cmd("set_tempo", tempo=bpm)
            except Exception:
                # Skip malformed tempo messages
                return

    elif addr == "/superdaw/track/create":
        name = str(args[0]) if len(args) > 0 else "Track"
        # Create track (index=-1 appends)
        resp = send_ableton_cmd("create_midi_track", index=-1)
        if not resp or resp.get("status") != "success":
            # creation failed
            return
        # Extract the new track index from the response
        new_index = int(resp.get("result", {}).get("index", 0))
        # Set the track name
        send_ableton_cmd("set_track_name", track_index=new_index, name=name)
        # Load a default instrument (Analog) so the track is playable
        send_ableton_cmd(
            "add_device_to_track",
            track_index=new_index,
            device_name="Analog",
            device_type="instrument",
        )
        # Create an initial empty clip (4 beats) on the first clip slot
        send_ableton_cmd("create_clip", track_index=new_index, clip_index=0, length=4.0)
        print(
            f"[Bridge] Created track {new_index} named '{name}' with instrument and empty clip."
        )

    elif addr == "/superdaw/track/volume":
        if len(args) >= 2:
            resp = send_ableton_cmd(
                "set_track_volume", track_index=int(args[0]), volume=float(args[1])
            )

    elif addr == "/superdaw/track/pan":
        if len(args) >= 2:
            resp = send_ableton_cmd(
                "set_track_pan", track_index=int(args[0]), pan=float(args[1])
            )

    elif addr == "/superdaw/clip/write":
        if len(args) >= 3:
            track_index = int(args[0])
            clip_index = int(args[1])

            # Parse notes, translate start_beat -> start_time
            raw_notes = (
                json.loads(args[2])
                if isinstance(args[2], str)
                else json.loads(args[2].decode("utf-8"))
            )
            notes = []
            for n in raw_notes:
                notes.append(
                    {
                        "pitch": n.get("pitch", 60),
                        "start_time": n.get("start_beat", n.get("start_time", 0.0)),
                        "duration": n.get("duration", 0.25),
                        "velocity": n.get("velocity", 100),
                        "mute": n.get("mute", False),
                    }
                )

            # 1) Create an empty clip first (4 beats long)
            # If clip already exists, that's fine - just write notes
            resp1 = send_ableton_cmd(
                "create_clip",
                track_index=track_index,
                clip_index=clip_index,
                length=4.0,
            )
            clip_ok = resp1 and resp1.get("status") == "success"
            clip_exists = (
                resp1
                and resp1.get("status") == "error"
                and "already has a clip" in resp1.get("message", "")
            )
            if clip_ok or clip_exists:
                # 2) Write notes into the clip
                resp2 = send_ableton_cmd(
                    "add_notes_to_clip",
                    track_index=track_index,
                    clip_index=clip_index,
                    notes=notes,
                )
                if resp2 and resp2.get("status") == "success":
                    print(
                        f"[Bridge] Wrote {len(notes)} notes to track {track_index}, clip {clip_index}"
                    )
                else:
                    print(f"[Bridge] Failed to write notes: {resp2}")
            else:
                print(f"[Bridge] Failed to create clip: {resp1}")

    elif addr == "/superdaw/clip/delete":
        if len(args) >= 2:
            resp = send_ableton_cmd(
                "stop_clip", track_index=int(args[0]), clip_index=int(args[1])
            )

    elif addr == "/live/scene/fire":
        if args:
            resp = send_ableton_cmd("fire_clip", track_index=int(args[0]), clip_index=0)

    else:
        print(f"[Bridge] Unhandled OSC address: {addr}")


def bridge_loop():
    """Listen for OSC packets on UDP port and bridge them."""
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    sock.bind(("127.0.0.1", OSC_LISTEN_PORT))
    print(f"[Bridge] Listening for SuperDAW OSC on port {OSC_LISTEN_PORT}")
    print(f"[Bridge] Forwarding to AbletonMCP on port {ABLETON_MCP_PORT}")

    while True:
        try:
            data, addr = sock.recvfrom(65536)
            for osc_addr, osc_args in decode_osc(data):
                handle_osc(osc_addr, osc_args)
        except KeyboardInterrupt:
            break
        except Exception as e:
            print(f"[Bridge] Error: {e}")


if __name__ == "__main__":
    bridge_loop()
