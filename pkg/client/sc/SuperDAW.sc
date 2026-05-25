/**
 * SuperDAW-MCP SuperCollider Client
 * Provides a high-level SC interface for DAW orchestration.
 */

SuperDAW {
    var <pipe, <id=1;

    *new { |serverPath="./bin/superdaw-mcp"|
        ^super.new.init(serverPath);
    }

    init { |serverPath|
        pipe = Pipe.new(serverPath, "w");
    }

    send { |method, params|
        var req = (
            jsonrpc: "2.0",
            method: method,
            params: params,
            id: id
        );
        id = id + 1;
        pipe.write(req.asJSON ++ "\n");
    }

    transport { |playing=true, bpm=120, daw="ableton"|
        this.send("tools/call", (
            name: "superdaw_transport_control",
            arguments: (playing: playing, bpm: bpm, daw: daw)
        ));
    }

    mixer { |trackID="0", vol=0.8, pan=0.0, daw="ableton"|
        this.send("tools/call", (
            name: "superdaw_set_mixer",
            arguments: (track_id: trackID, volume: vol, pan: pan, daw: daw)
        ));
    }

    close {
        pipe.close;
    }
}
