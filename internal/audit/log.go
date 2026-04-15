// Package audit provides structured audit logging for vaultpull operations.
// It records which secrets were synced, skipped, or failed during a run.
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// EventType represents the kind of audit event.
type EventType string

const (
	EventSynced  EventType = "synced"
	EventSkipped EventType = "skipped"
	EventFailed  EventType = "failed"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	Path      string    `json:"path"`
	Message   string    `json:"message,omitempty"`
}

// Logger writes structured audit entries to a destination.
type Logger struct {
	w io.Writer
}

// NewLogger creates a Logger writing to the given path.
// Pass an empty string to disable file logging (no-op logger).
func NewLogger(path string) (*Logger, error) {
	if path == "" {
		return &Logger{w: io.Discard}, nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return nil, fmt.Errorf("audit: open log file: %w", err)
	}
	return &Logger{w: f}, nil
}

// NewLoggerWithWriter creates a Logger writing to w (useful for testing).
func NewLoggerWithWriter(w io.Writer) *Logger {
	return &Logger{w: w}
}

// Record writes an audit entry for the given event and path.
func (l *Logger) Record(event EventType, path, message string) error {
	entry := Entry{
		Timestamp: time.Now().UTC(),
		Event:     event,
		Path:      path,
		Message:   message,
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	_, err = fmt.Fprintf(l.w, "%s\n", data)
	return err
}
