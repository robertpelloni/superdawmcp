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

## High-Level CLI Adapters
Located in `scripts/`. These provide a convenient way to control specific DAWs from the terminal.

### Ableton Adapter
```bash
./scripts/ableton_adapter.py transport play --bpm 120
./scripts/ableton_adapter.py mixer 1 --vol 0.8 --pan -0.2
./scripts/ableton_adapter.py create-track "Acid Bass" --type midi
./scripts/ableton_adapter.py gen-euclidean 1 --hits 5 --steps 16 --pitch 36
```

### REAPER Adapter
```bash
./scripts/reaper_adapter.py transport play --bpm 140
./scripts/reaper_adapter.py gen-euclidean 2 --hits 3 --steps 8 --pitch 42
./scripts/reaper_adapter.py stems input.wav ./output --count 5
```
