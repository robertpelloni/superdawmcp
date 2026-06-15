package client

import (
)

type Orchestrator struct {
	client *Client
	daws   []string
}

func NewOrchestrator(c *Client) *Orchestrator {
	return &Orchestrator{
		client: c,
		daws:   []string{"ableton", "reaper", "bitwig", "ardour", "flstudio", "logic", "cubase", "protools"},
	}
}

func (o *Orchestrator) StudioReset() {
	for _, daw := range o.daws {
		o.client.TransportControl(false, 120, daw)
		o.client.SetMixer("0", 0.8, 0, daw)
	}
}

func (o *Orchestrator) StudioPlay(bpm float64) {
	for _, daw := range o.daws {
		o.client.TransportControl(true, bpm, daw)
	}
}

func (o *Orchestrator) StudioSyncTempo(bpm float64) {
	for _, daw := range o.daws {
		o.client.TransportControl(true, bpm, daw) // Note: play=true to force sync on some daws
	}
}

func (o *Orchestrator) GetReport() map[string]string {
	report := make(map[string]string)
	for _, daw := range o.daws {
		_, err := o.client.ListClips("0", daw)
		if err == nil {
			report[daw] = "Online"
		} else {
			report[daw] = "Offline"
		}
	}
	return report
}
