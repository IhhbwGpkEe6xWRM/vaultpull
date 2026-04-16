package multipath_test

import (
	"errors"
	"testing"

	"github.com/your-org/vaultpull/internal/multipath"
)

// stubReader satisfies multipath.Reader for testing.
type stubReader struct {
	data map[string]map[string]string
	err  error
}

func (s *stubReader) ReadSecrets(path string) (map[string]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	if m, ok := s.data[path]; ok {
		return m, nil
	}
	return map[string]string{}, nil
}

func TestNew_NilReader(t *testing.T) {
	_, err := multipath.New(nil, []string{"a"})
	if err == nil {
		t.Fatal("expected error for nil reader")
	}
}

func TestNew_NoPaths(t *testing.T) {
	_, err := multipath.New(&stubReader{}, []string{})
	if err == nil {
		t.Fatal("expected error for empty paths")
	}
}

func TestMerge_SinglePath(t *testing.T) {
	r := &stubReader{data: map[string]map[string]string{
		"secret/app": {"KEY": "val"},
	}}
	m, _ := multipath.New(r, []string{"secret/app"})
	got, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "val" {
		t.Errorf("expected val, got %q", got["KEY"])
	}
}

func TestMerge_LaterPathWins(t *testing.T) {
	r := &stubReader{data: map[string]map[string]string{
		"secret/base":    {"KEY": "base", "ONLY_BASE": "yes"},
		"secret/overlay": {"KEY": "overlay"},
	}}
	m, _ := multipath.New(r, []string{"secret/base", "secret/overlay"})
	got, err := m.Merge()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "overlay" {
		t.Errorf("expected overlay, got %q", got["KEY"])
	}
	if got["ONLY_BASE"] != "yes" {
		t.Errorf("expected base-only key to survive merge")
	}
}

func TestMerge_ReaderError(t *testing.T) {
	r := &stubReader{err: errors.New("vault unavailable")}
	m, _ := multipath.New(r, []string{"secret/app"})
	_, err := m.Merge()
	if err == nil {
		t.Fatal("expected error from reader")
	}
}

func TestNew_TrimsSlashes(t *testing.T) {
	r := &stubReader{data: map[string]map[string]string{
		"secret/app": {"X": "1"},
	}}
	m, _ := multipath.New(r, []string{"/secret/app/"})
	if m.Paths()[0] != "secret/app" {
		t.Errorf("expected trimmed path, got %q", m.Paths()[0])
	}
}
