package envtag_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envtag"
)

func TestNew_ParsesPairs(t *testing.T) {
	tagger := envtag.New([]string{"env:production", "team:platform"})
	tags := tagger.Tags()
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[0].Key != "env" || tags[0].Value != "production" {
		t.Errorf("unexpected first tag: %+v", tags[0])
	}
}

func TestNew_IgnoresMalformed(t *testing.T) {
	tagger := envtag.New([]string{"nocoherence", ":emptykey", "valid:yes"})
	tags := tagger.Tags()
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Key != "valid" {
		t.Errorf("unexpected tag key: %s", tags[0].Key)
	}
}

func TestNew_EmptySlice(t *testing.T) {
	tagger := envtag.New(nil)
	if len(tagger.Tags()) != 0 {
		t.Error("expected no tags")
	}
}

func TestAnnotate_AddsTagKeys(t *testing.T) {
	tagger := envtag.New([]string{"env:staging"})
	secrets := map[string]string{"DB_URL": "postgres://localhost"}
	out := tagger.Annotate(secrets)
	if out["DB_URL"] != "postgres://localhost" {
		t.Error("original key should be preserved")
	}
	if out["__tag_ENV"] != "staging" {
		t.Errorf("expected __tag_ENV=staging, got %q", out["__tag_ENV"])
	}
}

func TestAnnotate_DoesNotMutateInput(t *testing.T) {
	tagger := envtag.New([]string{"env:prod"})
	secrets := map[string]string{"KEY": "val"}
	tagger.Annotate(secrets)
	if _, ok := secrets["__tag_ENV"]; ok {
		t.Error("original map should not be mutated")
	}
}

func TestMatchesAll_AllPresent(t *testing.T) {
	tagger := envtag.New([]string{"env:prod", "team:platform"})
	filter := []envtag.Tag{{Key: "env", Value: "prod"}}
	if !tagger.MatchesAll(filter) {
		t.Error("expected match")
	}
}

func TestMatchesAll_MissingTag(t *testing.T) {
	tagger := envtag.New([]string{"env:prod"})
	filter := []envtag.Tag{{Key: "team", Value: "platform"}}
	if tagger.MatchesAll(filter) {
		t.Error("expected no match")
	}
}

func TestMatchesAll_WrongValue(t *testing.T) {
	tagger := envtag.New([]string{"env:staging"})
	filter := []envtag.Tag{{Key: "env", Value: "prod"}}
	if tagger.MatchesAll(filter) {
		t.Error("expected no match due to value mismatch")
	}
}

func TestMatchesAll_EmptyFilter(t *testing.T) {
	tagger := envtag.New([]string{"env:prod"})
	if !tagger.MatchesAll(nil) {
		t.Error("empty filter should match everything")
	}
}
