# SuperDAW Client Libraries

## Python Client
Located in `pkg/client/py`.
Usage:
```python
from superdaw_client.client import SuperDAWClient
client = SuperDAWClient("./bin/superdaw-mcp")
client.connect()
client.transport_control(playing=True, bpm=128.0)
```

## Go Client
Located in `pkg/client/go`.
Usage:
```go
client, _ := client.NewClient("./bin/superdaw-mcp")
client.TransportControl(true, 128.0)
```

## TypeScript Client
Located in `pkg/client/ts`.
Usage:
```typescript
const client = new SuperDAWClient("./bin/superdaw-mcp");
await client.connect();
await client.transportControl(true, 128.0);
```

## C# Client
Located in `pkg/client/csharp`. Ideal for Unity and Godot.
Usage:
```csharp
using SuperDAW.Client;
var client = new SuperDAWClient("./bin/superdaw-mcp");
client.Connect();
await client.TransportControlAsync(true, 128.0);
```

## Ruby Client
Located in `pkg/client/rb`. Ideal for Sonic Pi integration.
Usage:
```ruby
require './pkg/client/rb/superdaw_client'
client = SuperDAW::Client.new("./bin/superdaw-mcp")
client.connect
client.transport_control(true, bpm: 128.0)
```

## Rust Client
Located in `pkg/client/rust`.
Usage:
```rust
let mut client = SuperDAWClient::new("./bin/superdaw-mcp")?;
client.transport_control(true, Some(128.0), None)?;
```

## High-Level CLI Adapters
Located in `scripts/`. These provide a convenient way to control specific DAWs from the terminal.

### Ableton Adapter
```bash
./scripts/ableton_adapter.py transport play --bpm 120
./scripts/ableton_adapter.py mixer 1 --vol 0.8 --pan -0.2
./scripts/ableton_adapter.py create-track "Acid Bass" --type midi
```

### REAPER Adapter
```bash
./scripts/reaper_adapter.py transport play --bpm 140
./scripts/reaper_adapter.py gen-euclidean 2 --hits 3 --steps 8 --pitch 42
```

### Bitwig Adapter
```bash
./scripts/bitwig_adapter.py transport play
./scripts/bitwig_adapter.py mixer 0 --vol 0.9
```

### FL Studio Adapter
```bash
./scripts/flstudio_adapter.py transport play
./scripts/flstudio_adapter.py create-track "Kick" --type audio
```

## Connectors (Real-time Integration)
Specialized scripts in `examples/` for external hardware and multi-DAW workflows.

### MIDI-to-MCP Bridge
Translates physical MIDI controller input into MCP commands.
```bash
pip install mido python-rtmidi
python3 examples/midi_bridge/midi_to_mcp.py --port "LPD8" --daw ableton
```

### Multi-DAW Sync
Keeps multiple DAWs in sync.
```bash
python3 examples/multi_daw_jam/sync_connector.py --daws ableton,reaper,bitwig
```
