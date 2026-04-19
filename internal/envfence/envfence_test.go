package envfence_test

import (
	"sort"
	"testing"

	"github.com/yourusername/vaultpull/internal/envfence"
)

func TestNew_InvalidPattern(t *testing.T) {
	_, err := envfence.New(envfence.Allow, []string{"[invalid"})
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestApply_AllowMode_PassesMatchedKeys(t *testing.T) {
	f, err := envfence.New(envfence.Allow, []string{"DB_.*"})
	if err != nil {
		t.Fatal(err)
	}
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432", "API_KEY": "secret"}
	out := f.Apply(secrets)
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to pass")
	}
	if _, ok := out["API_KEY"]; ok {
		t.Error("expected API_KEY to be blocked")
	}
}

func TestApply_DenyMode_BlocksMatchedKeys(t *testing.T) {
	f, err := envfence.New(envfence.Deny, []string{".*_SECRET", ".*_TOKEN"})
	if err != nil {
		t.Fatal(err)
	}
	secrets := map[string]string{"APP_SECRET": "x", "APP_TOKEN": "y", "APP_NAME": "z"}
	out := f.Apply(secrets)
	if len(out) != 1 {
		t.Fatalf("expected 1 key, got %d", len(out))
	}
	if _, ok := out["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to pass")
	}
}

func TestApply_NoPatterns_PassesAll(t *testing.T) {
	f, _ := envfence.New(envfence.Allow, nil)
	secrets := map[string]string{"A": "1", "B": "2"}
	out := f.Apply(secrets)
	if len(out) != 2 {
		t.Errorf("expected all keys, got %d", len(out))
	}
}

func TestApply_CaseInsensitive(t *testing.T) {
	f, _ := envfence.New(envfence.Allow, []string{"db_host"})
	out := f.Apply(map[string]string{"DB_HOST": "localhost"})
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected case-insensitive match")
	}
}

func TestBlocked_ReturnsRemovedKeys(t *testing.T) {
	f, _ := envfence.New(envfence.Deny, []string{"INTERNAL_.*"})
	secrets := map[string]string{"INTERNAL_FLAG": "1", "PUBLIC_KEY": "abc"}
	blocked := f.Blocked(secrets)
	if len(blocked) != 1 || blocked[0] != "INTERNAL_FLAG" {
		t.Errorf("unexpected blocked keys: %v", blocked)
	}
}

func TestApply_EmptySecrets(t *testing.T) {
	f, _ := envfence.New(envfence.Allow, []string{".*"})
	out := f.Apply(map[string]string{})
	if len(out) != 0 {
		t.Error("expected empty output")
	}
}

func TestApply_MultiplePatterns_AnyMatch(t *testing.T) {
	f, _ := envfence.New(envfence.Allow, []string{"DB_.*", "REDIS_.*"})
	secrets := map[string]string{"DB_HOST": "h", "REDIS_URL": "r", "OTHER": "o"}
	out := f.Apply(secrets)
	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "DB_HOST" || keys[1] != "REDIS_URL" {
		t.Errorf("unexpected keys: %v", keys)
	}
}
