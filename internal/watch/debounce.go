package watch

import (
	"time"
)

// Debouncer coalesces rapid Events into a single delivery after a quiet period.
type Debouncer struct {
	delay time.Duration
}

// NewDebouncer returns a Debouncer that waits delay after the last event
// before forwarding to out.
func NewDebouncer(delay time.Duration) *Debouncer {
	return &Debouncer{delay: delay}
}

// Run reads from in and writes debounced events to out.
// It returns when in is closed or the context signals done via the caller
// closing in.
func (d *Debouncer) Run(in <-chan Event, out chan<- Event) {
	var (
		pending *Event
		timer   *time.Timer
	)

	flush := func() {
		if pending != nil {
			out <- *pending
			pending = nil
		}
	}

	for {
		if timer == nil {
			ev, ok := <-in
			if !ok {
				flush()
				return
			}
			pending = &ev
			timer = time.NewTimer(d.delay)
			continue
		}

		select {
		case ev, ok := <-in:
			if !ok {
				flush()
				return
			}
			pending = &ev
			timer.Reset(d.delay)
		case <-timer.C:
			flush()
			timer = nil
		}
	}
}
