package vst

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type PluginMetadata struct {
	Name       string           `json:"name"`
	ID         string           `json:"id"`
	Vendor     string           `json:"vendor"`
	Parameters []ParamMetadata `json:"parameters"`
}

type ParamMetadata struct {
	Name  string  `json:"name"`
	Index int     `json:"index"`
	Min   float32 `json:"min"`
	Max   float32 `json:"max"`
}

type Scanner struct {
	cache     map[string]PluginMetadata
	cacheLock sync.RWMutex
	cachePath string
}

func NewScanner(cachePath string) *Scanner {
	s := &Scanner{
		cache:     make(map[string]PluginMetadata),
		cachePath: cachePath,
	}
	s.loadCache()
	return s
}

func (s *Scanner) loadCache() {
	data, err := os.ReadFile(s.cachePath)
	if err == nil {
		json.Unmarshal(data, &s.cache)
	}
}

func (s *Scanner) saveCache() {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()
	data, _ := json.Marshal(s.cache)
	os.WriteFile(s.cachePath, data, 0644)
}

func (s *Scanner) ScanDirectories(dirs []string) error {
	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if !info.IsDir() && (filepath.Ext(path) == ".vst3" || filepath.Ext(path) == ".vst") {
				s.extractMetadata(path)
			}
			return nil
		})
	}
	s.saveCache()
	return nil
}

func (s *Scanner) extractMetadata(path string) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	name := filepath.Base(path)
	vendor := "Unknown"
	if filepath.Base(filepath.Dir(path)) != "VST3" {
		vendor = filepath.Base(filepath.Dir(path))
	}

	params := []ParamMetadata{
		{Name: "Bypass", Index: 0},
	}

	lowerName := strings.ToLower(name)

	if strings.Contains(lowerName, "comp") || strings.Contains(lowerName, "limiter") {
		params = append(params, []ParamMetadata{
			{Name: "Threshold", Index: 1}, {Name: "Ratio", Index: 2},
			{Name: "Attack", Index: 3}, {Name: "Release", Index: 4},
			{Name: "Gain", Index: 5},
		}...)
	} else if strings.Contains(lowerName, "eq") || strings.Contains(lowerName, "filter") {
		params = append(params, []ParamMetadata{
			{Name: "Frequency", Index: 1}, {Name: "Gain", Index: 2},
			{Name: "Q", Index: 3}, {Name: "Type", Index: 4},
		}...)
	} else {
		params = append(params, []ParamMetadata{
			{Name: "Volume", Index: 1}, {Name: "Cutoff", Index: 2},
			{Name: "Resonance", Index: 3}, {Name: "Mix", Index: 4},
		}...)
	}

	s.cache[name] = PluginMetadata{
		Name:       name,
		ID:         fmt.Sprintf("vst-%s", name),
		Vendor:     vendor,
		Parameters: params,
	}
}

func (s *Scanner) GetPluginMetadata(name string) (PluginMetadata, bool) {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()
	meta, ok := s.cache[name]
	return meta, ok
}

func (s *Scanner) ListPlugins() []string {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()
	var plugins []string
	for name := range s.cache {
		plugins = append(plugins, name)
	}
	return plugins
}
