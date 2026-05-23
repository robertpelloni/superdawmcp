package daw
import ( "bufio"; "encoding/json"; "fmt"; "net"; "sync" )
type BitwigDriver struct { addr string; conn net.Conn; mu sync.Mutex; id int }
func NewBitwigDriver(h string, p int) *BitwigDriver { return &BitwigDriver{addr: fmt.Sprintf("%s:%d", h, p), id: 1} }
func (b *BitwigDriver) Connect(e string) error { b.mu.Lock(); defer b.mu.Unlock(); c, err := net.Dial("tcp", b.addr); b.conn = c; return err }
func (b *BitwigDriver) Disconnect() error { b.mu.Lock(); defer b.mu.Unlock(); if b.conn != nil { b.conn.Close(); b.conn = nil }; return nil }
func (b *BitwigDriver) call(m string, p map[string]interface{}) error {
	b.mu.Lock(); defer b.mu.Unlock(); if b.conn == nil { c, err := net.Dial("tcp", b.addr); if err != nil { return err }; b.conn = c }
	req := map[string]interface{}{"jsonrpc": "2.0", "method": m, "params": p, "id": b.id}; b.id++
	d, _ := json.Marshal(req); fmt.Fprintf(b.conn, "%s\n", string(d))
	_, err := bufio.NewReader(b.conn).ReadString('\n'); return err
}
func (b *BitwigDriver) SetTransportState(p bool, bpm float64) error {
	m := "transport.stop"; if p { m = "transport.play" }; return b.call(m, map[string]interface{}{"bpm": bpm})
}
func (b *BitwigDriver) CreateTrack(n, t string) (string, error) { return "id", b.call("track.create", map[string]interface{}{"name": n, "type": t}) }
func (b *BitwigDriver) SetTrackVolume(id string, v float32) error { return b.call("track.volume", map[string]interface{}{"index": id, "volume": v}) }
func (b *BitwigDriver) SetTrackPan(id string, p float32) error { return b.call("track.pan", map[string]interface{}{"index": id, "pan": p}) }
func (b *BitwigDriver) WriteMIDIClip(id string, idx int, n []MIDINote) error { return b.call("clip.set_notes", map[string]interface{}{"trackIndex": id, "slotIndex": idx, "notes": n}) }
func (b *BitwigDriver) ListClips(id string) ([]ClipInfo, error) { return []ClipInfo{}, b.call("clip.list", map[string]interface{}{"trackIndex": id}) }
func (b *BitwigDriver) DeleteClip(id string, idx int) error { return b.call("clip.delete", map[string]interface{}{"trackIndex": id, "slotIndex": idx}) }
