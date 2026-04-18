package watch

import (
	"testing"
	"time"
)

func sendEvents(ch chan<- Event, secrets []map[string]string, delay time.Duration) {
	for _, s := range secrets {
		ch <- Event{Path: "secret/app", Secrets: s, At: time.Now()}
		time.Sleep(delay)
	}
	close(ch)
}

func TestDebounce_CoalescesRapidEvents(t *testing.T) {
	in := make(chan Event, 10)
	out := make(chan Event, 10)

	d := NewDebouncer(40 * time.Millisecond)
	go d.Run(in, out)

	sendEvents(in, []map[string]string{
		{"K": "1"},
		{"K": "2"},
		{"K": "3"},
	}, 5*time.Millisecond)

	var received []Event
	timeout := time.After(200 * time.Millisecond)
	for {
		select {
		case e, ok := <-out:
			if !ok {
				goto done
			}
			received = append(received, e)
		case <-timeout:
			goto done
		}
	}
done:
	if len(received) != 1 {
		t.Errorf("expected 1 debounced event, got %d", len(received))
	}
	if received[0].Secrets["K"] != "3" {
		t.Errorf("expected last value '3', got %s", received[0].Secrets["K"])
	}
}

func TestDebounce_ForwardsAfterQuiet(t *testing.T) {
	in := make(chan Event, 4)
	out := make(chan Event, 4)

	d := NewDebouncer(20 * time.Millisecond)
	go d.Run(in, out)

	in <- Event{Path: "p", Secrets: map[string]string{"X": "a"}, At: time.Now()}
	time.Sleep(60 * time.Millisecond)
	in <- Event{Path: "p", Secrets: map[string]string{"X": "b"}, At: time.Now()}
	close(in)

	time.Sleep(60 * time.Millisecond)
	if len(out) != 2 {
		t.Errorf("expected 2 events (one per quiet period), got %d", len(out))
	}
}
