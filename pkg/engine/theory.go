package engine

import (
	"strings"
)

// Scale represents a musical scale.
type Scale struct {
	Name      string
	Intervals []int
}

var Scales = map[string]Scale{
	"major":      {Name: "Major", Intervals: []int{0, 2, 4, 5, 7, 9, 11}},
	"minor":      {Name: "Minor", Intervals: []int{0, 2, 3, 5, 7, 8, 10}},
	"pentatonic": {Name: "Pentatonic", Intervals: []int{0, 2, 4, 7, 9}},
	"dorian":     {Name: "Dorian", Intervals: []int{0, 2, 3, 5, 7, 9, 10}},
	"phrygian":   {Name: "Phrygian", Intervals: []int{0, 1, 3, 5, 7, 8, 10}},
}

// GetScaleNotes returns the MIDI notes for a given root and scale name.
func GetScaleNotes(root int, scaleName string) []int {
	scale, ok := Scales[strings.ToLower(scaleName)]
	if !ok {
		scale = Scales["major"]
	}
	notes := make([]int, len(scale.Intervals))
	for i, interval := range scale.Intervals {
		notes[i] = (root + interval) % 12
	}
	return notes
}

// Progression patterns (Tonic-Predominant-Dominant)
var RomanNumeralPatterns = map[string][]string{
	"basic":   {"I", "IV", "V", "I"},
	"pop":     {"I", "V", "vi", "IV"},
	"jazz":    {"ii", "V", "I", "vi"},
	"techno":  {"i", "iv", "v", "i"},
	"ambient": {"I", "vi", "IV", "V"},
}

// MapRomanToDegrees maps Roman numerals to scale degrees (0-indexed).
func MapRomanToDegrees(roman string) int {
	switch strings.ToLower(roman) {
	case "i": return 0
	case "ii": return 1
	case "iii": return 2
	case "iv": return 3
	case "v": return 4
	case "vi": return 5
	case "vii": return 6
	default: return 0
	}
}

// GenerateProgression generates a chord progression based on a style and key.
func GenerateProgression(style string, root int, scaleName string) [][]int {
	pattern, ok := RomanNumeralPatterns[strings.ToLower(style)]
	if !ok {
		pattern = RomanNumeralPatterns["basic"]
	}

	progression := [][]int{}

	for _, roman := range pattern {
		degree := MapRomanToDegrees(roman)
		chordRoot := (root + GetIntervalForDegree(degree, scaleName))

		// Simple Triad (Root, Third, Fifth)
		chord := []int{
			chordRoot,
			chordRoot + GetThird(degree, scaleName),
			chordRoot + 7, // Perfect fifth heuristic
		}
		progression = append(progression, chord)
	}

	return progression
}

func GetIntervalForDegree(degree int, scaleName string) int {
	scale, ok := Scales[strings.ToLower(scaleName)]
	if !ok { return 0 }
	if degree >= len(scale.Intervals) { return 0 }
	return scale.Intervals[degree]
}

func GetThird(degree int, scaleName string) int {
	// Simple major/minor third detection based on scale degrees
	// This is a simplified version of scribbletune logic
	scale, _ := Scales[strings.ToLower(scaleName)]
	first := scale.Intervals[degree]
	thirdIdx := (degree + 2) % len(scale.Intervals)
	third := scale.Intervals[thirdIdx]

	diff := third - first
	if diff < 0 { diff += 12 }
	return diff
}
