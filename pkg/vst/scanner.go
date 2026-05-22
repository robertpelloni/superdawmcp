package vst
import (
	"sync"
)
type PluginMetadata struct {
	Name string `json:"name"`
	Vendor string `json:"vendor"`
	Parameters []ParamMetadata `json:"parameters"`
}
type ParamMetadata struct { Name string `json:"name"`; Index int `json:"index"` }
type Scanner struct { cache map[string]PluginMetadata; cacheLock sync.RWMutex; cachePath string }
func NewScanner(cp string) *Scanner { return &Scanner{cache: make(map[string]PluginMetadata), cachePath: cp} }
func (s *Scanner) ScanDirectories(dirs []string) error {
	s.cacheLock.Lock(); defer s.cacheLock.Unlock()
	s.cache["MockSynth"] = PluginMetadata{Name: "MockSynth", Vendor: "SuperDAW", Parameters: []ParamMetadata{{Name: "Volume", Index: 0}, {Name: "Cutoff", Index: 1}}}
	return nil
}
