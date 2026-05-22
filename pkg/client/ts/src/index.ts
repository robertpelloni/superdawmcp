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
}
