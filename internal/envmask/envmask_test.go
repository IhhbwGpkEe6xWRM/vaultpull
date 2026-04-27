package envmask_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envmask"
)

func newMasker(t *testing.T, opts ...envmask.Option) *envmask.Masker {
	t.Helper()
	m, err := envmask.New(opts...)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return m
}

func TestIsSensitive_MatchesPassword(t *testing.T) {
	m := newMasker(t)
	if !m.IsSensitive("DB_PASSWORD") {
		t.Error("expected DB_PASSWORD to be sensitive")
	}
}

func TestIsSensitive_MatchesToken(t *testing.T) {
	m := newMasker(t)
	if !m.IsSensitive("GITHUB_TOKEN") {
		t.Error("expected GITHUB_TOKEN to be sensitive")
	}
}

func TestIsSensitive_NonSensitiveKey(t *testing.T) {
	m := newMasker(t)
	if m.IsSensitive("APP_ENV") {
		t.Error("expected APP_ENV to not be sensitive")
	}
}

func TestApply_MasksSensitiveValues(t *testing.T) {
	m := newMasker(t)
	src := map[string]string{
		"DB_PASSWORD": "s3cr3t",
		"APP_ENV":     "production",
	}
	out := m.Apply(src)
	if out["DB_PASSWORD"] != envmask.DefaultPlaceholder {
		t.Errorf("expected placeholder, got %q", out["DB_PASSWORD"])
	}
	if out["APP_ENV"] != "production" {
		t.Errorf("expected production, got %q", out["APP_ENV"])
	}
}

func TestApply_DoesNotMutateSource(t *testing.T) {
	m := newMasker(t)
	src := map[string]string{"API_KEY": "original"}
	_ = m.Apply(src)
	if src["API_KEY"] != "original" {
		t.Error("source map was mutated")
	}
}

func TestApply_CustomPlaceholder(t *testing.T) {
	m := newMasker(t, envmask.WithPlaceholder("[REDACTED]"))
	out := m.Apply(map[string]string{"SECRET_KEY": "abc"})
	if out["SECRET_KEY"] != "[REDACTED]" {
		t.Errorf("unexpected placeholder: %q", out["SECRET_KEY"])
	}
}

func TestApply_CustomPatterns(t *testing.T) {
	m := newMasker(t, envmask.WithPatterns([]string{"(?i)internal"}))
	src := map[string]string{
		"INTERNAL_HOST": "10.0.0.1",
		"DB_PASSWORD":   "should-not-mask",
	}
	out := m.Apply(src)
	if out["INTERNAL_HOST"] != envmask.DefaultPlaceholder {
		t.Error("expected INTERNAL_HOST to be masked")
	}
	if out["DB_PASSWORD"] != "should-not-mask" {
		t.Error("expected DB_PASSWORD to be unmasked with custom patterns")
	}
}

func TestMaskedKeys_ReturnsSortedSensitiveKeys(t *testing.T) {
	m := newMasker(t)
	src := map[string]string{
		"TOKEN":       "t",
		"APP_ENV":     "prod",
		"DB_PASSWORD": "p",
		"HOST":        "h",
	}
	keys := m.MaskedKeys(src)
	if len(keys) != 2 {
		t.Fatalf("expected 2 masked keys, got %d: %v", len(keys), keys)
	}
	if keys[0] != "DB_PASSWORD" || keys[1] != "TOKEN" {
		t.Errorf("unexpected order: %v", keys)
	}
}

func TestMaskedKeys_EmptyMap(t *testing.T) {
	m := newMasker(t)
	keys := m.MaskedKeys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("expected empty, got %v", keys)
	}
}
