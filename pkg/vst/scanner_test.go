package vst

import (
	"os"
	"strings"
	"testing"
)

func TestScanner_DiscoveryAndDeepScan(t *testing.T) {
	cachePath := "/tmp/test_vst3_cache.json"
	os.Remove(cachePath)
	defer os.Remove(cachePath)

	scanner := NewScanner(cachePath)
	err := scanner.ScanDirectories(nil)
	if err != nil {
		t.Fatalf("ScanDirectories failed: %v", err)
	}

	plugins := scanner.ListPlugins()
	if len(plugins) == 0 {
		t.Fatalf("Expected default plugins to be registered in cache, got 0")
	}

	// Verify Vital / Synth mapping
	synthMeta, ok := scanner.GetPluginMetadata("Vital")
	if !ok {
		t.Fatalf("Vital metadata not found in cache")
	}
	if len(synthMeta.Parameters) == 0 {
		t.Fatalf("Vital parameters should not be empty")
	}

	// Verify parameter classification, display formatting, and curve properties
	foundOsc := false
	foundFilter := false
	for _, p := range synthMeta.Parameters {
		if p.Name == "Osc 1 Waveform" {
			foundOsc = true
			if p.Category != "Filter/Tone" || p.Step != 1.0 {
				t.Errorf("Unexpected Osc 1 mapping: category=%s, step=%f", p.Category, p.Step)
			}
			formatted := p.FormatDisplayValue(0.0)
			if formatted != "Sine" {
				t.Errorf("Expected Sine waveform format, got %s", formatted)
			}
		}
		if p.Name == "Filter Cutoff" {
			foundFilter = true
			if p.Curve != "logarithmic" {
				t.Errorf("Filter Cutoff curve should be logarithmic, got %s", p.Curve)
			}
			formatted := p.FormatDisplayValue(0.55)
			if !strings.Contains(formatted, "Hz") && !strings.Contains(formatted, "kHz") {
				t.Errorf("Expected Hz or kHz formatting, got %s", formatted)
			}
		}
	}
	if !foundOsc || !foundFilter {
		t.Errorf("Expected Osc 1 Waveform and Filter Cutoff parameters to be present")
	}

	// Verify Compressor / Dynamics mapping
	compMeta, ok := scanner.GetPluginMetadata("FabFilter Pro-C 2")
	if !ok {
		t.Fatalf("FabFilter Pro-C 2 metadata not found")
	}
	if compMeta.DeepMeta.Category != "Fx|Dynamics" {
		t.Errorf("Expected Fx|Dynamics category, got %s", compMeta.DeepMeta.Category)
	}
	foundThresh := false
	for _, p := range compMeta.Parameters {
		if p.Name == "Threshold" {
			foundThresh = true
			if p.Unit != "dB" || p.Curve != "decibel" {
				t.Errorf("Unexpected Threshold unit/curve: %s / %s", p.Unit, p.Curve)
			}
			formatted := p.FormatDisplayValue(-18.0)
			if formatted != "-18.0 dB" {
				t.Errorf("Expected -18.0 dB formatted, got %s", formatted)
			}
		}
	}
	if !foundThresh {
		t.Errorf("Expected Threshold parameter in dynamics plugin")
	}

	// Verify DeepScan
	updated, err := scanner.DeepScan("Vital")
	if err != nil {
		t.Fatalf("DeepScan failed: %v", err)
	}
	if !updated.DeepMeta.IsDeepScanned {
		t.Errorf("Expected IsDeepScanned to be true")
	}
	if updated.DeepMeta.SymbolCount <= 0 {
		t.Errorf("Expected discovered symbols count > 0, got %d", updated.DeepMeta.SymbolCount)
	}
}

