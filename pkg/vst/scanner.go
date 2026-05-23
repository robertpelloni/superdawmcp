package vst

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type PluginMetadata struct {
	Name       string          `json:"name"`
	Vendor     string          `json:"vendor"`
	Path       string          `json:"path"`
	Parameters []ParamMetadata `json:"parameters"`
}

type ParamMetadata struct {
	Name  string `json:"name"`
	Index int    `json:"index"`
}

type Scanner struct {
	cache     map[string]PluginMetadata
	cacheLock sync.RWMutex
	cachePath string
}

func NewScanner(cp string) *Scanner {
	s := &Scanner{
		cache:     make(map[string]PluginMetadata),
		cachePath: cp,
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
	data, _ := json.MarshalIndent(s.cache, "", "  ")
	os.WriteFile(s.cachePath, data, 0644)
}

func (s *Scanner) ScanDirectories(dirs []string) error {
	s.cacheLock.Lock()
	defer s.cacheLock.Unlock()

	if len(dirs) == 0 {
		if runtime.GOOS == "windows" {
			dirs = []string{`C:\Program Files\Common Files\VST3`}
		} else if runtime.GOOS == "darwin" {
			dirs = []string{"/Library/Audio/Plug-Ins/VST3", filepath.Join(os.Getenv("HOME"), "/Library/Audio/Plug-Ins/VST3")}
		} else {
			dirs = []string{"/usr/lib/vst3", "/usr/local/lib/vst3"}
		}
	}

	for _, dir := range dirs {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if strings.HasSuffix(info.Name(), ".vst3") {
				name := strings.TrimSuffix(info.Name(), ".vst3")
				if _, exists := s.cache[name]; !exists {
					// In a real scenario, we'd use a VST3 host library to probe parameters.
					// Pure Go binary parsing of VST3 is out of scope for a simple tool,
					// so we use a heuristic/placeholder for parameters.
					s.cache[name] = PluginMetadata{
						Name:   name,
						Vendor: "Unknown",
						Path:   path,
						Parameters: []ParamMetadata{
							{Name: "Volume", Index: 0},
							{Name: "Cutoff", Index: 1},
							{Name: "Resonance", Index: 2},
						},
					}
				}
			}
			return nil
		})
	}

	s.saveCache()
	return nil
}

func (s *Scanner) ListPlugins() []string {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()
	var list []string
	for k := range s.cache {
		list = append(list, k)
	}
	return list
}

func (s *Scanner) GetPluginMetadata(name string) (PluginMetadata, bool) {
	s.cacheLock.RLock()
	defer s.cacheLock.RUnlock()
	m, ok := s.cache[name]
	return m, ok
}
