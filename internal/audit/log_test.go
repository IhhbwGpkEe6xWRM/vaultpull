package audit_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/audit"
)

func TestRecord_WritesJSONLine(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLoggerWithWriter(&buf)

	if err := l.Record(audit.EventSynced, "secret/app/db", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var entry audit.Entry
	if err := json.Unmarshal(bytes.TrimSpace(buf.Bytes()), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if entry.Event != audit.EventSynced {
		t.Errorf("expected event %q, got %q", audit.EventSynced, entry.Event)
	}
	if entry.Path != "secret/app/db" {
		t.Errorf("expected path %q, got %q", "secret/app/db", entry.Path)
	}
}

func TestRecord_IncludesMessage(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLoggerWithWriter(&buf)

	_ = l.Record(audit.EventFailed, "secret/app/key", "permission denied")

	if !strings.Contains(buf.String(), "permission denied") {
		t.Error("expected message in output")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	var buf bytes.Buffer
	l := audit.NewLoggerWithWriter(&buf)

	_ = l.Record(audit.EventSynced, "secret/a", "")
	_ = l.Record(audit.EventSkipped, "secret/b", "namespace mismatch")

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 lines, got %d", len(lines))
	}
}

func TestNewLogger_DisabledWhenEmpty(t *testing.T) {
	l, err := audit.NewLogger("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Should not panic or error — writes to discard.
	if err := l.Record(audit.EventSynced, "secret/x", ""); err != nil {
		t.Errorf("unexpected error on discard logger: %v", err)
	}
}

func TestNewLogger_CreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "audit.log")

	l, err := audit.NewLogger(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = l.Record(audit.EventSynced, "secret/app", "")

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	if !strings.Contains(string(data), "synced") {
		t.Error("expected 'synced' in log file")
	}
}
