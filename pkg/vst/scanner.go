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
	Version    string          `json:"version"`
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
					vendor := "Unknown"
					version := "1.0.0"

					// Improved Vendor Discovery via path heuristics
					pLower := strings.ToLower(path)
					if strings.Contains(pLower, "fabfilter") { vendor = "FabFilter" }
					else if strings.Contains(pLower, "waves") { vendor = "Waves" }
					else if strings.Contains(pLower, "u-he") { vendor = "u-he" }
					else if strings.Contains(pLower, "arturia") { vendor = "Arturia" }
					else if strings.Contains(pLower, "izotope") { vendor = "iZotope" }
					else if strings.Contains(pLower, "native instruments") { vendor = "Native Instruments" }

					// Parameter heuristics for common plugin types
					params := []ParamMetadata{
						{Name: "Volume", Index: 0},
						{Name: "Resonance", Index: 1},
						{Name: "Attack", Index: 2},
						{Name: "Release", Index: 3},
					}

					s.cache[name] = PluginMetadata{
						Name:    name,
						Vendor:  vendor,
						Version: version,
						Path:    path,
						Parameters: params,
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
