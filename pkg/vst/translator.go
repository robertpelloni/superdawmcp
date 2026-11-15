package vst

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PresetParam struct {
	Name  string  `json:"name"`
	Value float32 `json:"value"`
}

type Preset struct {
	Name       string        `json:"name"`
	Format     string        `json:"format"` // e.g. "fxp", "vstpreset", "vital"
	Parameters []PresetParam `json:"parameters"`
}

// TranslatePreset maps a preset parameter tree from a source synth (e.g. Serum)
// into a target synth (e.g. Vital or Ableton Wavetable) using heuristic mapping.
func TranslatePreset(source Preset, targetFormat string) (Preset, error) {
	translated := Preset{
		Name:       source.Name + " (Translated)",
		Format:     targetFormat,
		Parameters: []PresetParam{},
	}

	// Very simple heuristic mapping rules for demonstration
	for _, param := range source.Parameters {
		tName := strings.ToLower(param.Name)
		var newName string

		// Basic translation rules (Serum -> Vital style)
		if strings.Contains(tName, "osc") && strings.Contains(tName, "a") {
			if strings.Contains(tName, "vol") {
				newName = "Osc 1 Volume"
			} else if strings.Contains(tName, "pan") {
				newName = "Osc 1 Pan"
			} else if strings.Contains(tName, "detune") {
				newName = "Osc 1 Detune"
			} else if strings.Contains(tName, "wt pos") {
				newName = "Osc 1 Wavetable Frame"
			}
		} else if strings.Contains(tName, "filter") {
			if strings.Contains(tName, "cutoff") {
				newName = "Filter 1 Cutoff"
			} else if strings.Contains(tName, "res") {
				newName = "Filter 1 Resonance"
			}
		} else if strings.Contains(tName, "env1") {
			if strings.Contains(tName, "att") {
				newName = "Env 1 Attack"
			} else if strings.Contains(tName, "dec") {
				newName = "Env 1 Decay"
			} else if strings.Contains(tName, "sus") {
				newName = "Env 1 Sustain"
			} else if strings.Contains(tName, "rel") {
				newName = "Env 1 Release"
			}
		}

		if newName != "" {
			translated.Parameters = append(translated.Parameters, PresetParam{
				Name:  newName,
				Value: param.Value,
			})
		}
	}

	if len(translated.Parameters) == 0 {
		return translated, fmt.Errorf("no parameters could be translated from format %s to %s", source.Format, targetFormat)
	}

	return translated, nil
}

// DumpJSON serializes the translated preset for MCP tools or external UI.
func (p *Preset) DumpJSON() string {
	b, _ := json.Marshal(p)
	return string(b)
}
