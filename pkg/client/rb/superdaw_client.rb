require 'json'
require 'open3'

module SuperDAW
  class Client
    def initialize(server_path)
      @server_path = server_path
      @request_id = 1
    end

    def connect
      @stdin, @stdout, @wait_thr = Open3.popen2(@server_path)
    end

    def call(method, params)
      request = {
        jsonrpc: "2.0",
        method: method,
        params: params,
        id: @request_id
      }
      @request_id += 1

      @stdin.puts(JSON.generate(request))
      response = @stdout.gets
      JSON.parse(response)['result']
    end

    def set_mixer(track_id, volume, pan: 0.0, daw: nil)
      args = { track_id: track_id, volume: volume, pan: pan }
      args[:daw] = daw if daw
      call("tools/call", { name: "superdaw_set_mixer", arguments: args })
    end

    def transport_control(playing, bpm: nil, daw: nil)
      args = { playing: playing }
      args[:bpm] = bpm if bpm
      args[:daw] = daw if daw
      call("tools/call", { name: "superdaw_transport_control", arguments: args })
    end

    def disconnect
      @stdin.close
      @wait_thr.join
    end
  end
end
