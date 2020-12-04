package termux

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

// VolumeStream is a specific type of volume
type VolumeStream uint

func (v VolumeStream) String() string {
	return []string{"call", "system", "ring", "music", "alarm", "notification"}[v]
}

// UnmarshalJSON turns json into our enum type
func (v *VolumeStream) UnmarshalJSON(data []byte) error {
	for i, stream := range []string{"call", "system", "ring", "music", "alarm", "notification"} {
		if string(data) == stream {
			*v = VolumeStream(i)
		}
	}
	return fmt.Errorf("No value found for %s as VolumeStream", string(data))
}

// All available volume streams
const (
	VolumeStreamCall VolumeStream = iota
	VolumeStreamSystem
	VolumeStreamRing
	VolumeStreamMusic
	VolumeStreamAlarm
	VolumeStreamNotification
)

// Volume has info off a specific stream
type Volume struct {
	// call, system, ring, music, alarm, notification
	Volume    uint8
	MaxVolume uint8
}

type termuxVolumeRes struct {
	Stream    VolumeStream
	Volume    uint8
	MaxVolume uint8
}

// VolumeInfo returns the devices volume streams
func VolumeInfo() (map[VolumeStream]Volume, error) {
	streams := make(map[VolumeStream]Volume)

	cmd := exec.Command("termux-volume")
	output, err := cmd.Output()
	if err != nil {
		return streams, err
	}

	res := make([]termuxVolumeRes, 0)
	if err := json.Unmarshal(output, &res); err != nil {
		return streams, err
	}

	for _, stream := range res {
		streams[stream.Stream] = Volume{Volume: stream.Volume, MaxVolume: stream.MaxVolume}
	}

	return streams, nil
}

// VolumeOf returns the volume info of a specific stream
func VolumeOf(stream VolumeStream) (Volume, error) {
	info, err := VolumeInfo()
	if err != nil {
		return Volume{}, err
	}

	vol, exists := info[stream]
	if !exists {
		return Volume{}, errors.New("stream does not exist")
	}

	return vol, nil
}

// SetVolume sets the volume of the given stream
func SetVolume(stream VolumeStream, vol uint8) error {
	cmd := exec.Command("termux-volume", stream.String(), fmt.Sprintf("%d", vol))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
