package daw

import (
	"fmt"
	"os/exec"
)

// JackRouter utilizes the 'jack_connect' and 'jack_disconnect' CLI tools
// to manage audio routing on systems with JACK installed.
type JackRouter struct{}

func NewJackRouter() *JackRouter {
	return &JackRouter{}
}

func (j *JackRouter) Patch(source, dest string) error {
	// Example: jack_connect "Ableton:out1" "REAPER:in1"
	cmd := exec.Command("jack_connect", source, dest)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("jack_connect failed: %w", err)
	}
	return nil
}

func (j *JackRouter) Unpatch(source, dest string) error {
	cmd := exec.Command("jack_disconnect", source, dest)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("jack_disconnect failed: %w", err)
	}
	return nil
}
