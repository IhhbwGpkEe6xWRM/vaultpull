package envtag_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envtag"
)

func makeEntry(path string, tags []string) envtag.Entry {
	return envtag.Entry{
		Path:    path,
		Secrets: map[string]string{"KEY": "val"},
		Tagger:  envtag.New(tags),
	}
}

func TestFilter_EmptyFilter_ReturnsAll(t *testing.T) {
	entries := []envtag.Entry{
		makeEntry("a", []string{"env:prod"}),
		makeEntry("b", []string{"env:staging"}),
	}
	out := envtag.Filter(entries, nil)
	if len(out) != 2 {
		t.Fatalf("expected 2, got %d", len(out))
	}
}

func TestFilter_MatchingTag(t *testing.T) {
	entries := []envtag.Entry{
		makeEntry("a", []string{"env:prod"}),
		makeEntry("b", []string{"env:staging"}),
	}
	filter := envtag.ParseFilter([]string{"env:prod"})
	out := envtag.Filter(entries, filter)
	if len(out) != 1 || out[0].Path != "a" {
		t.Errorf("expected only entry 'a', got %v", out)
	}
}

func TestFilter_NoMatch_ReturnsEmpty(t *testing.T) {
	entries := []envtag.Entry{
		makeEntry("a", []string{"env:prod"}),
	}
	filter := envtag.ParseFilter([]string{"env:canary"})
	out := envtag.Filter(entries, filter)
	if len(out) != 0 {
		t.Errorf("expected empty result, got %d entries", len(out))
	}
}

func TestFilter_MultipleTagsRequired(t *testing.T) {
	entries := []envtag.Entry{
		makeEntry("a", []string{"env:prod", "team:platform"}),
		makeEntry("b", []string{"env:prod"}),
	}
	filter := envtag.ParseFilter([]string{"env:prod", "team:platform"})
	out := envtag.Filter(entries, filter)
	if len(out) != 1 || out[0].Path != "a" {
		t.Errorf("expected only entry 'a', got %v", out)
	}
}

func TestParseFilter_ReturnsTags(t *testing.T) {
	tags := envtag.ParseFilter([]string{"env:prod", "region:us-east"})
	if len(tags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags))
	}
	if tags[1].Key != "region" || tags[1].Value != "us-east" {
		t.Errorf("unexpected tag: %+v", tags[1])
	}
}
