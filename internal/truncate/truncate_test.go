package truncate_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/truncate"
)

func TestValue_ShortString_Unchanged(t *testing.T) {
	tr := truncate.New()
	got := tr.Value("short")
	if got != "short" {
		t.Fatalf("expected 'short', got %q", got)
	}
}

func TestValue_LongString_Truncated(t *testing.T) {
	tr := truncate.NewWithLimit(10)
	input := "this is a very long secret value"
	got := tr.Value(input)
	if len(got) != 13 { // 10 + len("...")
		t.Fatalf("expected length 13, got %d", len(got))
	}
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("expected suffix '...', got %q", got)
	}
}

func TestValue_ExactLimit_Unchanged(t *testing.T) {
	tr := truncate.NewWithLimit(5)
	got := tr.Value("hello")
	if got != "hello" {
		t.Fatalf("expected 'hello', got %q", got)
	}
}

func TestMap_TruncatesLongValues(t *testing.T) {
	tr := truncate.NewWithLimit(5)
	secrets := map[string]string{
		"SHORT": "hi",
		"LONG":  "this is too long",
	}
	out := tr.Map(secrets)
	if out["SHORT"] != "hi" {
		t.Errorf("SHORT should be unchanged")
	}
	if !strings.HasSuffix(out["LONG"], "...") {
		t.Errorf("LONG should be truncated")
	}
}

func TestContainsLong_True(t *testing.T) {
	tr := truncate.NewWithLimit(5)
	secrets := map[string]string{"KEY": "way too long value"}
	if !tr.ContainsLong(secrets) {
		t.Error("expected ContainsLong to return true")
	}
}

func TestContainsLong_False(t *testing.T) {
	tr := truncate.NewWithLimit(50)
	secrets := map[string]string{"KEY": "short"}
	if tr.ContainsLong(secrets) {
		t.Error("expected ContainsLong to return false")
	}
}

func TestSummary_NoLongValues(t *testing.T) {
	tr := truncate.NewWithLimit(50)
	secrets := map[string]string{"A": "x"}
	if got := tr.Summary(secrets); got != "" {
		t.Errorf("expected empty summary, got %q", got)
	}
}

func TestSummary_WithLongValues(t *testing.T) {
	tr := truncate.NewWithLimit(5)
	secrets := map[string]string{
		"A": "toolong",
		"B": "alsotoolong",
	}
	got := tr.Summary(secrets)
	if !strings.Contains(got, "2") {
		t.Errorf("expected summary to mention 2, got %q", got)
	}
}

func TestNewWithLimit_ZeroUsesDefault(t *testing.T) {
	tr := truncate.NewWithLimit(0)
	long := strings.Repeat("x", 200)
	got := tr.Value(long)
	if len(got) <= 80 {
		// truncated correctly
		return
	}
	t.Errorf("expected truncation at default limit, got length %d", len(got))
}
