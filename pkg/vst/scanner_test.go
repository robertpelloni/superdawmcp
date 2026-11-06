package vst

import (
	"os"
	"testing"
)

func TestScanner_DiscoveryAndDeepScan(t *testing.T) {
	cacheFile := "test_vst_cache.json"
	defer os.Remove(cacheFile)

	scanner := NewScanner(cacheFile)
	if err := scanner.ScanDirectories([]string{}); err != nil {
		t.Fatalf("ScanDirectories failed: %v", err)
	}

	plugins := scanner.ListPlugins()
	if len(plugins) == 0 {
		t.Fatalf("Expected discovered plugins, got 0")
	}

	foundSurge := false
	for _, p := range plugins {
		if p == "Surge XT" {
			foundSurge = true
			break
		}
	}
	if !foundSurge {
		t.Errorf("Surge XT was not discovered in scanner defaults")
	}

	meta, ok := scanner.GetPluginMetadata("Surge XT")
	if !ok {
		t.Fatalf("Failed to retrieve metadata for Surge XT")
	}

	if len(meta.Parameters) == 0 {
		t.Errorf("Expected parameters for Surge XT, got 0")
	}

	deepMeta, err := scanner.DeepScan("Surge XT")
	if err != nil {
		t.Fatalf("DeepScan failed: %v", err)
	}

	if !deepMeta.DeepMeta.IsDeepScanned {
		t.Errorf("Expected IsDeepScanned to be true")
	}

	params, ok := scanner.GetPluginParameters("Surge XT")
	if !ok || len(params) == 0 {
		t.Errorf("Expected GetPluginParameters to succeed and return parameters")
	}
}
