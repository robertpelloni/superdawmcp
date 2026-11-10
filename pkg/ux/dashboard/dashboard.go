package dashboard

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/robertpelloni/superdaw-mcp/pkg/vst"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type DashboardState struct {
	DAWs    map[string]DAWState `json:"daws"`
	Patches []AudioPatch        `json:"patches"`
	Jobs    []interface{}       `json:"jobs"`
	mu      sync.RWMutex
	clients map[*websocket.Conn]bool
}

type DAWState struct {
	Name        string  `json:"name"`
	IsPlaying   bool    `json:"is_playing"`
	BPM         float64 `json:"bpm"`
	Arrangement string  `json:"arrangement"`
}

type AudioPatch struct {
	SourceDAW   string `json:"source_daw"`
	SourceTrack string `json:"source_track"`
	DestDAW     string `json:"dest_daw"`
	DestTrack   string `json:"dest_track"`
}

type InternalCommand struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

var CommandBus = make(chan InternalCommand, 64)

// StartDashboard initiates the web dashboard server, websocket hub, and REST endpoints
func StartDashboard(port int, scanner *vst.Scanner) (*DashboardState, *http.ServeMux) {
	state := &DashboardState{
		DAWs:    make(map[string]DAWState),
		Patches: []AudioPatch{},
		clients: make(map[*websocket.Conn]bool),
	}

	mux := http.NewServeMux()

	// Asynchronous MCP command bus trigger
	mux.HandleFunc("/api/call", func(w http.ResponseWriter, r *http.Request) {
		var req InternalCommand
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		CommandBus <- req
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "queued", "command": req.Name})
	})

	// Synchronous VST plugin & parameter query endpoint
	mux.HandleFunc("/api/plugins", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if scanner == nil {
			json.NewEncoder(w).Encode(map[string]interface{}{"plugins": []string{}})
			return
		}

		pluginName := r.URL.Query().Get("name")
		deep := r.URL.Query().Get("deep") == "true"

		if pluginName != "" {
			if deep {
				meta, err := scanner.DeepScan(pluginName)
				if err != nil {
					http.Error(w, err.Error(), http.StatusNotFound)
					return
				}
				json.NewEncoder(w).Encode(meta)
				return
			}

			meta, ok := scanner.GetPluginMetadata(pluginName)
			if !ok {
				http.Error(w, "Plugin not found", http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(meta)
			return
		}

		// Return all plugins
		all := scanner.GetAllPlugins()
		json.NewEncoder(w).Encode(all)
	})

	// Graceful shutdown / Quit endpoint for tray and dashboard control
	mux.HandleFunc("/api/quit", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "shutting_down",
			"message": "SuperDAW core servers shutting down gracefully.",
		})
		go func() {
			time.Sleep(200 * time.Millisecond)
			os.Exit(0)
		}()
	})

	// WebSocket handler for real-time telemetry
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		state.mu.Lock()
		state.clients[conn] = true
		state.mu.Unlock()
		state.broadcast()
	})

	// OBS Overlay Subview
	mux.HandleFunc("/obs", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `
			<!DOCTYPE html>
			<html>
				<head>
					<meta charset="utf-8">
					<title>SuperDAW OBS Overlay</title>
					<style>
						body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; color: #fff; background: transparent; padding: 12px; font-weight: bold; text-shadow: 2px 2px 4px #000; margin: 0; }
						.obs-stat { font-size: 26px; color: #00ff88; margin-left: 8px; }
						.obs-badge { background: rgba(18,18,18,0.85); padding: 8px 14px; border-radius: 8px; margin: 6px 0; border-left: 5px solid #00bcd4; display: inline-flex; align-items: center; box-shadow: 0 4px 10px rgba(0,0,0,0.5); }
					</style>
				</head>
				<body>
					<div id="stats"></div>
					<script>
						const ws = new WebSocket('ws://' + window.location.host + '/ws');
						ws.onmessage = (event) => {
							const state = JSON.parse(event.data);
							let html = '';
							for (const name in state.daws) {
								const d = state.daws[name];
								if (d.is_playing) {
									html += '<div class="obs-badge">DAW: ' + name.toUpperCase() + ' <span class="obs-stat">' + d.bpm.toFixed(1) + ' BPM</span></div>';
								}
							}
							document.getElementById('stats').innerHTML = html || '<div class="obs-badge" style="border-left-color: #666; color: #aaa;">DAW Transport Standby</div>';
						};
					</script>
				</body>
			</html>
		`)
	})

	// Unified Dashboard Homepage
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		renderUnifiedDashboard(w)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf("127.0.0.1:%d", port),
		Handler: mux,
	}
	go server.ListenAndServe()
	return state, mux
}

func (s *DashboardState) UpdateDAW(name string, playing bool, bpm float64) {
	s.mu.Lock()
	state := s.DAWs[name]
	state.Name = name
	state.IsPlaying = playing
	state.BPM = bpm
	s.DAWs[name] = state
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) UpdateJobs(jobs []interface{}) {
	s.mu.Lock()
	s.Jobs = jobs
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) GetState() DashboardState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return *s
}

func (s *DashboardState) SetState(state DashboardState) {
	s.mu.Lock()
	s.DAWs = state.DAWs
	s.Patches = state.Patches
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) UpdateArrangement(name string, arrangement string) {
	s.mu.Lock()
	state := s.DAWs[name]
	state.Arrangement = arrangement
	s.DAWs[name] = state
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) RemovePatch(p AudioPatch) {
	s.mu.Lock()
	newPatches := []AudioPatch{}
	for _, patch := range s.Patches {
		if patch.SourceDAW == p.SourceDAW && patch.SourceTrack == p.SourceTrack &&
			patch.DestDAW == p.DestDAW && patch.DestTrack == p.DestTrack {
			continue
		}
		newPatches = append(newPatches, patch)
	}
	s.Patches = newPatches
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) AddPatch(p AudioPatch) {
	s.mu.Lock()
	s.Patches = append(s.Patches, p)
	s.mu.Unlock()
	s.broadcast()
}

func (s *DashboardState) broadcast() {
	s.mu.RLock()
	data, _ := json.Marshal(s)
	for conn := range s.clients {
		conn.WriteMessage(websocket.TextMessage, data)
	}
	s.mu.RUnlock()
}
