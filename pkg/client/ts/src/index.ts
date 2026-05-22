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
  async createTrack(name: string, type: "audio" | "midi", daw?: string) {
    const args: any = { name, type };
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_create_track", arguments: args });
  }
  async transportControl(playing: boolean, bpm?: number, daw?: string) {
    const args: any = { playing };
    if (bpm !== undefined) args.bpm = bpm;
    if (daw) args.daw = daw;
    return await this.client.callTool({ name: "superdaw_transport_control", arguments: args });
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
