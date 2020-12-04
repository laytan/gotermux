package termux

import (
	"encoding/json"
	"os/exec"
)

// LocationInfo is the information about the device
type LocationInfo struct {
	Latitude         float64
	Longitude        float64
	Altitude         float64
	Accuracy         float64
	VerticalAccuracy float64
	Bearing          float64
	Speed            float64
	// Time on which the speed is calculated
	ElapsedMS uint
	// gps, wifi etc
	Provider string
}

// Location returns the device's location info (requires permission)
func Location() (LocationInfo, error) {
	info := LocationInfo{}

	cmd := exec.Command("termux-location")
	bytes, err := cmd.Output()
	if err != nil {
		return info, err
	}

	err = json.Unmarshal(bytes, &info)
	return info, err
}
