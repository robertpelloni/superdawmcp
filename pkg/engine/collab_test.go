package engine

import (
	"testing"
	"time"
)

func TestCollabHub_PeerLifecycle(t *testing.T) {
	hub := NewCollabHub("Test-Jam")
	if hub.GetPeerCount() != 0 {
		t.Fatalf("expected 0 initial peers, got %d", hub.GetPeerCount())
	}

	// Register peer 1
	ch1, err := hub.RegisterPeer("peer-1", "Alice", true)
	if err != nil {
		t.Fatalf("failed to register peer-1: %v", err)
	}
	if ch1 == nil {
		t.Fatalf("expected valid channel for peer-1")
	}

	// Register peer 2
	ch2, err := hub.RegisterPeer("peer-2", "Bob", false)
	if err != nil {
		t.Fatalf("failed to register peer-2: %v", err)
	}
	if ch2 == nil {
		t.Fatalf("expected valid channel for peer-2")
	}

	if hub.GetPeerCount() != 2 {
		t.Fatalf("expected 2 peers, got %d", hub.GetPeerCount())
	}

	peers := hub.GetPeers()
	if len(peers) != 2 {
		t.Fatalf("expected 2 peers in list, got %d", len(peers))
	}

	// Unregister peer 1
	hub.UnregisterPeer("peer-1")
	if hub.GetPeerCount() != 1 {
		t.Fatalf("expected 1 peer after unregistering, got %d", hub.GetPeerCount())
	}

	// Verify channel for peer 1 was closed
	_, ok := <-ch1
	if ok {
		t.Fatalf("expected ch1 to be closed")
	}
}

func TestCollabHub_BroadcastAndLWW(t *testing.T) {
	hub := NewCollabHub("Test-Jam-LWW")
	ch1, _ := hub.RegisterPeer("peer-1", "Alice", true)
	ch2, _ := hub.RegisterPeer("peer-2", "Bob", false)

	now := time.Now()

	// Peer 1 broadcasts Cutoff = 0.75
	applied, err := hub.BroadcastParamChange("peer-1", "Surge XT", 5, 0.75, now)
	if err != nil || !applied {
		t.Fatalf("BroadcastParamChange failed: applied=%v, err=%v", applied, err)
	}

	// Peer 2 should receive the event
	select {
	case event := <-ch2:
		if event.SenderID != "peer-1" || event.PluginName != "Surge XT" || event.ParamIndex != 5 || event.Value != 0.75 {
			t.Fatalf("unexpected event received on peer-2: %+v", event)
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatalf("timed out waiting for event on peer-2")
	}

	// Peer 1 should NOT receive their own event
	select {
	case event := <-ch1:
		t.Fatalf("peer-1 should not receive own event: %+v", event)
	default:
		// Expected
	}

	// Verify consensus state
	state, ok := hub.GetParamState("Surge XT", 5)
	if !ok || state.Value != 0.75 || state.UpdatedBy != "peer-1" {
		t.Fatalf("unexpected consensus state: %+v", state)
	}

	// Test LWW: Peer 2 attempts to broadcast an OLDER event (now - 1 second)
	staleTime := now.Add(-1 * time.Second)
	appliedStale, err := hub.BroadcastParamChange("peer-2", "Surge XT", 5, 0.40, staleTime)
	if err != nil {
		t.Fatalf("stale broadcast returned error: %v", err)
	}
	if appliedStale {
		t.Fatalf("expected stale change to be rejected by LWW")
	}

	// Consensus value must remain 0.75
	stateAfter, _ := hub.GetParamState("Surge XT", 5)
	if stateAfter.Value != 0.75 {
		t.Fatalf("consensus state was overwritten by stale change: %f", stateAfter.Value)
	}

	// Newer event from Peer 2 (now + 1 second)
	newerTime := now.Add(1 * time.Second)
	appliedNewer, err := hub.BroadcastParamChange("peer-2", "Surge XT", 5, 0.90, newerTime)
	if err != nil || !appliedNewer {
		t.Fatalf("newer broadcast failed: applied=%v, err=%v", appliedNewer, err)
	}

	stateNewer, _ := hub.GetParamState("Surge XT", 5)
	if stateNewer.Value != 0.90 || stateNewer.UpdatedBy != "peer-2" {
		t.Fatalf("expected updated consensus value 0.90 from peer-2: %+v", stateNewer)
	}
}
