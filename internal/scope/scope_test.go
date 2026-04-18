package scope_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/scope"
)

func TestNew_TrimsSlashes(t *testing.T) {
	s := scope.New([]string{"/secrets/app/", "//infra"})
	prefixes := s.Prefixes()
	if len(prefixes) != 2 || prefixes[0] != "secrets/app" || prefixes[1] != "infra" {
		t.Fatalf("unexpected prefixes: %v", prefixes)
	}
}

func TestNew_EmptyPrefixesIgnored(t *testing.T) {
	s := scope.New([]string{"/", "", "   "})
	if len(s.Prefixes()) != 0 {
		t.Fatalf("expected no prefixes, got %v", s.Prefixes())
	}
}

func TestAllows_NoPrefixes_AllowsAll(t *testing.T) {
	s := scope.New(nil)
	if !s.Allows("anything/goes") {
		t.Fatal("expected all paths allowed when no prefixes set")
	}
}

func TestAllows_ExactMatch(t *testing.T) {
	s := scope.New([]string{"secrets/app"})
	if !s.Allows("secrets/app") {
		t.Fatal("expected exact match to be allowed")
	}
}

func TestAllows_ChildPath(t *testing.T) {
	s := scope.New([]string{"secrets/app"})
	if !s.Allows("secrets/app/db") {
		t.Fatal("expected child path to be allowed")
	}
}

func TestAllows_UnrelatedPath(t *testing.T) {
	s := scope.New([]string{"secrets/app"})
	if s.Allows("secrets/other") {
		t.Fatal("expected unrelated path to be denied")
	}
}

func TestAllows_PrefixSubstring_NotAllowed(t *testing.T) {
	s := scope.New([]string{"sec"})
	if s.Allows("secrets/app") {
		t.Fatal("prefix substring without slash boundary should not match")
	}
}

func TestFilter_ReturnsAllowed(t *testing.T) {
	s := scope.New([]string{"prod"})
	input := []string{"prod/db", "staging/db", "prod/cache", "dev/api"}
	got := s.Filter(input)
	if len(got) != 2 || got[0] != "prod/db" || got[1] != "prod/cache" {
		t.Fatalf("unexpected filter result: %v", got)
	}
}

func TestFilter_EmptyInput(t *testing.T) {
	s := scope.New([]string{"prod"})
	got := s.Filter(nil)
	if len(got) != 0 {
		t.Fatalf("expected empty result, got %v", got)
	}
}
