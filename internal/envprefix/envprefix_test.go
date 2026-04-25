package envprefix_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envprefix"
)

func TestNew_TrimsUnderscores(t *testing.T) {
	tx := envprefix.New("_APP_")
	got := tx.Apply(map[string]string{"KEY": "val"})
	if _, ok := got["APP_KEY"]; !ok {
		t.Fatalf("expected key APP_KEY, got %v", got)
	}
}

func TestApply_EmptyPrefix_ReturnsOriginal(t *testing.T) {
	tx := envprefix.New("")
	src := map[string]string{"FOO": "bar"}
	got := tx.Apply(src)
	if got["FOO"] != "bar" {
		t.Fatalf("expected FOO=bar, got %v", got)
	}
}

func TestApply_AddsPrefixToAllKeys(t *testing.T) {
	tx := envprefix.New("SVC")
	src := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	got := tx.Apply(src)
	for _, k := range []string{"SVC_DB_HOST", "SVC_DB_PORT"} {
		if _, ok := got[k]; !ok {
			t.Errorf("missing key %s in %v", k, got)
		}
	}
	if len(got) != 2 {
		t.Errorf("expected 2 keys, got %d", len(got))
	}
}

func TestApply_DoesNotMutateSource(t *testing.T) {
	tx := envprefix.New("SVC")
	src := map[string]string{"DB_HOST": "localhost"}
	_ = tx.Apply(src)
	if _, ok := src["SVC_DB_HOST"]; ok {
		t.Error("Apply must not mutate the source map")
	}
	if _, ok := src["DB_HOST"]; !ok {
		t.Error("Apply must preserve original keys in source map")
	}
}

func TestStrip_RemovesPrefixFromMatchingKeys(t *testing.T) {
	tx := envprefix.New("SVC")
	src := map[string]string{"SVC_DB_HOST": "localhost", "OTHER": "x"}
	got := tx.Strip(src)
	if got["DB_HOST"] != "localhost" {
		t.Errorf("expected DB_HOST=localhost, got %v", got)
	}
	if got["OTHER"] != "x" {
		t.Errorf("expected OTHER=x, got %v", got)
	}
}

func TestStrip_EmptyPrefix_ReturnsOriginal(t *testing.T) {
	tx := envprefix.New("")
	src := map[string]string{"FOO": "1"}
	got := tx.Strip(src)
	if got["FOO"] != "1" {
		t.Fatalf("unexpected result %v", got)
	}
}

func TestHasPrefix_Match(t *testing.T) {
	tx := envprefix.New("APP")
	if !tx.HasPrefix("APP_SECRET") {
		t.Error("expected HasPrefix to return true")
	}
}

func TestHasPrefix_NoMatch(t *testing.T) {
	tx := envprefix.New("APP")
	if tx.HasPrefix("OTHER_SECRET") {
		t.Error("expected HasPrefix to return false")
	}
}

func TestHasPrefix_EmptyPrefix_AlwaysTrue(t *testing.T) {
	tx := envprefix.New("")
	if !tx.HasPrefix("ANYTHING") {
		t.Error("empty prefix should match everything")
	}
}
