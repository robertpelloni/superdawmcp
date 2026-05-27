using UnityEngine;
using SuperDAW.Client;
using System.Threading.Tasks;

/**
 * SuperDAW-MCP Unity Integration Example
 * This script allows Unity projects to orchestrate real DAWs
 * for interactive audio experiences or game-audio sync.
 */
public class StudioManager : MonoBehaviour
{
    private SuperDAWClient _client;
    public string serverPath = "./bin/superdaw-mcp";

    async void Start()
    {
        _client = new SuperDAWClient(serverPath);
        _client.Connect();

        Debug.Log("Unity: Connected to SuperDAW-MCP Core.");

        // Sync studio to game tempo
        await _client.TransportControlAsync(true, 120.0, "ableton");
    }

    public async void OnGameTriggerEnter()
    {
        // Change mixer state based on game event
        await _client.SetMixerAsync("1", 1.0f, 0.5f, "ableton");
        Debug.Log("Unity: Studio Mixer adjusted by game trigger.");
    }

    void OnApplicationQuit()
    {
        _client?.Dispose();
    }
}
