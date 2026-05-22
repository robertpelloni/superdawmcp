package engine

import (
	"fmt"
	"os/exec"
)

// SeparateStems invokes the Spleeter CLI to separate an audio file into stems.
func SeparateStems(inputPath string, outputDir string, stems int) error {
	// Example command: spleeter separate -p spleeter:4stems -o output/ input.mp3
	stemsStr := fmt.Sprintf("spleeter:%dstems", stems)
	cmd := exec.Command("spleeter", "separate", "-p", stemsStr, "-o", outputDir, inputPath)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run spleeter: %w", err)
	}

	return nil
}
