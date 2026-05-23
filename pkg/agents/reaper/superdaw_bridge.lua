-- SuperDAW REAPER Bridge (OSC + ReaScript)
-- This script provides a native OSC listener within REAPER to handle
-- the unified SuperDAW protocol.

local function log(msg)
    reaper.ShowConsoleMsg("SuperDAW: " .. tostring(msg) .. "\n")
end

log("SuperDAW REAPER Bridge (OSC) Initializing...")

-- REAPER's built-in OSC is preferred for performance, but this script
-- can be used to handle custom logic that native OSC doesn't support.

function Main()
    -- Placeholder for per-tick logic
    reaper.defer(Main)
end

-- If using a Lua-based OSC server (requires LuaSocket), it would be initialized here.
-- For the purpose of this integration, we rely on REAPER's native OSC support
-- configured to listen on port 8000 and map to the SuperDAW schema.

log("SuperDAW REAPER Bridge Ready")
Main()
