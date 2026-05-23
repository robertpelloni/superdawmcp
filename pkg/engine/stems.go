package engine

import (
	"fmt"
	"os/exec"
)

// SeparateStems invokes the Spleeter CLI to separate an audio file into stems.
func SeparateStems(inputPath, outputDir string, stems int) (string, error) {
	// Ensure Spleeter is in PATH.
	// Command: spleeter separate -p spleeter:2stems -o output_dir input_path

	params := fmt.Sprintf("spleeter:%dstems", stems)
	if stems != 2 && stems != 4 && stems != 5 {
		params = "spleeter:4stems" // Fallback to 4 stems
	}

	cmd := exec.Command("spleeter", "separate", "-p", params, "-o", outputDir, inputPath)

	// We run the command but don't block indefinitely in the MCP loop.
	// In a real scenario, this would be a background task with status reporting.
	err := cmd.Start()
	if err != nil {
		return "", fmt.Errorf("failed to start spleeter: %w. ensures spleeter is installed: pip install spleeter", err)
	}

	return fmt.Sprintf("Started Spleeter separation for %s into %s (stems: %d). PID: %d", inputPath, outputDir, stems, cmd.Process.Pid), nil
}
