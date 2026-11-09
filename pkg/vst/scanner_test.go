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

func TestParamMetadata_ClampingAndReset(t *testing.T) {
	p := ParamMetadata{
		Name:         "Filter Cutoff",
		DefaultValue: 0.55,
		MinValue:     0.0,
		MaxValue:     1.0,
		IsBipolar:    false,
		SmoothingMs:  8.0,
		CanAutomate:  true,
	}

	if clamped := p.ClampValue(1.5); clamped != 1.0 {
		t.Errorf("Expected 1.5 clamped to 1.0, got %f", clamped)
	}

	if clamped := p.ClampValue(-0.5); clamped != 0.0 {
		t.Errorf("Expected -0.5 clamped to 0.0, got %f", clamped)
	}

	if def := p.ResetToDefault(); def != 0.55 {
		t.Errorf("Expected default value 0.55, got %f", def)
	}
}

func TestParamMetadata_GenerateCurvePoints(t *testing.T) {
	p := ParamMetadata{
		Name:     "Cutoff",
		MinValue: 20.0,
		MaxValue: 20000.0,
		Curve:    "logarithmic",
	}

	pts := p.GenerateCurvePoints(10)
	if len(pts) != 10 {
		t.Fatalf("Expected 10 curve points, got %d", len(pts))
	}
	if pts[0] != 0.0 {
		t.Errorf("Expected first normalized curve point to be 0.0, got %f", pts[0])
	}
	if pts[9] != 1.0 {
		t.Errorf("Expected last normalized curve point to be 1.0, got %f", pts[9])
	}
	// Check strictly monotonic increasing for logarithmic frequency response
	for i := 1; i < len(pts); i++ {
		if pts[i] <= pts[i-1] {
			t.Errorf("Expected strictly increasing points at %d: %f <= %f", i, pts[i], pts[i-1])
		}
	}
}

func TestScanner_ParamLockAndBatchUpdate(t *testing.T) {
	cachePath := "/tmp/test_vst3_lock_cache.json"
	os.Remove(cachePath)
	defer os.Remove(cachePath)

	scanner := NewScanner(cachePath)
	err := scanner.ScanDirectories(nil)
	if err != nil {
		t.Fatalf("ScanDirectories failed: %v", err)
	}

	pluginName := "Surge XT"
	// Verify initial lock status is false
	locked, err := scanner.GetParamLock(pluginName, 0)
	if err != nil {
		t.Fatalf("GetParamLock failed: %v", err)
	}
	if locked {
		t.Errorf("Expected initial lock to be false, got true")
	}

	// Lock parameter 0 (Master Volume)
	err = scanner.SetParamLock(pluginName, 0, true)
	if err != nil {
		t.Fatalf("SetParamLock true failed: %v", err)
	}

	locked, err = scanner.GetParamLock(pluginName, 0)
	if err != nil || !locked {
		t.Errorf("Expected parameter 0 to be locked, got locked=%v, err=%v", locked, err)
	}

	// Batch update parameters 0 (locked) and 1 (unlocked)
	origVal0 := float32(0.80)
	updates := map[int]float32{
		0: 0.15, // Should be skipped because parameter 0 is locked
		1: 0.90, // Should be applied because parameter 1 is unlocked
	}

	applied, err := scanner.BatchUpdateParams(pluginName, updates)
	if err != nil {
		t.Fatalf("BatchUpdateParams failed: %v", err)
	}

	if _, exists := applied[0]; exists {
		t.Errorf("Expected locked param 0 to not be in applied map")
	}
	if applied[1] != 0.90 {
		t.Errorf("Expected param 1 to be applied with 0.90, got %+v", applied)
	}

	// Verify metadata reflects unchanged param 0 and updated param 1
	meta, ok := scanner.GetPluginMetadata(pluginName)
	if !ok {
		t.Fatalf("Failed to retrieve metadata for %s", pluginName)
	}
	if meta.Parameters[0].DefaultValue != origVal0 {
		t.Errorf("Locked parameter 0 default value changed from %f to %f!", origVal0, meta.Parameters[0].DefaultValue)
	}
	if meta.Parameters[1].DefaultValue != 0.90 {
		t.Errorf("Unlocked parameter 1 default value was not updated, got %f", meta.Parameters[1].DefaultValue)
	}

	// Unlock parameter 0
	err = scanner.SetParamLock(pluginName, 0, false)
	if err != nil {
		t.Fatalf("SetParamLock false failed: %v", err)
	}

	// Re-attempt batch update now that parameter 0 is unlocked
	applied2, err := scanner.BatchUpdateParams(pluginName, map[int]float32{0: 0.45})
	if err != nil {
		t.Fatalf("BatchUpdateParams after unlock failed: %v", err)
	}
	if applied2[0] != 0.45 {
		t.Errorf("Expected param 0 to be applied after unlocking, got applied=%+v", applied2)
	}
}

