package sanitize_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/sanitize"
)

func TestValue_NoOptions_Unchanged(t *testing.T) {
	s := sanitize.New()
	if got := s.Value("  hello\x00world  "); got != "  hello\x00world  " {
		t.Fatalf("unexpected: %q", got)
	}
}

func TestValue_TrimSpace(t *testing.T) {
	s := sanitize.New(sanitize.WithTrimSpace())
	if got := s.Value("  hello  "); got != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestValue_StripControl(t *testing.T) {
	s := sanitize.New(sanitize.WithStripControl())
	got := s.Value("hel\x00lo\x01")
	if strings.ContainsAny(got, "\x00\x01") {
		t.Fatalf("control chars not stripped: %q", got)
	}
	if got != "hello" {
		t.Fatalf("got %q", got)
	}
}

func TestValue_StripControl_PreservesTab(t *testing.T) {
	s := sanitize.New(sanitize.WithStripControl())
	if got := s.Value("a\tb"); got != "a\tb" {
		t.Fatalf("tab should be preserved, got %q", got)
	}
}

func TestValue_MaxLen_Truncates(t *testing.T) {
	s := sanitize.New(sanitize.WithMaxLen(5))
	if got := s.Value("abcdefgh"); got != "abcde" {
		t.Fatalf("got %q", got)
	}
}

func TestValue_MaxLen_ShortValue_Unchanged(t *testing.T) {
	s := sanitize.New(sanitize.WithMaxLen(10))
	if got := s.Value("hi"); got != "hi" {
		t.Fatalf("got %q", got)
	}
}

func TestValue_CombinedOptions(t *testing.T) {
	s := sanitize.New(sanitize.WithTrimSpace(), sanitize.WithStripControl(), sanitize.WithMaxLen(4))
	got := s.Value("  ab\x00cdef  ")
	if got != "abcd" {
		t.Fatalf("got %q", got)
	}
}

func TestMap_AppliesValueToAll(t *testing.T) {
	s := sanitize.New(sanitize.WithTrimSpace())
	in := map[string]string{"A": "  foo  ", "B": " bar "}
	out := s.Map(in)
	if out["A"] != "foo" || out["B"] != "bar" {
		t.Fatalf("got %v", out)
	}
}

func TestMap_DoesNotMutateInput(t *testing.T) {
	s := sanitize.New(sanitize.WithTrimSpace())
	in := map[string]string{"K": "  v  "}
	_ = s.Map(in)
	if in["K"] != "  v  " {
		t.Fatal("input map was mutated")
	}
}
