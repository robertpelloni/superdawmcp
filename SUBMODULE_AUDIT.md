# SuperDAW Submodule Audit & Implementation Status

This document tracks the integration progress of architectural reference submodules into the SuperDAW-MCP core.

## 1. Ableton Live Integration
| Submodule | Key Functionality | Status | Integration Target |
|-----------|-------------------|--------|--------------------|
| `third_party/pylive` | OSC framework for Live | Integrated | `pkg/daw/ableton.go` |
| `third_party/ableton-osc` | OSC Address Schema | Integrated | `pkg/daw/ableton.go` |
| `third_party/ableton-mcp-extended` | MCP Tool Mapping | Partially Integrated | `pkg/mcp/protocol.go` |
| `third_party/launchpad95` | Hardware surface logic | Referenced | `scripts/launchpad_mcp.py` |
| `third_party/ableton-live-tools` | CLI Utilities | Redundant | `scripts/ableton_adapter.py` |

## 2. REAPER Integration
| Submodule | Key Functionality | Status | Integration Target |
|-----------|-------------------|--------|--------------------|
| `third_party/reaper-reapy-mcp` | Python API Bridge | Integrated | `pkg/agents/reaper/` |
| `third_party/total-reaper-mcp` | Comprehensive MCP schema | Partially Integrated | `pkg/mcp/protocol.go` |
| `third_party/reaper-daw-mcp-server` | HTTP/Web Interface | Integrated | `pkg/daw/reaper.go` |

## 3. Universal & Utility
| Submodule | Key Functionality | Status | Integration Target |
|-----------|-------------------|--------|--------------------|
| `third_party/daw-mcp-ptaczek` | Abstract transport logic | Integrated | `pkg/daw/driver.go` |
| `third_party/scribbletune` | Music theory algorithms | Ported | `pkg/engine/theory.go` |
| `third_party/ableton-link` | Clock synchronization | Integrated | `pkg/engine/link.go` |

## Removal Roadmap
- **Batch 1 (Ableton):** `pylive`, `ableton-osc`, `ableton-live-tools` (Scheduled after Step 8.1)
- **Batch 2 (REAPER):** `reaper-reapy-mcp`, `reaper-daw-mcp-server` (Scheduled after Step 8.2)
- **Batch 3 (Universal):** `daw-mcp-ptaczek`, `scribbletune` (Scheduled after Step 8.3)

## 4. Newly Processed Architecture Submodules
| Submodule | Key Functionality | Status | Integration Target |
|-----------|-------------------|--------|--------------------|
| `third_party/logic-pro-mcp` | Logic Pro specific tools | Integrated | `data/schemas/logic-pro-mcp.json` |
| `third_party/ATRI_AGENT` | ATRI MIDI Tools | Integrated | `data/schemas/ATRI_AGENT_mcp.json` |
| `third_party/franz` | Arturia Pigments parameters | Integrated | `data/schemas/franz_mcp.json` |
| `third_party/cleaper` | REAPER full session IPC | Integrated | `data/schemas/cleaper_mcp.json` |

## 5. Extracted Ableton & Universal Extensions
All remaining `third_party/` DAWMCP submodules (`ableton-mcp-ahujasid`, `dawpilot-mcp`, `ardour-mcp`, `opendaw-mcp`, etc.) have been completely extracted into `data/schemas/` to keep the Go core lean, enabling their generic capabilities inside the `DAWDriver` ExecuteCustomCommand mapping loop.
