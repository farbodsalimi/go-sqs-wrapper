package util

import "os/exec"

// IsCommandAvailable checks if a command is available or not
func IsCommandAvailable(name string) bool {
	cmd := exec.Command("command", "-v", name)
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}