func TestScanner_ParamModulationRouting(t *testing.T) {
	cachePath := "/tmp/test_vst3_mod_cache.json"
	os.Remove(cachePath)
	defer os.Remove(cachePath)

	scanner := NewScanner(cachePath)
	err := scanner.ScanDirectories(nil)
	if err != nil {
		t.Fatalf("ScanDirectories failed: %v", err)
	}

	// Assign LFO 1 modulation to Filter Cutoff (index 5 on Surge XT)
	err = scanner.SetParamModulation("Surge XT", 5, "LFO 1", 0.65)
	if err != nil {
		t.Fatalf("SetParamModulation failed: %v", err)
	}

	mod, err := scanner.GetParamModulation("Surge XT", 5)
	if err != nil {
		t.Fatalf("GetParamModulation failed: %v", err)
	}
	if mod.Source != "LFO 1" || mod.Depth != 0.65 {
		t.Errorf("Expected LFO 1 with depth 0.65, got %+v", mod)
	}

	// Test boundary clamp on depth
	err = scanner.SetParamModulation("Surge XT", 5, "Envelope 1", 2.5) // clamps to 1.0
	if err != nil {
		t.Fatalf("SetParamModulation failed: %v", err)
	}
	mod, _ = scanner.GetParamModulation("Surge XT", 5)
	if mod.Depth != 1.0 {
		t.Errorf("Expected depth clamp to 1.0, got %f", mod.Depth)
	}

	// Test invalid parameter index error
	err = scanner.SetParamModulation("Surge XT", 999, "Velocity", 0.5)
	if err == nil {
		t.Errorf("Expected error for non-existent parameter index, got nil")
	}
}

func TestParamMetadata_CurveConversionsAndSnapping(t *testing.T) {
	// 1. Logarithmic Frequency Curve (20 Hz to 20,000 Hz)
	logParam := ParamMetadata{
		Name:     "Cutoff",
		MinValue: 20.0,
		MaxValue: 20000.0,
		Curve:    "logarithmic",
		Step:     1.0,
	}

	normVal := float32(0.5)
	physVal := logParam.NormalizedToPhysical(normVal)
	roundTripNorm := logParam.PhysicalToNormalized(physVal)
	if (roundTripNorm - normVal) > 0.001 || (normVal - roundTripNorm) > 0.001 {
		t.Errorf("Log curve round-trip failed: initial %f, physical %f, round-trip %f", normVal, physVal, roundTripNorm)
	}

	// 2. Decibel Curve (-60 dB to +6 dB)
	dbParam := ParamMetadata{
		Name:     "Gain",
		MinValue: -60.0,
		MaxValue: 6.0,
		Curve:    "decibel",
		Step:     0.1,
	}

	normDb := float32(0.75)
	physDb := dbParam.NormalizedToPhysical(normDb)
	roundTripDb := dbParam.PhysicalToNormalized(physDb)
	if (roundTripDb - normDb) > 0.001 || (normDb - roundTripDb) > 0.001 {
		t.Errorf("Decibel curve round-trip failed: initial %f, physical %f, round-trip %f", normDb, physDb, roundTripDb)
	}

	// 3. S-Curve Mapping
	sParam := ParamMetadata{
		Name:     "Morph",
		MinValue: 0.0,
		MaxValue: 100.0,
		Curve:    "scurve",
		Step:     0.5,
	}

	normS := float32(0.3)
	physS := sParam.NormalizedToPhysical(normS)
	roundTripS := sParam.PhysicalToNormalized(physS)
	if (roundTripS - normS) > 0.005 || (normS - roundTripS) > 0.005 {
		t.Errorf("S-curve round-trip failed: initial %f, physical %f, round-trip %f", normS, physS, roundTripS)
	}

	// 4. Stepped Snapping
	stepParam := ParamMetadata{
		Name:     "Octave",
		MinValue: -2.0,
		MaxValue: 2.0,
		Step:     1.0,
	}

	snapped := stepParam.SnapToStep(0.78)
	if snapped != 1.0 {
		t.Errorf("Expected 0.78 snapped to 1.0, got %f", snapped)
	}

	snappedNeg := stepParam.SnapToStep(-1.3)
	if snappedNeg != -1.0 {
		t.Errorf("Expected -1.3 snapped to -1.0, got %f", snappedNeg)
	}
}
