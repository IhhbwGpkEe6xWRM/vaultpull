package label_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/label"
)

func TestNew_ParsesPairs(t *testing.T) {
	tag := label.New([]string{"env=prod", "team=platform"})
	meta := tag.Apply(nil)
	if meta["env"] != "prod" {
		t.Fatalf("expected env=prod, got %q", meta["env"])
	}
	if meta["team"] != "platform" {
		t.Fatalf("expected team=platform, got %q", meta["team"])
	}
}

func TestNew_IgnoresMalformed(t *testing.T) {
	tag := label.New([]string{"nodivider", "ok=yes"})
	meta := tag.Apply(nil)
	if _, ok := meta["nodivider"]; ok {
		t.Fatal("malformed entry should be ignored")
	}
	if meta["ok"] != "yes" {
		t.Fatal("valid pair should be present")
	}
}

func TestApply_DoesNotOverwriteExisting(t *testing.T) {
	tag := label.New([]string{"env=prod"})
	meta := tag.Apply(map[string]string{"env": "staging"})
	if meta["env"] != "staging" {
		t.Fatalf("existing key should not be overwritten, got %q", meta["env"])
	}
}

func TestApply_NilMeta(t *testing.T) {
	tag := label.New([]string{"x=1"})
	meta := tag.Apply(nil)
	if meta["x"] != "1" {
		t.Fatal("expected x=1")
	}
}

func TestMatches_AllPresent(t *testing.T) {
	meta := map[string]string{"env": "prod", "region": "us-east"}
	filter := label.Set{"env": "prod"}
	if !label.Matches(meta, filter) {
		t.Fatal("expected match")
	}
}

func TestMatches_MissingKey(t *testing.T) {
	meta := map[string]string{"env": "prod"}
	filter := label.Set{"team": "platform"}
	if label.Matches(meta, filter) {
		t.Fatal("expected no match")
	}
}

func TestMatches_WrongValue(t *testing.T) {
	meta := map[string]string{"env": "staging"}
	filter := label.Set{"env": "prod"}
	if label.Matches(meta, filter) {
		t.Fatal("expected no match on wrong value")
	}
}

func TestMatches_EmptyFilter(t *testing.T) {
	meta := map[string]string{"env": "prod"}
	if !label.Matches(meta, label.Set{}) {
		t.Fatal("empty filter should match everything")
	}
}
