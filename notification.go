package termux

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
)

// NotificationAction specifies a command to execute
type NotificationAction string

// NotificationPriority is how the notification is handled by the os high/low/max/min/default
type NotificationPriority uint

func (n NotificationPriority) String() string {
	return []string{"default", "high", "low", "max", "min"}[n]
}

// Available notification priorities
const (
	NotificationPriorityDefault NotificationPriority = iota
	NotificationPriorityHigh
	NotificationPriorityLow
	NotificationPriorityMax
	NotificationPriorityMin
)

// NotificationType is the type of notification, media or default
type NotificationType uint

func (n NotificationType) String() string {
	return []string{"default", "media"}[n]
}

// Available notification types
const (
	NotificationTypeDefault NotificationType = iota
	NotificationTypeMedia
)

// NotificationLED sets led options
type NotificationLED struct {
	// Set an LED color, should be rrggbb without #
	Color string `json:"--led-color,omitempty"`
	// How long the led should be off
	Off time.Duration `json:"--led-off,omitempty"`
	// How long the led should be on
	On time.Duration `json:"--led-on,omitempty"`
}

// NotificationOptions are the options to use when showing a notification
// More info: https://wiki.termux.com/wiki/Termux-notification
type NotificationOptions struct {
	// Notification title
	Title string `json:"-t,omitempty"`
	// Notification content
	Content string `json:"-c,omitempty"`
	// On clicking the notification
	Action NotificationAction `json:"--action,omitempty"`
	// Only alert on initial send, not on edits etc
	AlertOnce bool `json:"--alert-once"`
	// Buttons, there can be 3 buttons in a notification
	Buttons []NotificationButton `json:"buttons,omitempty"`
	// Notifications with the same group will be grouped together
	Group string `json:"--group,omitempty"`
	// Override an existing notification by id
	ID uint `json:"-i,omitempty"`
	// Show an image, should be an absolute path
	ImagePath string `json:"--image-path,omitempty"`
	// Led options
	LED *NotificationLED `json:"led,omitempty"`
	// Action to execute on deleting/clearing the notification
	OnDelete NotificationAction `json:"--on-delete,omitempty"`
	// Pin the notification
	OnGoing bool `json:"--ongoing"`
	// What priority is the notification
	Priority NotificationPriority `json:"--priority,omitempty"`
	// Play a sound with the notification
	Sound bool `json:"--sound"`
	// Specify a pattern to vibrate on
	VibratePattern []time.Duration `json:"vibratePattern,omitempty"`
	// Type default or media notification
	Type NotificationType `json:"--type,omitempty"`
}

func (o *NotificationOptions) String() string {
	// Convert the struct to json
	js, _ := json.Marshal(o)
	// Convert the json to a map
	asMap := make(map[string]interface{})
	json.Unmarshal(js, &asMap)

	var opts string

	for k, v := range asMap {
		switch k {
		// Boolean case, add the option if the value is true
		case "--alert-once", "--ongoing", "--sound":
			if v.(bool) == true {
				opts += " " + k
			}
		case "buttons":
			asSlice := v.([]interface{})

			// Max of 3 buttons
			if len(asSlice) > 3 {
				asSlice = asSlice[:2]
			}

			for i, btn := range asSlice {
				asMap := btn.(map[string]interface{})

				text, exists := asMap["Text"]
				action, existsA := asMap["Action"]
				if !exists || !existsA {
					log.Println("Button without text or without action is not allowed")
					break
				}

				opts += fmt.Sprintf(" --button%d %s --button%d-action %s", i+1, text, i+1, action)
			}
		case "led":
			for mk, mv := range v.(map[string]interface{}) {
				if mk == "--led-color" {
					opts += fmt.Sprintf(" %s %s", mk, mv)
				} else {
					// Turn the time.Duration nanoseconds into milliseconds
					opts += fmt.Sprintf(" %s %d", mk, uint(mv.(float64)/1000000))
				}
			}
		case "vibratePattern":
			opts += " --vibrate "
			asSlice := v.([]interface{})
			opts += fmt.Sprintf("%d", uint(asSlice[0].(float64)/1000000))
			for _, t := range asSlice[1:] {
				opts += fmt.Sprintf(",%d", uint(t.(float64)/1000000))
			}
			opts += " "
		case "--priority":
			var prio NotificationPriority
			json.Unmarshal([]byte(fmt.Sprintf("%v", v)), &prio)
			opts += fmt.Sprintf(" %s %s", k, prio.String())
		case "--type":
			var t NotificationType
			json.Unmarshal([]byte(fmt.Sprintf("%v", v)), &t)
			opts += fmt.Sprintf(" %s %s", k, t.String())
		default:
			opts += fmt.Sprintf(" %s %v", k, v)
		}
	}

	return opts
}

// NotificationButton shown in a notification
type NotificationButton struct {
	// Command to execute on clicking the button
	Action NotificationAction
	// Button text
	Text string
}

// ShowNotification sends a notification
func ShowNotification(title string, content string, opts *NotificationOptions) error {
	if opts == nil {
		opts = &NotificationOptions{
			Title:   title,
			Content: content,
		}
	}

	cmd := exec.Command("termux-notification", opts.String())
	return cmd.Run()
}

// RemoveNotification removes a notification that is currently shown
func RemoveNotification(id uint) error {
	cmd := exec.Command("termux-notification-remove", fmt.Sprintf("%d", id))
	return cmd.Run()
}
