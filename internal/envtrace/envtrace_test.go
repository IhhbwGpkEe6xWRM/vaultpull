package envtrace

import (
	"testing"
)

func TestRecord_StoresEntry(t *testing.T) {
	tr := New()
	if err := tr.Record("DB_HOST", "secret/app"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	e, ok := tr.Lookup("DB_HOST")
	if !ok {
		t.Fatal("expected entry to be found")
	}
	if e.Path != "secret/app" {
		t.Errorf("expected path %q, got %q", "secret/app", e.Path)
	}
}

func TestRecord_EmptyKey_ReturnsError(t *testing.T) {
	tr := New()
	if err := tr.Record("", "secret/app"); err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestRecord_WhitespaceKey_ReturnsError(t *testing.T) {
	tr := New()
	if err := tr.Record("   ", "secret/app"); err == nil {
		t.Fatal("expected error for whitespace-only key")
	}
}

func TestLookup_MissingKey_ReturnsFalse(t *testing.T) {
	tr := New()
	_, ok := tr.Lookup("MISSING")
	if ok {
		t.Fatal("expected lookup to return false for unknown key")
	}
}

func TestRecordAll_RecordsAllKeys(t *testing.T) {
	tr := New()
	secrets := map[string]string{
		"API_KEY": "abc123",
		"API_SECRET": "xyz789",
	}
	if err := tr.RecordAll(secrets, "secret/api"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tr.Len() != 2 {
		t.Errorf("expected 2 entries, got %d", tr.Len())
	}
	for k := range secrets {
		e, ok := tr.Lookup(k)
		if !ok {
			t.Errorf("expected key %q to be recorded", k)
		}
		if e.Path != "secret/api" {
			t.Errorf("key %q: expected path %q, got %q", k, "secret/api", e.Path)
		}
	}
}

func TestEntries_SortedByKey(t *testing.T) {
	tr := New()
	_ = tr.Record("Z_KEY", "path/z")
	_ = tr.Record("A_KEY", "path/a")
	_ = tr.Record("M_KEY", "path/m")

	entries := tr.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}
	expected := []string{"A_KEY", "M_KEY", "Z_KEY"}
	for i, e := range entries {
		if e.Key != expected[i] {
			t.Errorf("position %d: expected %q, got %q", i, expected[i], e.Key)
		}
	}
}

func TestLen_EmptyTracer(t *testing.T) {
	tr := New()
	if tr.Len() != 0 {
		t.Errorf("expected 0, got %d", tr.Len())
	}
}
