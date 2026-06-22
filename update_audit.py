import re

audit_file = "SUBMODULE_AUDIT.md"
with open(audit_file, "r") as f:
    content = f.read()

new_content = content + """
## 4. Newly Processed Architecture Submodules
| Submodule | Key Functionality | Status | Integration Target |
|-----------|-------------------|--------|--------------------|
| `third_party/logic-pro-mcp` | Logic Pro specific tools | Integrated | `data/schemas/logic-pro-mcp.json` |
| `third_party/ATRI_AGENT` | ATRI MIDI Tools | Integrated | `data/schemas/ATRI_AGENT_mcp.json` |
| `third_party/franz` | Arturia Pigments parameters | Integrated | `data/schemas/franz_mcp.json` |
| `third_party/cleaper` | REAPER full session IPC | Integrated | `data/schemas/cleaper_mcp.json` |
"""

with open(audit_file, "w") as f:
    f.write(new_content)
