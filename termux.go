package termux

import (
	"os/exec"
	"strings"
)

// IsInstalled returns wether or not termux commands are available
func IsInstalled() bool {
	_, err := exec.LookPath("termux-volume")
	return err == nil
}

func toEnumStr(data []byte) string {
	return strings.Trim(strings.ToLower(string(data)), "\"")
}
