package batch_test

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/your-org/vaultpull/internal/batch"
)

type mockReader struct {
	data map[string]map[string]string
	err  error
}

func (m *mockReader) ReadSecrets(_ context.Context, path string) (map[string]string, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.data[path], nil
}

func TestFetchAll_ReturnsAllResults(t *testing.T) {
	r := &mockReader{data: map[string]map[string]string{
		"a": {"K1": "v1"},
		"b": {"K2": "v2"},
	}}
	f := batch.New(r, 2)
	results := f.FetchAll(context.Background(), []string{"a", "b"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	for _, res := range results {
		if res.Err != nil {
			t.Errorf("unexpected error for %q: %v", res.Path, res.Err)
		}
	}
}

func TestFetchAll_PropagatesError(t *testing.T) {
	r := &mockReader{err: errors.New("vault down")}
	f := batch.New(r, 2)
	results := f.FetchAll(context.Background(), []string{"a"})
	if results[0].Err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFetchAll_CancelledContext(t *testing.T) {
	r := &mockReader{data: map[string]map[string]string{}}
	f := batch.New(r, 1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	results := f.FetchAll(ctx, []string{"a", "b"})
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
}

func TestMerge_CombinesSecrets(t *testing.T) {
	results := []batch.Result{
		{Path: "a", Secrets: map[string]string{"K1": "v1"}},
		{Path: "b", Secrets: map[string]string{"K2": "v2"}},
	}
	out, err := batch.Merge(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	keys := make([]string, 0, len(out))
	for k := range out {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if len(keys) != 2 || keys[0] != "K1" || keys[1] != "K2" {
		t.Errorf("unexpected keys: %v", keys)
	}
}

func TestMerge_ErrorOnFailedResult(t *testing.T) {
	results := []batch.Result{
		{Path: "a", Secrets: map[string]string{"K1": "v1"}},
		{Path: "b", Err: errors.New("read failed")},
	}
	_, err := batch.Merge(results)
	if err == nil {
		t.Fatal("expected error from Merge")
	}
}

func TestNew_DefaultConcurrency(t *testing.T) {
	r := &mockReader{data: map[string]map[string]string{"x": {"A": "1"}}}
	f := batch.New(r, 0)
	results := f.FetchAll(context.Background(), []string{"x"})
	if results[0].Err != nil {
		t.Errorf("unexpected error: %v", results[0].Err)
	}
}
