import { SuperDAWClient } from "./index.js";

export class SuperDAWOrchestrator {
    private client: SuperDAWClient;
    private daws: string[] = ["ableton", "reaper", "bitwig", "ardour", "flstudio", "logic", "cubase", "protools"];

    constructor(client: SuperDAWClient) {
        this.client = client;
    }

    async studioReset() {
        for (const daw of this.daws) {
            try {
                await this.client.transportControl(false, undefined, daw);
                await this.client.setMixer("0", 0.8, 0, daw);
            } catch (e) {
                // Ignore offline DAWs
            }
        }
    }

    async studioPlay(bpm: number = 120) {
        for (const daw of this.daws) {
            try {
                await this.client.transportControl(true, bpm, daw);
            } catch (e) { }
        }
    }

    async studioSyncTempo(bpm: number) {
        for (const daw of this.daws) {
            try {
                await this.client.transportControl(undefined, bpm, daw);
            } catch (e) { }
        }
    }
}
