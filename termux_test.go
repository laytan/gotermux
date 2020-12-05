package termux

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"
)

type testBatteryPayload struct {
	input        string
	shouldOutput BatteryInfo
}

func TestBattery(t *testing.T) {
	payloads := []testBatteryPayload{
		{
			input: `{"health": "GOOD", "percentage": 37, "plugged": "UNPLUGGED", "status": "DISCHARGING", "temperature": 23.23, "current": -222}`,
			shouldOutput: BatteryInfo{
				Health:      BatteryHealthGood,
				Percentage:  37,
				Plugged:     BatteryPluggedUnplugged,
				Status:      BatteryStatusDischarging,
				Temperature: 23.23,
				Current:     -222,
			},
		},
	}

	for _, payload := range payloads {
		out := BatteryInfo{}
		if err := json.Unmarshal([]byte(payload.input), &out); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(out, payload.shouldOutput) {
			t.Error("Out and shouldOutput not the same")
		}
	}
}

type testLocationPayload struct {
	input        string
	shouldOutput LocationInfo
}

func TestLocation(t *testing.T) {
	payloads := []testLocationPayload{
		{
			input: `{"latitude": 22.123123123, "longitude": 22.123123, "altitude": 22.123123, "accuracy": 15.0, "vertical_accuracy": 16.111111, "bearing": 0.0, "speed": 0.0, "elapsedMs": 30, "provider": "network"}`,
			shouldOutput: LocationInfo{
				Latitude:         22.123123123,
				Longitude:        22.123123,
				Altitude:         22.123123,
				Accuracy:         15.0,
				VerticalAccuracy: 16.111111,
				Bearing:          0.0,
				Speed:            0.0,
				ElapsedMS:        30,
				Provider:         LocationProviderNetwork,
			},
		},
	}

	for _, payload := range payloads {
		out := LocationInfo{}
		if err := json.Unmarshal([]byte(payload.input), &out); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(out, payload.shouldOutput) {
			t.Error("Out and shouldOutput not the same")
		}
	}
}

type testVolumePayload struct {
	input        string
	shouldOutput map[VolumeStream]Volume
}

func TestVolumeInfo(t *testing.T) {
	payloads := []testVolumePayload{
		{
			input: `[{"stream":"call","volume":4,"max_volume":5},{"stream":"system","volume":2,"max_volume":15},{"stream":"ring","volume":1,"max_volume":15},{"stream":"music","volume":7,"max_volume":15},{"stream":"alarm","volume":11,"max_volume":15},{"stream":"notification","volume":2,"max_volume":15}]`,
			shouldOutput: map[VolumeStream]Volume{
				VolumeStreamCall: Volume{
					MaxVolume: 5,
					Volume:    4,
				},
				VolumeStreamSystem: Volume{
					MaxVolume: 15,
					Volume:    2,
				},
				VolumeStreamRing: Volume{
					MaxVolume: 15,
					Volume:    1,
				},
				VolumeStreamMusic: Volume{
					MaxVolume: 15,
					Volume:    7,
				},
				VolumeStreamAlarm: Volume{
					MaxVolume: 15,
					Volume:    11,
				},
				VolumeStreamNotification: Volume{
					MaxVolume: 15,
					Volume:    2,
				},
			},
		},
	}

	for _, payload := range payloads {
		streams := make(map[VolumeStream]Volume)
		res := make([]termuxVolumeRes, 0)
		if err := json.Unmarshal([]byte(payload.input), &res); err != nil {
			t.Error(err)
		}

		for _, stream := range res {
			streams[stream.Stream] = Volume{Volume: stream.Volume, MaxVolume: stream.MaxVolume}
		}

		if !reflect.DeepEqual(streams, payload.shouldOutput) {
			t.Error("Out and shouldOutput not the same")
		}
	}
}

type testNotificationOptionsPayload struct {
	input               NotificationOptions
	shouldOutput        string
	shouldNotBeInOutput []string
}

func TestNotificationOptions(t *testing.T) {
	payloads := []testNotificationOptionsPayload{
		{
			input:        NotificationOptions{},
			shouldOutput: "",
		},
		// Everything present, true and correct
		{
			input: NotificationOptions{
				Title:     "Test",
				Content:   "Testing",
				Action:    "Actie!",
				AlertOnce: true,
				Buttons: []NotificationButton{
					{
						Text:   "btn1",
						Action: "Action1",
					},
					{
						Text:   "btn2",
						Action: "Action2",
					},
					{
						Text:   "btn3",
						Action: "Action3",
					},
				},
				Group:     "test",
				ID:        5,
				ImagePath: "testpath",
				LED: &NotificationLED{
					Color: "123123",
					Off:   time.Duration(time.Millisecond * 100),
					On:    time.Duration(time.Millisecond * 500),
				},
				OnDelete: "testactie",
				OnGoing:  true,
				Priority: NotificationPriorityHigh,
				Sound:    true,
				VibratePattern: []time.Duration{
					time.Duration(time.Millisecond),
					time.Duration(time.Second),
				},
				Type: NotificationTypeMedia,
			},
			shouldOutput:        "-t Test -c Testing --action Actie! --alert-once --button1 btn1 --button1-action Action1 --button2 btn2 --button2-action Action2 --button3 btn3 --button3-action Action3 --group test -i 5 --image-path testpath --led-color 123123 --led-off 100 --led-on 500 --on-delete testactie --ongoing --sound --vibrate 1,1000 --type media --priority high",
			shouldNotBeInOutput: []string{},
		},
		// TODO: Test > 3 buttons
		// TODO: Test false booleans not included
	}

	for _, payload := range payloads {
		fmt.Println(payload.input.String())
		for _, part := range strings.Split(payload.shouldOutput, " ") {
			if !strings.Contains(payload.input.String(), part) {
				t.Errorf("%s not in input", part)
			}
		}

		for _, notbe := range payload.shouldNotBeInOutput {
			if strings.Contains(payload.input.String(), notbe) {
				t.Errorf("%s found in input", notbe)
			}
		}
	}
}
