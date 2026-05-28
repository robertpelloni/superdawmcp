import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";
export class SuperDAWClient {
  private client: Client;
  private transport: StdioClientTransport;
  constructor(serverPath: string, args: string[] = []) {
    this.transport = new StdioClientTransport({ command: serverPath, args: args });
    this.client = new Client({ name: "superdaw-ts-client", version: "1.0.0" }, { capabilities: {} });
  }
  async connect() { await this.client.connect(this.transport); }
  async disconnect() { await this.transport.close(); }
  async setMixer(trackId: string, volume: number, pan: number = 0, daw?: string) {
    const args: any = { track_id: trackId, volume, pan };
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_set_mixer", arguments: args });
  }
  async writeMIDI(trackId: string, notes: any[], daw?: string) {
    const args: any = { track_id: trackId, notes };
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_write_midi", arguments: args });
  }
  async createTrack(name: string, type: "audio" | "midi", daw?: string) {
    const args: any = { name, type };
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_create_track", arguments: args });
  }
  async transportControl(playing: boolean, bpm?: number, daw?: string) {
    const args: any = { playing };
    if (bpm) args.bpm = bpm;
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_transport_control", arguments: args });
  }
  async generateEuclidean(trackId: string, hits: number, steps: number, pitch: number, daw?: string) {
    const args: any = { track_id: trackId, hits, steps, pitch };
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_generate_euclidean", arguments: args });
  }
  async separateStems(inputPath: string, outputDir: string, stems: number = 4) {
    return await this.client.callTool({ name: "superdaw_separate_stems", arguments: { input_path: inputPath, output_dir: outputDir, stems } });
  }
  async listClips(trackId: string, daw?: string) {
    const args: any = { track_id: trackId };
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_list_clips", arguments: args });
  }
  async deleteClip(trackId: string, clipIdx: number, daw?: string) {
    const args: any = { track_id: trackId, clip_idx: clipIdx };
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_delete_clip", arguments: args });
  }
  async listPlugins() {
    return await this.client.callTool({ name: "superdaw_list_plugins", arguments: {} });
  }
  async getPluginParams(pluginName: string) {
    return await this.client.callTool({ name: "superdaw_get_plugin_params", arguments: { plugin_name: pluginName } });
  }
}

export class DAWAdapter {
  constructor(protected client: SuperDAWClient, protected dawName: string) {}
  async play() { return this.client.transportControl(true, undefined, this.dawName); }
  async stop() { return this.client.transportControl(false, undefined, this.dawName); }
}

export class AbletonLive extends DAWAdapter {
  constructor(client: SuperDAWClient) { super(client, "ableton"); }
  async fireScene(index: number) {
    // Custom command implementation
  }
}

export class Reaper extends DAWAdapter {
  constructor(client: SuperDAWClient) { super(client, "reaper"); }
}

export class LogicPro extends DAWAdapter {
  constructor(client: SuperDAWClient) { super(client, "logic"); }
}

export class Bitwig extends DAWAdapter {
  constructor(client: SuperDAWClient) { super(client, "bitwig"); }
}

export class FLStudio extends DAWAdapter {
  constructor(client: SuperDAWClient) { super(client, "flstudio"); }
}

export class Cubase extends DAWAdapter {
  constructor(client: SuperDAWClient) { super(client, "cubase"); }
}

export class Ardour extends DAWAdapter {
  constructor(client: SuperDAWClient) { super(client, "ardour"); }
}

export class ProTools extends DAWAdapter {
  constructor(client: SuperDAWClient) { super(client, "protools"); }
}
