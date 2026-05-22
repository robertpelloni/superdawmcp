# SuperDAW Client Libraries

SuperDAW provides high-level client libraries in Go and TypeScript to interact with the SuperDAW-Universal-MCP server.

## Go Client

The Go client library is located in `pkg/client/go`.

### Usage

```go
import "github.com/robertpelloni/superdaw-mcp/pkg/client/go"

func main() {
    client, err := client.NewClient("./bin/superdaw-mcp")
    if err != nil {
        panic(err)
    }
    defer client.Close()

    // Set mixer volume for track 1 to 0.8
    client.SetMixer("1", 0.8, 0.0)
}
```

## TypeScript Client

The TypeScript client library is located in `pkg/client/ts`. It is built on top of the `@modelcontextprotocol/sdk`.

### Usage

```typescript
import { SuperDAWClient } from "./pkg/client/ts/src/index";

async function run() {
    const client = new SuperDAWClient("./bin/superdaw-mcp");
    await client.connect();

    // Set mixer volume for track 1 to 0.8
    await client.setMixer("1", 0.8, 0.0);

    await client.disconnect();
}
```

## Integration Tests

You can verify the connectivity between the client libraries and the server using the integration test suite:

```bash
make integration-test
```
