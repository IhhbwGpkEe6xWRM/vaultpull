package watch

import (
	"context"
	"errors"
	"testing"
	"time"
)

type mockReader struct {
	calls   int
	results []map[string]string
	err     error
}

func (m *mockReader) ReadSecrets(_ context.Context, _ string) (map[string]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	idx := m.calls
	if idx >= len(m.results) {
		idx = len(m.results) - 1
	}
	m.calls++
	return m.results[idx], nil
}

func fastClock() func() time.Time {
	t := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	return func() time.Time { return t }
}

func TestWatch_EmitsEventOnChange(t *testing.T) {
	reader := &mockReader{
		results: []map[string]string{
			{"KEY": "v1"},
			{"KEY": "v2"},
		},
	}
	cfg := Config{Path: "secret/app", Interval: 10 * time.Millisecond}
	w := newWithClock(cfg, reader, fastClock())

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	ch := make(chan Event, 4)
	go w.Watch(ctx, ch) //nolint:errcheck

	var events []Event
	timer := time.NewTimer(150 * time.Millisecond)
	defer timer.Stop()
	for {
		select {
		case e := <-ch:
			events = append(events, e)
		case <-timer.C:
			goto done
		}
	}
done:
	if len(events) == 0 {
		t.Fatal("expected at least one event")
	}
	if events[0].Path != "secret/app" {
		t.Errorf("unexpected path: %s", events[0].Path)
	}
}

func TestWatch_NoEventWhenUnchanged(t *testing.T) {
	reader := &mockReader{
		results: []map[string]string{{"KEY": "same"}},
	}
	cfg := Config{Path: "secret/app", Interval: 10 * time.Millisecond}
	w := newWithClock(cfg, reader, fastClock())

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	ch := make(chan Event, 4)
	go w.Watch(ctx, ch) //nolint:errcheck

	time.Sleep(70 * time.Millisecond)
	if len(ch) > 1 {
		t.Errorf("expected at most 1 event (initial), got %d", len(ch))
	}
}

func TestWatch_SkipsOnReaderError(t *testing.T) {
	reader := &mockReader{err: errors.New("vault down")}
	cfg := Config{Path: "secret/app", Interval: 10 * time.Millisecond}
	w := newWithClock(cfg, reader, fastClock())

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()

	ch := make(chan Event, 4)
	go w.Watch(ctx, ch) //nolint:errcheck

	time.Sleep(50 * time.Millisecond)
	if len(ch) != 0 {
		t.Errorf("expected no events on error, got %d", len(ch))
	}
}

func TestChanged_DetectsAddedKey(t *testing.T) {
	prev := map[string]string{"A": "1"}
	curr := map[string]string{"A": "1", "B": "2"}
	if !changed(prev, curr) {
		t.Error("expected changed to be true")
	}
}

func TestChanged_IdenticalMaps(t *testing.T) {
	m := map[string]string{"A": "1"}
	if changed(m, m) {
		t.Error("expected changed to be false")
	}
}
