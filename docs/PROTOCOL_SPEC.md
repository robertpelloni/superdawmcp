# SuperDAW Unified Protocol Specification (v1.0.0)

## 1. Transport Control
| OSC Address | Arguments | Description |
|---|---|---|
| /superdaw/transport/play | (int: 0/1) | Start (1) or Stop (0) playback |
| /superdaw/transport/tempo | (float: bpm) | Set transport tempo |

## 2. Mixer & Track Management
| OSC Address | Arguments | Description |
|---|---|---|
| /superdaw/track/volume | (str: track_id, float: volume) | Set track volume (0.0 to 1.0) |
| /superdaw/track/pan | (str: track_id, float: pan) | Set track panning (-1.0 to 1.0) |
| /superdaw/track/create | (str: name, str: type) | Create new track |

## 3. MIDI Clip Operations
| OSC Address | Arguments | Description |
|---|---|---|
| /superdaw/clip/write | (str: track_id, int: clip_idx, str: notes_json) | Inject MIDI notes |