func TestScanner_ParamFavoritesAndTagging(t *testing.T) {
	cachePath := "/tmp/test_vst3_fav_cache.json"
	os.Remove(cachePath)
	defer os.Remove(cachePath)

	scanner := NewScanner(cachePath)
	err := scanner.ScanDirectories(nil)
	if err != nil {
		t.Fatalf("ScanDirectories failed: %v", err)
	}

	pluginName := "Surge XT"
	// Check initial favorites
	favs, err := scanner.GetFavoriteParams(pluginName)
	if err != nil {
		t.Fatalf("GetFavoriteParams failed: %v", err)
	}
	if len(favs) == 0 {
		t.Fatalf("Expected default favorite parameters, got 0")
	}

	// Toggle parameter 0 (Master Volume) from favorite (true) to false
	err = scanner.SetParamFavorite(pluginName, 0, false)
	if err != nil {
		t.Fatalf("SetParamFavorite to false failed: %v", err)
	}

	meta, ok := scanner.GetPluginMetadata(pluginName)
	if !ok {
		t.Fatalf("Failed to retrieve metadata for %s", pluginName)
	}
	if meta.Parameters[0].IsFavorite {
		t.Errorf("Expected parameter 0 IsFavorite to be false, got true")
	}

	// Add custom semantic tag
	err = scanner.AddParamTag(pluginName, 0, "LeadSynth")
	if err != nil {
		t.Fatalf("AddParamTag failed: %v", err)
	}

	meta, _ = scanner.GetPluginMetadata(pluginName)
	foundTag := false
	for _, tg := range meta.Parameters[0].Tags {
		if tg == "LeadSynth" {
			foundTag = true
			break
		}
	}
	if !foundTag {
		t.Errorf("Expected 'LeadSynth' tag in parameter 0 tags: %+v", meta.Parameters[0].Tags)
	}

	// Test persistence across new scanner instance loading same cache
	scanner2 := NewScanner(cachePath)
	meta2, ok := scanner2.GetPluginMetadata(pluginName)
	if !ok {
		t.Fatalf("Failed to retrieve metadata from cache reload")
	}
	if meta2.Parameters[0].IsFavorite {
		t.Errorf("Expected parameter 0 IsFavorite to persist as false across reload")
	}
}

func TestScanner_MacroAssignments(t *testing.T) {
	cachePath := "/tmp/test_vst3_macro_cache.json"
	os.Remove(cachePath)
	defer os.Remove(cachePath)

	scanner := NewScanner(cachePath)
	err := scanner.ScanDirectories(nil)
	if err != nil {
		t.Fatalf("ScanDirectories failed: %v", err)
	}

	pluginName := "Surge XT"

	// 1. Verify default seeded macro assignments (Filter Cutoff should have 'brighten': 1.0)
	meta, ok := scanner.GetPluginMetadata(pluginName)
	if !ok {
		t.Fatalf("Failed to retrieve metadata for %s", pluginName)
	}

	foundCutoff := false
	for _, p := range meta.Parameters {
		if p.Name == "Filter Cutoff" {
			foundCutoff = true
			if p.MacroAssignments == nil || p.MacroAssignments["brighten"] != 1.0 {
				t.Errorf("Expected Cutoff to have 'brighten': 1.0, got: %+v", p.MacroAssignments)
			}
		}
	}
	if !foundCutoff {
		t.Fatalf("Filter Cutoff parameter not found in %s", pluginName)
	}

	// 2. Query parameters assigned to 'brighten' macro
	brightenParams, err := scanner.GetMacroParams(pluginName, "brighten")
	if err != nil {
		t.Fatalf("GetMacroParams failed: %v", err)
	}
	if len(brightenParams) < 2 { // Cutoff (1.0) and Resonance (0.3)
		t.Errorf("Expected at least 2 brighten params, got %d", len(brightenParams))
	}

	// 3. Assign a custom macro weight with clamping
	// Param 0 (Master Volume) -> assign 'warmth' with weight 1.5 (should clamp to 1.0)
	err = scanner.AssignParamMacro(pluginName, 0, "warmth", 1.5)
	if err != nil {
		t.Fatalf("AssignParamMacro failed: %v", err)
	}

	meta, _ = scanner.GetPluginMetadata(pluginName)
	if meta.Parameters[0].MacroAssignments["warmth"] != 1.0 {
		t.Errorf("Expected clamped weight 1.0 for warmth, got %f", meta.Parameters[0].MacroAssignments["warmth"])
	}

	// 4. Delete macro assignment with weight 0
	err = scanner.AssignParamMacro(pluginName, 0, "warmth", 0.0)
	if err != nil {
		t.Fatalf("AssignParamMacro delete failed: %v", err)
	}

	meta, _ = scanner.GetPluginMetadata(pluginName)
	if _, exists := meta.Parameters[0].MacroAssignments["warmth"]; exists {
		t.Errorf("Expected warmth macro to be removed from parameter 0 assignments")
	}

	// 5. Test persistence across scanner instances
	err = scanner.AssignParamMacro(pluginName, 0, "punch", 0.75)
	if err != nil {
		t.Fatalf("AssignParamMacro failed: %v", err)
	}

	scanner2 := NewScanner(cachePath)
	meta2, ok := scanner2.GetPluginMetadata(pluginName)
	if !ok {
		t.Fatalf("Failed to reload cache")
	}
	if meta2.Parameters[0].MacroAssignments["punch"] != 0.75 {
		t.Errorf("Expected punch macro weight 0.75 to persist across cache reload, got %f", meta2.Parameters[0].MacroAssignments["punch"])
	}
}
