package dashboard

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
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

type ChatMessage struct {
	Sender    string `json:"sender"`
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

type RoomEvent struct {
	Text      string `json:"text"`
	Timestamp int64  `json:"timestamp"`
}

type ClientPresence struct {
	ID          string `json:"id"`
	ActiveTrack string `json:"active_track"`
}

// RoomData holds only the serializable state of a room
type RoomData struct {
	ID       string              `json:"id"`
	DAWs     map[string]DAWState `json:"daws"`
	Patches  []AudioPatch        `json:"patches"`
	Jobs     []interface{}       `json:"jobs"`
	Messages []ChatMessage       `json:"messages"`
	Events   []RoomEvent         `json:"events"`
}

// RoomState holds the state and connected clients for a specific collaborative session
type RoomState struct {
	RoomData
	Presences map[*websocket.Conn]*ClientPresence `json:"-"`
	mu       sync.RWMutex
}

type DashboardState struct {
	rooms map[string]*RoomState
	mu    sync.RWMutex
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
	RoomID    string                 `json:"room_id"`
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
	// Create default room
	orchestrator.getOrCreateRoom("default")

	mux := http.NewServeMux()

	// Asynchronous MCP command bus trigger
	mux.HandleFunc("/api/call", func(w http.ResponseWriter, r *http.Request) {
		var req InternalCommand
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if req.RoomID == "" { req.RoomID = "default" }
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
		roomID := r.URL.Query().Get("room")
		if roomID == "" { roomID = "default" }

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
		roomID := r.URL.Query().Get("room")
		if roomID == "" { roomID = "default" }
		safeRoomID := html.EscapeString(roomID)
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
						const roomID = '%s';
						const ws = new WebSocket('ws://' + window.location.host + '/ws?room=' + roomID);
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
		`, safeRoomID)
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
	return orchestrator, mux
}

func (s *DashboardState) getOrCreateRoom(id string) *RoomState {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room, ok := s.rooms[id]; ok {
		return room
	}
	room := &RoomState{
		RoomData: RoomData{
			ID:       id,
			DAWs:     make(map[string]DAWState),
			Patches:  []AudioPatch{},
			Messages: []ChatMessage{},
			Events:   []RoomEvent{},
		},
		Presences: make(map[*websocket.Conn]*ClientPresence),
	}
	s.rooms[id] = room
	return room
}

func (r *RoomState) AddMessage(sender, text string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Messages = append(r.Messages, ChatMessage{
		Sender:    sender,
		Text:      text,
		Timestamp: time.Now().Unix(),
	})
	if len(r.Messages) > 50 {
		r.Messages = r.Messages[1:]
	}
}

func (r *RoomState) AddEvent(text string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.Events = append(r.Events, RoomEvent{
		Text:      text,
		Timestamp: time.Now().Unix(),
	})
	if len(r.Events) > 50 {
		r.Events = r.Events[1:]
	}
}

// Global update methods that default to "default" room or can be updated to take roomID
func (s *DashboardState) UpdateDAW(roomID string, name string, playing bool, bpm float64) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	state := room.DAWs[name]
	state.Name = name
	state.IsPlaying = playing
	state.BPM = bpm
	room.DAWs[name] = state
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) UpdateJobs(roomID string, jobs []interface{}) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	room.Jobs = jobs
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) GetState(roomID string) RoomData {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.RLock()
	defer room.mu.RUnlock()
	return room.RoomData
}

func (s *DashboardState) SetState(roomID string, state RoomData) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	room.RoomData = state
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) UpdateArrangement(roomID string, name string, arrangement string) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	state := room.DAWs[name]
	state.Arrangement = arrangement
	room.DAWs[name] = state
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) RemovePatch(roomID string, p AudioPatch) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	newPatches := []AudioPatch{}
	for _, patch := range room.Patches {
		if patch.SourceDAW == p.SourceDAW && patch.SourceTrack == p.SourceTrack &&
			patch.DestDAW == p.DestDAW && patch.DestTrack == p.DestTrack {
			continue
		}
		newPatches = append(newPatches, patch)
	}
	room.Patches = newPatches
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) AddPatch(roomID string, p AudioPatch) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.mu.Lock()
	room.Patches = append(room.Patches, p)
	room.mu.Unlock()
	room.broadcast()
}

func (s *DashboardState) AddEvent(roomID string, text string) {
	if roomID == "" { roomID = "default" }
	room := s.getOrCreateRoom(roomID)
	room.AddEvent(text)
	room.broadcast()
}

func (r *RoomState) broadcast() {
	r.mu.RLock()

	presences := make([]*ClientPresence, 0, len(r.Presences))
	for _, p := range r.Presences {
		presences = append(presences, p)
	}

	msg := struct {
		RoomData
		UserCount int               `json:"user_count"`
		Presences []*ClientPresence `json:"presences"`
	}{
		RoomData: r.RoomData,
		UserCount: len(r.Presences),
		Presences: presences,
	}

	data, _ := json.Marshal(msg)

	// Copy connections to avoid holding lock during I/O
	conns := make([]*websocket.Conn, 0, len(r.Presences))
	for conn := range r.Presences {
		conns = append(conns, conn)
	}
	r.mu.RUnlock()

	for _, conn := range conns {
		conn.WriteMessage(websocket.TextMessage, data)
	}
}
