package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/checkpoint"
)

func tempFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "checkpoint.json")
}

func TestNewStore_EmptyWhenMissing(t *testing.T) {
	s, err := checkpoint.NewStore(tempFile(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := s.Get("secret/app")
	if ok {
		t.Fatal("expected no entry")
	}
}

func TestSetAndGet_RoundTrip(t *testing.T) {
	path := tempFile(t)
	s, _ := checkpoint.NewStore(path)
	now := time.Now().UTC().Truncate(time.Second)
	if err := s.Set("secret/app", "abc123", now); err != nil {
		t.Fatalf("Set: %v", err)
	}
	e, ok := s.Get("secret/app")
	if !ok {
		t.Fatal("expected entry")
	}
	if e.Checksum != "abc123" {
		t.Errorf("checksum = %q, want abc123", e.Checksum)
	}
}

func TestSet_PersistsToDisk(t *testing.T) {
	path := tempFile(t)
	s, _ := checkpoint.NewStore(path)
	now := time.Now().UTC()
	s.Set("secret/db", "xyz", now)

	s2, err := checkpoint.NewStore(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	e, ok := s2.Get("secret/db")
	if !ok || e.Checksum != "xyz" {
		t.Errorf("expected persisted entry, got ok=%v checksum=%q", ok, e.Checksum)
	}
}

func TestIsFresh_ReturnsFalseWhenMissing(t *testing.T) {
	s, _ := checkpoint.NewStore(tempFile(t))
	if s.IsFresh("secret/x", "sum", time.Hour, time.Now()) {
		t.Fatal("expected not fresh")
	}
}

func TestIsFresh_ReturnsTrueWithinTTL(t *testing.T) {
	s, _ := checkpoint.NewStore(tempFile(t))
	now := time.Now().UTC()
	s.Set("secret/app", "sum1", now.Add(-10*time.Minute))
	if !s.IsFresh("secret/app", "sum1", time.Hour, now) {
		t.Fatal("expected fresh")
	}
}

func TestIsFresh_ReturnsFalseWhenExpired(t *testing.T) {
	s, _ := checkpoint.NewStore(tempFile(t))
	now := time.Now().UTC()
	s.Set("secret/app", "sum1", now.Add(-2*time.Hour))
	if s.IsFresh("secret/app", "sum1", time.Hour, now) {
		t.Fatal("expected not fresh")
	}
}

func TestIsFresh_ReturnsFalseOnChecksumMismatch(t *testing.T) {
	s, _ := checkpoint.NewStore(tempFile(t))
	now := time.Now().UTC()
	s.Set("secret/app", "old", now.Add(-1*time.Minute))
	if s.IsFresh("secret/app", "new", time.Hour, now) {
		t.Fatal("expected not fresh on mismatch")
	}
}

func TestNewStore_InvalidJSON(t *testing.T) {
	path := tempFile(t)
	os.WriteFile(path, []byte("not json"), 0600)
	_, err := checkpoint.NewStore(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
