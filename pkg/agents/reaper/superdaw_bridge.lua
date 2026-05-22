-- SuperDAW REAPER Bridge (Lua)
-- Exposes unified OSC control and ReaScript logic for REAPER

local function log(msg)
    reaper.ShowConsoleMsg("SuperDAW: " .. tostring(msg) .. "\n")
end

log("Bridge Initializing...")

-- Unified Handlers
local function SetTrackVolume(track_idx, volume)
    local track = reaper.GetTrack(0, tonumber(track_idx))
    if track then
        reaper.SetMediaTrackInfo_Value(track, "D_VOL", volume)
    end
end

local function SetTrackPan(track_idx, pan)
    local track = reaper.GetTrack(0, tonumber(track_idx))
    if track then
        reaper.SetMediaTrackInfo_Value(track, "D_PAN", pan)
    end
end

local function SetTrackMute(track_idx, mute)
    local track = reaper.GetTrack(0, tonumber(track_idx))
    if track then
        reaper.SetMediaTrackInfo_Value(track, "B_MUTE", mute)
    end
end

local function SetTrackSolo(track_idx, solo)
    local track = reaper.GetTrack(0, tonumber(track_idx))
    if track then
        reaper.SetMediaTrackInfo_Value(track, "I_SOLO", solo)
    end
end

local function WriteMIDI(track_idx, clip_idx, notes_json)
    local track = reaper.GetTrack(0, tonumber(track_idx))
    if not track then return end

    -- In REAPER, we typically create or find a MIDI item
    local item = reaper.GetTrackMediaItem(track, tonumber(clip_idx))
    if not item then
        -- Create a new 4-bar MIDI item if not found
        item = reaper.CreateNewMIDIItemInProj(track, 0, 16)
    end

    local take = reaper.GetActiveTake(item)
    if not take or not reaper.TakeIsMIDI(take) then return end

    -- Clear existing notes
    reaper.MIDI_DeleteNote(take, -1)

    -- Parse JSON (Mock parsing as REAPER doesn't have native JSON)
    -- In a real scenario, we'd use a Lua JSON library or pass flattened args
    -- For now, we assume a flattened string format or simple CSV
    log("WriteMIDI called for track " .. track_idx)
end

-- OSC Configuration is handled by REAPER's native OSC engine.
-- This script provides the ReaScript backend for custom logic.

log("SuperDAW REAPER Bridge Ready")
