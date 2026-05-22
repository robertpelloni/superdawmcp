# SuperDAW Unified Protocol Specification (v1.0.0)

This document defines the unified OSC address schema and JSON-RPC parameters for the SuperDAW-Universal-MCP system. All client adapters must adhere to these standards to ensure seamless cross-DAW orchestration.

## 1. Global Addressing Principles
- Addresses are lowercase.
- Parameters are positional unless otherwise specified.
- Boolean values are represented as `0` (False) and `1` (True) in OSC.

## 2. Transport Control
| OSC Address | Arguments | Description |
|---|---|---|
| `/superdaw/transport/play` | `(int: 0/1)` | Start (1) or Stop (0) playback |
| `/superdaw/transport/record` | `(int: 0/1)` | Toggle record mode |
| `/superdaw/transport/tempo` | `(float: bpm)` | Set transport tempo |
| `/superdaw/transport/loop` | `(int: 0/1)` | Toggle loop mode |

## 3. Mixer & Track Management
| OSC Address | Arguments | Description |
|---|---|---|
| `/superdaw/track/volume` | `(str: track_id, float: volume)` | Set track volume (0.0 to 1.0) |
| `/superdaw/track/pan` | `(str: track_id, float: pan)` | Set track panning (-1.0 to 1.0) |
| `/superdaw/track/mute` | `(str: track_id, int: 0/1)` | Set track mute state |
| `/superdaw/track/solo` | `(str: track_id, int: 0/1)` | Set track solo state |
| `/superdaw/track/create` | `(str: name, str: type)` | Create new track (type: "audio", "midi") |

## 4. MIDI Clip Operations
| OSC Address | Arguments | Description |
|---|---|---|
| `/superdaw/clip/write` | `(str: track_id, int: clip_idx, str: notes_json)` | Inject MIDI notes from JSON string |
| `/superdaw/clip/clear` | `(str: track_id, int: clip_idx)` | Clear all notes in a clip |

## 5. Plugin / VST Control
| OSC Address | Arguments | Description |
|---|---|---|
| `/superdaw/device/instantiate` | `(str: track_id, str: plugin_name)` | Instantiate a VST/AU plugin |
| `/superdaw/device/param` | `(str: track_id, str: device_id, int: idx, float: val)` | Set plugin parameter |

## 6. Bidirectional State Synchronization
DAW Adapters should broadcast state changes back to the Go server on the following addresses:
| OSC Address | Arguments | Description |
|---|---|---|
| `/superdaw/state/beat` | `(int: beat)` | Broadcast current beat |
| `/superdaw/state/tempo` | `(float: bpm)` | Broadcast current tempo |
| `/superdaw/state/track_added`| `(str: track_id, str: name)` | Broadcast when a track is created |
