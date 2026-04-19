package inherit_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpull/internal/inherit"
)

func makeReader(data map[string]map[string]string) func(string) (map[string]string, error) {
	return func(path string) (map[string]string, error) {
		if m, ok := data[path]; ok {
			return m, nil
		}
		return map[string]string{}, nil
	}
}

func TestResolve_SinglePath(t *testing.T) {
	r := inherit.New(makeReader(map[string]map[string]string{
		"app": {"KEY": "base"},
	}))
	got, err := r.Resolve("app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "base" {
		t.Errorf("expected base, got %s", got["KEY"])
	}
}

func TestResolve_ChildOverridesParent(t *testing.T) {
	r := inherit.New(makeReader(map[string]map[string]string{
		"app":         {"KEY": "parent", "SHARED": "from-parent"},
		"app/service": {"KEY": "child"},
	}))
	got, err := r.Resolve("app/service")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "child" {
		t.Errorf("expected child, got %s", got["KEY"])
	}
	if got["SHARED"] != "from-parent" {
		t.Errorf("expected from-parent, got %s", got["SHARED"])
	}
}

func TestResolve_ThreeLevels(t *testing.T) {
	r := inherit.New(makeReader(map[string]map[string]string{
		"a":     {"X": "1"},
		"a/b":   {"Y": "2"},
		"a/b/c": {"Z": "3"},
	}))
	got, err := r.Resolve("a/b/c")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["X"] != "1" || got["Y"] != "2" || got["Z"] != "3" {
		t.Errorf("unexpected result: %v", got)
	}
}

func TestResolve_LeadingSlashStripped(t *testing.T) {
	r := inherit.New(makeReader(map[string]map[string]string{
		"app": {"K": "v"},
	}))
	got, err := r.Resolve("/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["K"] != "v" {
		t.Errorf("expected v, got %s", got["K"])
	}
}

func TestResolve_ReaderError(t *testing.T) {
	readErr := errors.New("vault unavailable")
	r := inherit.New(func(string) (map[string]string, error) {
		return nil, readErr
	})
	_, err := r.Resolve("app/svc")
	if !errors.Is(err, readErr) {
		t.Errorf("expected readErr, got %v", err)
	}
}

func TestResolve_EmptyPath(t *testing.T) {
	r := inherit.New(makeReader(nil))
	got, err := r.Resolve("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}
