package termux

import (
	"os/exec"
)

// IsInstalled returns wether or not termux commands are available
func IsInstalled() bool {
	_, err := exec.LookPath("termux-volume")
	return err == nil
}
