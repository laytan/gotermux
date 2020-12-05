package termux

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// LocationProvider is the provider of a location
type LocationProvider uint

func (l LocationProvider) String() string {
	return []string{"gps", "network", "passive"}[l]
}

// UnmarshalJSON turns the json data into one of the enum types
func (l *LocationProvider) UnmarshalJSON(data []byte) error {
	for i, provider := range []string{"gps", "network", "passive"} {
		if toEnumStr(data) == provider {
			*l = LocationProvider(i)
			return nil
		}
	}
	return fmt.Errorf("No provider enum type for %s", string(data))
}

// All possible location providers
const (
	LocationProviderGPS LocationProvider = iota
	LocationProviderNetwork
	LocationProviderPassive
)

// LocationInfo is the information about the device
type LocationInfo struct {
	Latitude         float64
	Longitude        float64
	Altitude         float64
	Accuracy         float64
	VerticalAccuracy float64 `json:"vertical_accuracy"`
	Bearing          float64
	Speed            float64
	// Time on which the speed is calculated
	ElapsedMS uint
	Provider  LocationProvider
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
