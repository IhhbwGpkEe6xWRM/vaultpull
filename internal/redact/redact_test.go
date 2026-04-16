package redact_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/redact"
)

func TestIsSensitive_MatchesDefault(t *testing.T) {
	r := redact.New()
	cases := []string{"DB_PASSWORD", "API_TOKEN", "aws_secret", "PRIVATE_KEY", "auth_header"}
	for _, k := range cases {
		if !r.IsSensitive(k) {
			t.Errorf("expected %q to be sensitive", k)
		}
	}
}

func TestIsSensitive_NonSensitive(t *testing.T) {
	r := redact.New()
	cases := []string{"APP_ENV", "PORT", "LOG_LEVEL", "DATABASE_HOST"}
	for _, k := range cases {
		if r.IsSensitive(k) {
			t.Errorf("expected %q to NOT be sensitive", k)
		}
	}
}

func TestRedact_ReplacesValues(t *testing.T) {
	r := redact.New()
	input := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_ENV":     "production",
		"API_TOKEN":   "tok_abc",
	}
	out := r.Redact(input)
	if out["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("expected DB_PASSWORD to be redacted, got %q", out["DB_PASSWORD"])
	}
	if out["API_TOKEN"] != "[REDACTED]" {
		t.Errorf("expected API_TOKEN to be redacted, got %q", out["API_TOKEN"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected APP_ENV to be unchanged, got %q", out["APP_ENV"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	r := redact.New()
	input := map[string]string{"DB_PASSWORD": "s3cr3t"}
	_ = r.Redact(input)
	if input["DB_PASSWORD"] != "s3cr3t" {
		t.Error("Redact must not mutate the input map")
	}
}

func TestNewWithPatterns_CustomPlaceholder(t *testing.T) {
	r := redact.NewWithPatterns([]string{"pin"}, "***")
	out := r.Redact(map[string]string{"USER_PIN": "1234", "NAME": "alice"})
	if out["USER_PIN"] != "***" {
		t.Errorf("expected USER_PIN to be ***, got %q", out["USER_PIN"])
	}
	if out["NAME"] != "alice" {
		t.Errorf("expected NAME unchanged, got %q", out["NAME"])
	}
}

func TestRedact_EmptyMap(t *testing.T) {
	r := redact.New()
	out := r.Redact(map[string]string{})
	if len(out) != 0 {
		t.Errorf("expected empty map, got %v", out)
	}
}
