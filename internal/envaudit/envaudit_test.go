package envaudit_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/envaudit"
)

type fixedClock struct {
	t time.Time
}

func (c fixedClock) Now() time.Time { return c.t }

var epoch = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func newAuditor(t *testing.T) (*envaudit.Auditor, *bytes.Buffer) {
	t.Helper()
	var buf bytes.Buffer
	a := envaudit.New(
		envaudit.WithClock(fixedClock{t: epoch}),
		envaudit.WithWriter(&buf),
	)
	return a, &buf
}

func TestRecord_WritesJSONLine(t *testing.T) {
	a, buf := newAuditor(t)

	if err := a.Record("READ", "secret/app", nil); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("expected output, got none")
	}

	var entry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
}

func TestRecord_ContainsActionAndPath(t *testing.T) {
	a, buf := newAuditor(t)

	_ = a.Record("WRITE", "secret/db", nil)

	line := buf.String()
	if !strings.Contains(line, "WRITE") {
		t.Errorf("expected action WRITE in output, got: %s", line)
	}
	if !strings.Contains(line, "secret/db") {
		t.Errorf("expected path secret/db in output, got: %s", line)
	}
}

func TestRecord_ContainsTimestamp(t *testing.T) {
	a, buf := newAuditor(t)

	_ = a.Record("READ", "secret/app", nil)

	var entry map[string]interface{}
	_ = json.Unmarshal(buf.Bytes(), &entry)

	ts, ok := entry["timestamp"]
	if !ok {
		t.Fatal("expected timestamp field in audit entry")
	}
	if !strings.Contains(ts.(string), "2024-01-15") {
		t.Errorf("unexpected timestamp value: %v", ts)
	}
}

func TestRecord_WithMeta_IncludesMeta(t *testing.T) {
	a, buf := newAuditor(t)

	meta := map[string]string{"user": "alice", "env": "prod"}
	_ = a.Record("DELETE", "secret/old", meta)

	line := buf.String()
	if !strings.Contains(line, "alice") {
		t.Errorf("expected meta value alice in output, got: %s", line)
	}
	if !strings.Contains(line, "prod") {
		t.Errorf("expected meta value prod in output, got: %s", line)
	}
}

func TestRecord_EmptyAction_ReturnsError(t *testing.T) {
	a, _ := newAuditor(t)

	if err := a.Record("", "secret/app", nil); err == nil {
		t.Fatal("expected error for empty action, got nil")
	}
}

func TestRecord_EmptyPath_ReturnsError(t *testing.T) {
	a, _ := newAuditor(t)

	if err := a.Record("READ", "", nil); err == nil {
		t.Fatal("expected error for empty path, got nil")
	}
}

func TestRecord_MultipleEntries_EachOnOwnLine(t *testing.T) {
	a, buf := newAuditor(t)

	_ = a.Record("READ", "secret/a", nil)
	_ = a.Record("READ", "secret/b", nil)
	_ = a.Record("READ", "secret/c", nil)

	lines := strings.Split(strings.TrimSpace(buf.String()), "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}

	for i, line := range lines {
		var entry map[string]interface{}
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			t.Errorf("line %d is not valid JSON: %v", i, err)
		}
	}
}
