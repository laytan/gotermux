package termux

import (
	"fmt"
	"os/exec"
	"strings"
)

// ToastPosition specifies where to show a toast
type ToastPosition uint

func (p ToastPosition) String() string {
	return []string{"middle", "top", "bottom"}[p]
}

// The possible toast positions
const (
	ToastMiddle ToastPosition = iota
	ToastTop
	ToastBottom
)

// ToastOptions are the options for a toast
// Note: Colors can be default color names or a 6/8 character hex color
// (i.e. "#FF0000" or "#FFFF0000") where order is (AA)RRGGBB. Invalid color will revert to default value.
type ToastOptions struct {
	// Default: gray
	BackgroundColor string
	// Default: white
	TextColor string
	// top, middle or bottom. Default: middle
	Position ToastPosition
	Short    bool
}

// ShowToast shows a toast(message) on the screen
func ShowToast(msg string, opts *ToastOptions) error {
	var optsStr string
	if opts != nil {
		if len(opts.BackgroundColor) > 0 {
			optsStr += fmt.Sprintf(" -b %s", opts.BackgroundColor)
		}
		if len(opts.TextColor) > 0 {
			optsStr += fmt.Sprintf(" -c %s", opts.TextColor)
		}
		optsStr += fmt.Sprintf(" -g %s", opts.Position)
		if opts.Short {
			optsStr += " -s"
		}
	}

	cmd := exec.Command("termux-toast", msg, strings.Trim(optsStr, " "))
	return cmd.Run()
}
