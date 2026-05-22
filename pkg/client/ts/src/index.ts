import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { StdioClientTransport } from "@modelcontextprotocol/sdk/client/stdio.js";

export interface MIDINote {
  pitch: number;
  velocity: number;
  start_beat: number;
  duration: number;
}

export class SuperDAWClient {
  private client: Client;
  private transport: StdioClientTransport;

  constructor(serverPath: string, args: string[] = []) {
    this.transport = new StdioClientTransport({
      command: serverPath,
      args: args,
    });

    this.client = new Client(
      {
        name: "superdaw-ts-client",
        version: "1.0.0",
      },
      {
        capabilities: {},
      }
    );
  }

  async connect() {
    await this.client.connect(this.transport);
  }

  async setMixer(trackId: string, volume: number, pan: number) {
    return await this.client.callTool({
      name: "superdaw_set_mixer",
      arguments: {
        track_id: trackId,
        volume: volume,
        pan: pan,
      },
    });
  }

  async writeMIDI(trackId: string, notes: MIDINote[]) {
    return await this.client.callTool({
      name: "superdaw_write_midi",
      arguments: {
        track_id: trackId,
        notes: notes,
      },
    });
  }

  async disconnect() {
    await this.transport.close();
  }
}
