-- SuperDAW REAPER Bridge (Lua)
-- Exposes unified OSC control for REAPER

local function log(msg)
    reaper.ShowConsoleMsg("SuperDAW: " .. tostring(msg) .. "\n")
end

log("Bridge Initializing...")

-- OSC Configuration
local osc_ip = "127.0.0.1"
local osc_port = 8000 -- REAPER listens here

local function handle_osc(msg)
    local addr = msg.address
    log("Received OSC: " .. addr)

    if addr == "/superdaw/transport/play" then
        local val = msg.args[1]
        if val == 1 then reaper.OnPlayButton() else reaper.OnStopButton() end
    elseif addr == "/superdaw/transport/tempo" then
        reaper.SetCurrentBPM(0, msg.args[1], true)
    elseif addr == "/superdaw/track/volume" then
        local track_idx = msg.args[1]
        local volume = msg.args[2]
        local track = reaper.GetTrack(0, tonumber(track_idx))
        if track then
            reaper.SetMediaTrackInfo_Value(track, "D_VOL", volume)
        end
    end
end

-- Main loop to keep script alive if necessary, though REAPER OSC is usually handled by its engine.
-- For a pure Lua implementation of OSC parsing, we would use a socket library.
-- However, REAPER's built-in OSC support is more efficient.
-- This script serves as a placeholder for custom ReaScript logic that standard OSC cannot reach.

log("SuperDAW REAPER Bridge Ready (Placeholder for advanced ReaScript logic)")
