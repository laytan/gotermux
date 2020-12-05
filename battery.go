package termux

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// BatteryHealth is a status of the battery in the device
type BatteryHealth uint

func (b BatteryHealth) String() string {
	return []string{"unknown", "unspecified_failure", "good", "cold", "dead", "overheat", "over_voltage"}[b]
}

// UnmarshalJSON turns the json string into our enum
func (b *BatteryHealth) UnmarshalJSON(data []byte) error {
	for i, healthStatus := range []string{"unknown", "unspecified_failure", "good", "cold", "dead", "overheat", "over_voltage"} {
		if toEnumStr(data) == healthStatus {
			*b = BatteryHealth(i)
			return nil
		}
	}
	return fmt.Errorf("No BatteryHealth value for %s", string(data))
}

// All possible BatteryHealth statusses
const (
	BatteryHealthUnknown BatteryHealth = iota
	BatteryHealthUnspecifiedFailure
	BatteryHealthGood
	BatteryHealthCold
	BatteryHealthDead
	BatteryHealthOverHeat
	BatteryHealthOverVoltage
)

// BatteryPlugged is the plugged state of the device
type BatteryPlugged uint

func (b BatteryPlugged) String() string {
	return []string{"unplugged", "plugged_ac", "plugged_usb", "plugged_wireless"}[b]
}

// UnmarshalJSON turns the json into our enum type
func (b *BatteryPlugged) UnmarshalJSON(data []byte) error {
	for i, pluggedStatus := range []string{"unplugged", "plugged_ac", "plugged_usb", "plugged_wireless"} {
		if toEnumStr(data) == pluggedStatus {
			*b = BatteryPlugged(i)
			return nil
		}
	}
	return fmt.Errorf("No BatteryPlugged value for %s", string(data))
}

// All possible BatteryPlugged states
const (
	BatteryPluggedUnplugged BatteryPlugged = iota
	BatteryPluggedAC
	BatteryPluggedUSB
	BatteryPluggedWireless
)

// BatteryStatus is the status of the battery of the device
type BatteryStatus uint

func (b BatteryStatus) String() string {
	return []string{"unknown", "charging", "discharging", "full", "not_charging"}[b]
}

// UnmarshalJSON turns json into our enum type
func (b *BatteryStatus) UnmarshalJSON(data []byte) error {
	for i, batteryStatus := range []string{"unknown", "charging", "discharging", "full", "not_charging"} {
		if toEnumStr(data) == batteryStatus {
			*b = BatteryStatus(i)
			return nil
		}
	}
	return fmt.Errorf("No BatteryStatus value for %s", string(data))
}

// All possible battery statusses
const (
	BatteryStatusUnknown BatteryStatus = iota
	BatteryStatusCharging
	BatteryStatusDischarging
	BatteryStatusFull
	BatteryStatusNotCharging
)

// BatteryInfo is info about the battery
// See: https://developer.android.com/reference/android/os/BatteryManager
// See: https://wiki.termux.com/wiki/Termux-battery-status
type BatteryInfo struct {
	Health BatteryHealth
	// 0 to 100
	Percentage uint8
	Plugged    BatteryPlugged
	Status     BatteryStatus
	// Battery temperature
	Temperature float64
	// Current coming in or out of the device, depending on usage in microamperes
	Current int
}

// Battery returns information about the current battery status
func Battery() (BatteryInfo, error) {
	info := BatteryInfo{}

	cmd := exec.Command("termux-battery-status")
	bytes, err := cmd.Output()
	if err != nil {
		return info, err
	}

	err = json.Unmarshal(bytes, &info)
	return info, err
}
