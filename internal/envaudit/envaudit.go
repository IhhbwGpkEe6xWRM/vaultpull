// Package envaudit provides change-tracking for environment variable maps.
// It records when keys are added, removed, or modified and produces a
// structured audit trail suitable for logging or compliance reporting.
package envaudit

import (
	"fmt"
	"sort"
	"time"
)

// EventKind describes the type of change that occurred to a key.
type EventKind string

const (
	EventAdded    EventKind = "added"
	EventRemoved  EventKind = "removed"
	EventModified EventKind = "modified"
)

// Event represents a single change to an environment variable.
type Event struct {
	Kind      EventKind `json:"kind"`
	Key       string    `json:"key"`
	OldValue  string    `json:"old_value,omitempty"`
	NewValue  string    `json:"new_value,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}

// Tracker records audit events for changes between two env maps.
type Tracker struct {
	events []Event
	clock  func() time.Time
}

// Option configures a Tracker.
type Option func(*Tracker)

// WithClock overrides the clock used to timestamp events.
func WithClock(fn func() time.Time) Option {
	return func(t *Tracker) {
		t.clock = fn
	}
}

// New returns a new Tracker with optional configuration.
func New(opts ...Option) *Tracker {
	t := &Tracker{
		clock: time.Now,
	}
	for _, o := range opts {
		o(t)
	}
	return t
}

// Record compares prev and next, appending audit events for every
// key that was added, removed, or whose value changed. Values are
// stored verbatim; callers should redact sensitive data before calling.
func (t *Tracker) Record(prev, next map[string]string) {
	now := t.clock()

	for key, newVal := range next {
		if oldVal, exists := prev[key]; !exists {
			t.events = append(t.events, Event{
				Kind:      EventAdded,
				Key:       key,
				NewValue:  newVal,
				Timestamp: now,
			})
		} else if oldVal != newVal {
			t.events = append(t.events, Event{
				Kind:      EventModified,
				Key:       key,
				OldValue:  oldVal,
				NewValue:  newVal,
				Timestamp: now,
			})
		}
	}

	for key, oldVal := range prev {
		if _, exists := next[key]; !exists {
			t.events = append(t.events, Event{
				Kind:      EventRemoved,
				Key:       key,
				OldValue:  oldVal,
				Timestamp: now,
			})
		}
	}
}

// Events returns a sorted copy of all recorded events (sorted by key,
// then by timestamp for deterministic output).
func (t *Tracker) Events() []Event {
	out := make([]Event, len(t.events))
	copy(out, t.events)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Key != out[j].Key {
			return out[i].Key < out[j].Key
		}
		return out[i].Timestamp.Before(out[j].Timestamp)
	})
	return out
}

// Reset clears all recorded events.
func (t *Tracker) Reset() {
	t.events = t.events[:0]
}

// Summary returns a human-readable one-line summary of recorded changes.
func (t *Tracker) Summary() string {
	var added, removed, modified int
	for _, e := range t.events {
		switch e.Kind {
		case EventAdded:
			added++
		case EventRemoved:
			removed++
		case EventModified:
			modified++
		}
	}
	if added+removed+modified == 0 {
		return "no changes"
	}
	return fmt.Sprintf("%d added, %d removed, %d modified", added, removed, modified)
}
