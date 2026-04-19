package depcheck_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/depcheck"
)

func TestCheck_AllPresent(t *testing.T) {
	c := depcheck.New([]string{"DB_HOST", "DB_PASS"})
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PASS": "s3cr3t"}
	if v := c.Check(secrets); len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheck_MissingKey(t *testing.T) {
	c := depcheck.New([]string{"DB_HOST", "DB_PASS"})
	secrets := map[string]string{"DB_HOST": "localhost"}
	v := c.Check(secrets)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "DB_PASS" || v[0].Reason != "missing" {
		t.Errorf("unexpected violation: %+v", v[0])
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	c := depcheck.New([]string{"API_KEY"})
	v := c.Check(map[string]string{"API_KEY": ""})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Reason != "empty value" {
		t.Errorf("expected 'empty value', got %q", v[0].Reason)
	}
}

func TestCheck_NoRequiredKeys(t *testing.T) {
	c := depcheck.New(nil)
	if v := c.Check(map[string]string{"X": "y"}); len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheck_MultipleViolations(t *testing.T) {
	c := depcheck.New([]string{"A", "B", "C"})
	v := c.Check(map[string]string{"A": "ok"})
	if len(v) != 2 {
		t.Fatalf("expected 2 violations, got %d", len(v))
	}
}

func TestSummary_Empty(t *testing.T) {
	if s := depcheck.Summary(nil); s != "" {
		t.Errorf("expected empty summary, got %q", s)
	}
}

func TestSummary_WithViolations(t *testing.T) {
	v := []depcheck.Violation{{Key: "TOKEN", Reason: "missing"}}
	s := depcheck.Summary(v)
	if s == "" {
		t.Error("expected non-empty summary")
	}
}

func TestViolation_Error(t *testing.T) {
	v := depcheck.Violation{Key: "FOO", Reason: "missing"}
	if v.Error() != "FOO: missing" {
		t.Errorf("unexpected error string: %s", v.Error())
	}
}
