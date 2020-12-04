package termux

import (
	"fmt"
	"os/exec"
)

// SetBrightness sets the screen brightness
// Values can be 0 to 255
func SetBrightness(brightness uint8) error {
	cmd := exec.Command("termux-brightness", fmt.Sprintf("%d", brightness))
	return cmd.Run()
}
