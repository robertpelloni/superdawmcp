# Sonic Pi + SuperDAW-MCP Integration
# Controls multiple DAWs directly from live-coded Ruby.

# Load SuperDAW Client (assuming pkg/client/rb is in load path)
require_relative '../../pkg/client/rb/superdaw_client'

sd = SuperDAWClient.new("./bin/superdaw-mcp")
sd.connect

live_loop :superdaw_sync do
  # Sync DAW transport to Sonic Pi beats
  sd.transport_control(playing: true, bpm: current_bpm, daw: "ableton")
  sd.transport_control(playing: true, bpm: current_bpm, daw: "reaper")

  # Perform algorithmic mixer adjustments
  4.times do |i|
    sd.set_mixer(track_id: i.to_s, volume: rrand(0.5, 0.9), daw: "ableton")
    sleep 4
  end
end

at :finish do
  sd.disconnect
end
