package engine

import (
	"testing"
)

func TestGenerateEuclidean(t *testing.T) {
	// E(3, 8) -> [1, 0, 0, 1, 0, 0, 1, 0]
	notes := GenerateEuclidean(3, 8, 60, 100, 0, 4.0)

	if len(notes) != 3 {
		t.Errorf("Expected 3 notes, got %d", len(notes))
	}

	expectedStarts := []float32{0.0, 1.5, 3.0}
	for i, n := range notes {
		if n.StartBeat != expectedStarts[i] {
			t.Errorf("Note %d: expected start %f, got %f", i, expectedStarts[i], n.StartBeat)
		}
	}
}
