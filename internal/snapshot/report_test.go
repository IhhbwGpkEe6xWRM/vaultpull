package snapshot_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/snapshot"
)

func TestReport_NoPrevious_AllAdded(t *testing.T) {
	current := map[string]string{"API_KEY": "abc", "DB_URL": "postgres://"}
	var buf strings.Builder

	changes, err := snapshot.Report(&buf, nil, current)
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	for _, c := range changes {
		if c.Kind != snapshot.Added {
			t.Errorf("key %q kind = %q, want %q", c.Key, c.Kind, snapshot.Added)
		}
	}
	if !strings.Contains(buf.String(), "added") {
		t.Errorf("output %q missing 'added'", buf.String())
	}
}

func TestReport_RemovedKey(t *testing.T) {
	prev := &snapshot.Snapshot{
		Secrets: map[string]string{"OLD_KEY": "value", "KEEP": "yes"},
	}
	current := map[string]string{"KEEP": "yes"}
	var buf strings.Builder

	changes, err := snapshot.Report(&buf, prev, current)
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	found := false
	for _, c := range changes {
		if c.Key == "OLD_KEY" && c.Kind == snapshot.Removed {
			found = true
		}
	}
	if !found {
		t.Error("expected OLD_KEY to be reported as removed")
	}
}

func TestReport_ModifiedKey(t *testing.T) {
	prev := &snapshot.Snapshot{
		Secrets: map[string]string{"TOKEN": "old"},
	}
	current := map[string]string{"TOKEN": "new"}
	var buf strings.Builder

	changes, err := snapshot.Report(&buf, prev, current)
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}

	for _, c := range changes {
		if c.Key == "TOKEN" && c.Kind != snapshot.Modified {
			t.Errorf("TOKEN kind = %q, want modified", c.Kind)
		}
	}
}

func TestReport_Unchanged_NotInOutput(t *testing.T) {
	prev := &snapshot.Snapshot{
		Secrets: map[string]string{"STABLE": "same"},
	}
	current := map[string]string{"STABLE": "same"}
	var buf strings.Builder

	_, err := snapshot.Report(&buf, prev, current)
	if err != nil {
		t.Fatalf("Report() error = %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("expected empty output for unchanged secrets, got %q", buf.String())
	}
}
