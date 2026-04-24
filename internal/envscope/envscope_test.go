package envscope_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envscope"
)

func newScope(t *testing.T, prefixes []string) *envscope.Scope {
	t.Helper()
	s, err := envscope.New(prefixes)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNew_TrimsAndUppercases(t *testing.T) {
	s := newScope(t, []string{"app_", " db "})
	if !s.Allows("APP_HOST") {
		t.Error("expected APP_HOST to be allowed")
	}
	if !s.Allows("DB_PASS") {
		t.Error("expected DB_PASS to be allowed")
	}
}

func TestNew_EmptyPrefixesIgnored(t *testing.T) {
	s := newScope(t, []string{"", "  "})
	if !s.Allows("ANYTHING") {
		t.Error("empty scope should allow all keys")
	}
}

func TestAllows_NoPrefixes_AllowsAll(t *testing.T) {
	s := newScope(t, nil)
	for _, k := range []string{"FOO", "BAR_BAZ", "X"} {
		if !s.Allows(k) {
			t.Errorf("expected %q to be allowed with no prefixes", k)
		}
	}
}

func TestAllows_ExactMatch(t *testing.T) {
	s := newScope(t, []string{"TOKEN"})
	if !s.Allows("TOKEN") {
		t.Error("exact match should be allowed")
	}
	if s.Allows("TOKENIZER") {
		t.Error("TOKENIZER should not match prefix TOKEN without underscore boundary")
	}
}

func TestAllows_ChildKey(t *testing.T) {
	s := newScope(t, []string{"APP"})
	if !s.Allows("APP_HOST") {
		t.Error("APP_HOST should be allowed under APP prefix")
	}
}

func TestAllows_UnrelatedKey(t *testing.T) {
	s := newScope(t, []string{"APP"})
	if s.Allows("DB_HOST") {
		t.Error("DB_HOST should not be allowed under APP prefix")
	}
}

func TestFilter_RemovesOutOfScopeKeys(t *testing.T) {
	s := newScope(t, []string{"APP"})
	input := map[string]string{
		"APP_HOST": "localhost",
		"DB_PASS":  "secret",
		"APP_PORT": "8080",
	}
	out := s.Filter(input)
	if len(out) != 2 {
		t.Fatalf("expected 2 keys, got %d", len(out))
	}
	if _, ok := out["DB_PASS"]; ok {
		t.Error("DB_PASS should have been filtered out")
	}
}

func TestValidate_NoViolations(t *testing.T) {
	s := newScope(t, []string{"APP"})
	err := s.Validate(map[string]string{"APP_HOST": "localhost"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_ReturnsViolations(t *testing.T) {
	s := newScope(t, []string{"APP"})
	err := s.Validate(map[string]string{
		"APP_HOST": "localhost",
		"DB_PASS":  "secret",
	})
	if err == nil {
		t.Fatal("expected error for out-of-scope key")
	}
	if !contains(err.Error(), "DB_PASS") {
		t.Errorf("error should mention DB_PASS, got: %v", err)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
			return false
		})())
}
