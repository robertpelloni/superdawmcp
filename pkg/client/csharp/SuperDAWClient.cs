using System;
using System.IO;
using System.Diagnostics;
using System.Text.Json;
using System.Threading.Tasks;
using System.Collections.Generic;

namespace SuperDAW.Client
{
    public class SuperDAWClient : IDisposable
    {
        private Process _process;
        private StreamWriter _stdin;
        private StreamReader _stdout;
        private int _requestId = 1;

        public SuperDAWClient(string serverPath)
        {
            _process = new Process
            {
                StartInfo = new ProcessStartInfo
                {
                    FileName = serverPath,
                    RedirectStandardInput = true,
                    RedirectStandardOutput = true,
                    UseShellExecute = false,
                    CreateNoWindow = true
                }
            };
        }

        public void Connect()
        {
            _process.Start();
            _stdin = _process.StandardInput;
            _stdout = _process.StandardOutput;
        }

        public async Task<JsonElement> CallAsync(string method, object parameters)
        {
            var request = new
            {
                jsonrpc = "2.0",
                method = method,
                params = parameters,
                id = _requestId++
            };

            string json = JsonSerializer.Serialize(request);
            await _stdin.WriteLineAsync(json);

            string response = await _stdout.ReadLineAsync();
            using var doc = JsonDocument.Parse(response);
            return doc.RootElement.GetProperty("result").Clone();
        }

        public async Task SetMixerAsync(string trackId, float volume, float pan = 0.0f, string daw = null)
        {
            var args = new Dictionary<string, object>
            {
                { "track_id", trackId },
                { "volume", volume },
                { "pan", pan }
            };
            if (daw != null) args["daw"] = daw;

            await CallAsync("tools/call", new { name = "superdaw_set_mixer", arguments = args });
        }

        public async Task TransportControlAsync(bool playing, double? bpm = null, string daw = null)
        {
            var args = new Dictionary<string, object> { { "playing", playing } };
            if (bpm.HasValue) args["bpm"] = bpm.Value;
            if (daw != null) args["daw"] = daw;

            await CallAsync("tools/call", new { name = "superdaw_transport_control", arguments = args });
        }

        public void Dispose()
        {
            _stdin?.Close();
            _process?.WaitForExit();
            _process?.Dispose();
        }
    }
}
