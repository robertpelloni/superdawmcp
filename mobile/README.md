# SuperDAWMobile

A lightweight, cross-platform React Native mobile client that connects to the **SuperDAW-MCP** TCP gateway (port **12002**) using the built-in JSON-RPC 2.0 protocol.

Control transport (play/stop), adjust tempo, and launch scenes — all from your phone.

---

## Prerequisites

| Tool | Version |
|------|---------|
| Node.js | >= 18 |
| npm or yarn | latest |
| React Native CLI | `npm i -g react-native-cli` (optional) |
| Xcode (iOS) | 15+ (macOS only) |
| Android Studio (Android) | Hedgehog+ |
| CocoaPods (iOS) | >= 1.13 (`gem install cocoapods`) |

---

## Getting Started

```bash
# 1. Enter the mobile directory
cd mobile

# 2. Install JavaScript dependencies
npm install

# 3. Install iOS Pods (macOS only)
cd ios && pod install && cd ..
```

---

## Running the App

### iOS (macOS)

```bash
npx react-native run-ios
```

Opens the iOS Simulator automatically. The app attempts to connect to **localhost:12002**.

### Android

```bash
npx react-native run-android
```

Requires an Android emulator or connected device. If connecting to a host machine, use `10.0.2.2:12002` instead of `127.0.0.1` (edit `TcpGateway.ts` line 36).

---

## How It Works

```
┌─────────────────┐     TCP (JSON-RPC 2.0)     ┌──────────────────────┐
│  SuperDAWMobile  │ ◄─────────────────────────► │  SuperDAW-MCP Server │
│  (React Native)  │       port 12002            │  (Go TCP Gateway)   │
└─────────────────┘                              └──────────────────────┘
                                                          │
                                                          ▼
                                                  ┌──────────────────┐
                                                  │  DAW Adapter     │
                                                  │  (Ableton, etc.) │
                                                  └──────────────────┘
```

1. **Connect** — The app opens a TCP socket to `127.0.0.1:12002`.
2. **Authenticate** — The MCP gateway accepts the connection (no authentication in development mode).
3. **Send Commands** — Each UI action (Play, Stop, Tempo change) builds a JSON-RPC 2.0 request and sends it over the socket.
4. **Receive Responses** — The gateway responds with JSON-RPC results. Pending promises resolve as responses arrive.

---

## UI Overview

| Screen | Controls |
|--------|----------|
| **Home** | Connect/Disconnect button, Play/Stop toggle, Tempo slider (80–200 BPM) |

### Future screens (easily added)

- **Scenes** — Launch Ableton/Bitwig session scenes via `superdaw_fire_scene`
- **Mixer** — Adjust individual track volume/pan via `superdaw_set_mixer`
- **Track List** — View and select tracks via `superdaw_get_tracks`
- **Transport State** — Display current playing status and BPM via `superdaw_get_transport_state`

---

## Architecture

```
mobile/
├─ App.tsx                          # Root: NavigationContainer + Stack navigator
├─ index.js                         # Entry point (registers app)
├─ package.json                     # Dependencies & scripts
├─ babel.config.js                  # Babel (metro preset)
├─ tsconfig.json                    # TypeScript config
├─ metro.config.js                  # Metro bundler config
├─ jest.config.js                   # Jest test config
├─ react-native.config.js           # RN CLI config
├─ .gitignore
├─ Gemfile                          # Ruby gems for CocoaPods
├─ README.md
├─ __tests__/
│   └─ App.test.tsx                 # Smoke test
└─ src/
    ├─ types/
    │   └─ react-native-tcp-socket.d.ts  # Type declarations for RN packages
    ├─ services/
    │   └─ TcpGateway.ts            # TCP socket singleton (JSON-RPC)
    ├─ hooks/
    │   └─ useMcp.ts                # React hook: connection state, transport, tempo
    ├─ screens/
    │   └─ HomeScreen.tsx           # Main control screen
    └─ components/
        ├─ PlayButton.tsx           # Play / Stop toggle button
        └─ TempoSlider.tsx          # BPM slider (80–200)
```

---

## JSON-RPC Protocol

All requests use JSON-RPC 2.0. Example:

```json
// Request (client → server)
{
  "jsonrpc": "2.0",
  "method": "superdaw_transport_control",
  "params": { "playing": true, "bpm": 145 },
  "id": 1
}

// Response (server → client)
{
  "jsonrpc": "2.0",
  "result": { "success": true },
  "id": 1
}
```

### Available methods used by this app

| Method | Params | Description |
|--------|--------|-------------|
| `superdaw_transport_control` | `{ playing: bool, bpm?: number }` | Start/stop playback and set tempo |
| `superdaw_get_transport_state` | `{ daw?: string }` | Get current transport state |
| `superdaw_fire_scene` | `{ scene_index: number }` | Launch a session scene |
| `superdaw_get_tracks` | `{ daw?: string }` | List all tracks in the DAW |

> Full API: see `pkg/mcp/protocol.go` in the SuperDAW-MCP root.

---

## Testing

```bash
npm test
```

Runs Jest with the `react-native` preset. Tests are in `__tests__/`.

---

## Extending the App

1. **Add a new screen** — Create `src/screens/MixerScreen.tsx`, import `TcpGateway.sendRequest`, add a route in `App.tsx`.
2. **Send custom commands** — Use `TcpGateway.sendRequest('superdaw_custom_command', { command: '...', args: {...} })`.
3. **Add authentication** — Call a custom method before other requests, and store a session token.

---

## Troubleshooting

| Problem | Fix |
|---------|-----|
| **Connection refused** | Ensure the MCP server is running: `./bin/superdaw-mcp` or the gateway process. Check port 12002 with `telnet 127.0.0.1 12002`. |
| **Android can't connect to localhost** | Android emulator loops back to itself. Use `10.0.2.2:12002` instead, or connect to the host's LAN IP. |
| **iOS Simulator can't connect** | The simulator shares the host network — `127.0.0.1` should work. |
| **App crashes on startup** | Run `npm install` and `cd ios && pod install && cd ..` again. |
| **TypeScript errors in IDE** | Run `npm install` to resolve `@types/*` declarations. |
