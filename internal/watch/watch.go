package watch

import (
	"context"
	"time"
)

// Watcher polls a secret path at a fixed interval and notifies via channel
// when secrets change compared to the previous read.
type Watcher struct {
	cfg    Config
	reader SecretReader
	clock  func() time.Time
}

// SecretReader reads secrets from a path.
type SecretReader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Config holds watcher configuration.
type Config struct {
	Path     string
	Interval time.Duration
}

// Event carries changed secrets and the time of detection.
type Event struct {
	Path    string
	Secrets map[string]string
	At      time.Time
}

// New returns a Watcher with default clock.
func New(cfg Config, r SecretReader) *Watcher {
	return &Watcher{cfg: cfg, reader: r, clock: time.Now}
}

// newWithClock returns a Watcher with an injectable clock for testing.
func newWithClock(cfg Config, r SecretReader, clock func() time.Time) *Watcher {
	return &Watcher{cfg: cfg, reader: r, clock: clock}
}

// Watch polls until ctx is cancelled, sending an Event on ch whenever secrets
// differ from the previously seen snapshot.
func (w *Watcher) Watch(ctx context.Context, ch chan<- Event) error {
	var prev map[string]string

	tick := time.NewTicker(w.cfg.Interval)
	defer tick.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			current, err := w.reader.ReadSecrets(ctx, w.cfg.Path)
			if err != nil {
				continue
			}
			if changed(prev, current) {
				prev = current
				ch <- Event{Path: w.cfg.Path, Secrets: current, At: w.clock()}
			}
		}
	}
}

func changed(prev, current map[string]string) bool {
	if len(prev) != len(current) {
		return true
	}
	for k, v := range current {
		if prev[k] != v {
			return true
		}
	}
	return false
}
