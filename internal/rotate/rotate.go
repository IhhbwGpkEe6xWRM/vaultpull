// Package rotate provides utilities for detecting and reporting
// secret rotation events by comparing current secrets against snapshots.
package rotate

import (
	"fmt"
	"time"
)

// Event describes a single secret rotation occurrence.
type Event struct {
	Key       string    `json:"key"`
	Path      string    `json:"path"`
	DetectedAt time.Time `json:"detected_at"`
}

// Detector compares two secret maps and returns rotation events
// for keys whose values have changed.
type Detector struct {
	path string
}

// NewDetector creates a Detector scoped to the given secret path.
func NewDetector(path string) *Detector {
	return &Detector{path: path}
}

// Detect returns an Event for each key whose value differs between
// previous and current. Keys only in current are not considered rotations.
func (d *Detector) Detect(previous, current map[string]string) []Event {
	var events []Event
	for k, newVal := range current {
		oldVal, existed := previous[k]
		if existed && oldVal != newVal {
			events = append(events, Event{
				Key:        k,
				Path:       d.path,
				DetectedAt: time.Now().UTC(),
			})
		}
	}
	return events
}

// Summary returns a human-readable summary line for a slice of events.
func Summary(events []Event) string {
	if len(events) == 0 {
		return "no rotated secrets detected"
	}
	return fmt.Sprintf("%d secret(s) rotated", len(events))
}

// Keys returns a slice of the key names from the given events.
// This is useful for logging or filtering which secrets were rotated.
func Keys(events []Event) []string {
	keys := make([]string, len(events))
	for i, e := range events {
		keys[i] = e.Key
	}
	return keys
}
