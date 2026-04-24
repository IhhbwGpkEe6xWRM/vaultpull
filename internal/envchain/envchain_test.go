package envchain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/example/vaultpull/internal/envchain"
)

// stubReader is a test double for envchain.Reader.
type stubReader struct {
	data map[string]string
	err  error
}

func (s *stubReader) Read(_ context.Context, _ string) (map[string]string, error) {
	if s.err != nil {
		return nil, s.err
	}
	out := make(map[string]string, len(s.data))
	for k, v := range s.data {
		out[k] = v
	}
	return out, nil
}

func TestResolve_SingleReader(t *testing.T) {
	r := &stubReader{data: map[string]string{"KEY": "value"}}
	chain := envchain.New(r)

	got, err := chain.Resolve(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "value" {
		t.Errorf("expected value, got %q", got["KEY"])
	}
}

func TestResolve_EarlierReaderWins(t *testing.T) {
	first := &stubReader{data: map[string]string{"KEY": "first"}}
	second := &stubReader{data: map[string]string{"KEY": "second", "OTHER": "x"}}
	chain := envchain.New(first, second)

	got, err := chain.Resolve(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "first" {
		t.Errorf("expected first reader to win, got %q", got["KEY"])
	}
	if got["OTHER"] != "x" {
		t.Errorf("expected OTHER from second reader, got %q", got["OTHER"])
	}
}

func TestResolve_EmptyValueNotOverwritten(t *testing.T) {
	// First reader has empty value; second reader should fill it in.
	first := &stubReader{data: map[string]string{"KEY": ""}}
	second := &stubReader{data: map[string]string{"KEY": "fallback"}}
	chain := envchain.New(first, second)

	got, err := chain.Resolve(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["KEY"] != "fallback" {
		t.Errorf("expected fallback for empty key, got %q", got["KEY"])
	}
}

func TestResolve_PropagatesError(t *testing.T) {
	boom := errors.New("vault unavailable")
	r := &stubReader{err: boom}
	chain := envchain.New(r)

	_, err := chain.Resolve(context.Background(), "secret/app")
	if !errors.Is(err, boom) {
		t.Errorf("expected boom error, got %v", err)
	}
}

func TestResolve_NoReaders_ReturnsEmpty(t *testing.T) {
	chain := envchain.New()

	got, err := chain.Resolve(context.Background(), "secret/app")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestLen_ReturnsReaderCount(t *testing.T) {
	chain := envchain.New(&stubReader{}, &stubReader{}, &stubReader{})
	if chain.Len() != 3 {
		t.Errorf("expected 3, got %d", chain.Len())
	}
}
