package vst

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
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
				// In a real implementation, we would load the plugin and query parameters.
				// Here we mock the metadata extraction.
				s.mockExtractMetadata(path)
			}
			return nil
		})
	}
	s.saveCache()
	return nil
}

func (s *Scanner) mockExtractMetadata(path string) {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	name := filepath.Base(path)
	s.cache[name] = PluginMetadata{
		Name: name,
		ID:   fmt.Sprintf("vst-%s", name),
		Parameters: []ParamMetadata{
			{Name: "Volume", Index: 0, Min: 0.0, Max: 1.0},
			{Name: "Cutoff", Index: 1, Min: 0.0, Max: 1.0},
		},
	}
}

func (s *Scanner) GetPluginMetadata(name string) (PluginMetadata, bool) {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()
	meta, ok := s.cache[name]
	return meta, ok
}
