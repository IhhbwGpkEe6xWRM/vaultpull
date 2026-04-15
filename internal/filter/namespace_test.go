package filter_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/filter"
)

func TestNewMatcher_TrimsSlashes(t *testing.T) {
	m := filter.NewMatcher("/team/backend/")
	if !m.Match("team/backend/db") {
		t.Error("expected match after trimming slashes")
	}
}

func TestMatch_EmptyNamespace_MatchesAll(t *testing.T) {
	m := filter.NewMatcher("")
	paths := []string{"secret/app", "team/backend/db", "other"}
	for _, p := range paths {
		if !m.Match(p) {
			t.Errorf("expected empty namespace to match %q", p)
		}
	}
}

func TestMatch_ExactPath(t *testing.T) {
	m := filter.NewMatcher("team/backend")
	if !m.Match("team/backend") {
		t.Error("expected exact path to match")
	}
}

func TestMatch_ChildPath(t *testing.T) {
	m := filter.NewMatcher("team/backend")
	if !m.Match("team/backend/db") {
		t.Error("expected child path to match")
	}
}

func TestMatch_UnrelatedPath(t *testing.T) {
	m := filter.NewMatcher("team/backend")
	if m.Match("team/frontend/db") {
		t.Error("expected unrelated path not to match")
	}
}

func TestMatch_PartialPrefixNoSlash(t *testing.T) {
	m := filter.NewMatcher("team/back")
	// "team/backend" should NOT match "team/back" as a child
	if m.Match("team/backend") {
		t.Error("partial prefix without slash separator should not match")
	}
}

func TestStripNamespace_RemovesPrefix(t *testing.T) {
	m := filter.NewMatcher("team/backend")
	got := m.StripNamespace("team/backend/db")
	if got != "db" {
		t.Errorf("expected %q, got %q", "db", got)
	}
}

func TestStripNamespace_NoMatch_ReturnsClean(t *testing.T) {
	m := filter.NewMatcher("team/backend")
	got := m.StripNamespace("/other/path/")
	if got != "other/path" {
		t.Errorf("expected %q, got %q", "other/path", got)
	}
}

func TestFilterPaths_ReturnsMatchingOnly(t *testing.T) {
	m := filter.NewMatcher("team/backend")
	input := []string{
		"team/backend/db",
		"team/frontend/ui",
		"team/backend/cache",
		"infra/network",
	}
	got := m.FilterPaths(input)
	if len(got) != 2 {
		t.Fatalf("expected 2 matches, got %d: %v", len(got), got)
	}
}

func TestFilterPaths_EmptyInput(t *testing.T) {
	m := filter.NewMatcher("team/backend")
	got := m.FilterPaths([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty result, got %v", got)
	}
}
