package envflow_test

import (
	"errors"
	"testing"

	"github.com/example/vaultpull/internal/envflow"
)

func identity(m map[string]string) (map[string]string, error) {
	return m, nil
}

func addKey(key, val string) func(map[string]string) (map[string]string, error) {
	return func(m map[string]string) (map[string]string, error) {
		out := make(map[string]string, len(m)+1)
		for k, v := range m {
			out[k] = v
		}
		out[key] = val
		return out, nil
	}
}

func failStage(msg string) func(map[string]string) (map[string]string, error) {
	return func(m map[string]string) (map[string]string, error) {
		return nil, errors.New(msg)
	}
}

func TestNew_EmptyStageName(t *testing.T) {
	_, err := envflow.New([]envflow.Stage{{Name: "", Apply: identity}})
	if err == nil {
		t.Fatal("expected error for empty stage name")
	}
}

func TestNew_NilApply(t *testing.T) {
	_, err := envflow.New([]envflow.Stage{{Name: "s", Apply: nil}})
	if err == nil {
		t.Fatal("expected error for nil Apply")
	}
}

func TestRun_EmptyPipeline(t *testing.T) {
	p, _ := envflow.New(nil)
	in := map[string]string{"A": "1"}
	out, results, err := p.Run(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
	if out["A"] != "1" {
		t.Errorf("expected A=1, got %q", out["A"])
	}
}

func TestRun_SingleStage(t *testing.T) {
	p, _ := envflow.New([]envflow.Stage{
		{Name: "add", Apply: addKey("NEW", "val")},
	})
	out, results, err := p.Run(map[string]string{"X": "1"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 || results[0].Stage != "add" {
		t.Errorf("unexpected results: %+v", results)
	}
	if out["NEW"] != "val" {
		t.Errorf("expected NEW=val, got %q", out["NEW"])
	}
}

func TestRun_StopsOnError(t *testing.T) {
	p, _ := envflow.New([]envflow.Stage{
		{Name: "ok", Apply: identity},
		{Name: "boom", Apply: failStage("exploded")},
		{Name: "never", Apply: identity},
	})
	_, results, err := p.Run(map[string]string{})
	if err == nil {
		t.Fatal("expected error")
	}
	if len(results) != 2 {
		t.Errorf("expected 2 results (ok + boom), got %d", len(results))
	}
	if results[1].Err == nil {
		t.Error("expected boom stage to carry error")
	}
}

func TestRun_ResultKeyCount(t *testing.T) {
	p, _ := envflow.New([]envflow.Stage{
		{Name: "add2", Apply: func(m map[string]string) (map[string]string, error) {
			return map[string]string{"A": "1", "B": "2"}, nil
		}},
	})
	_, results, err := p.Run(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if results[0].Keys != 2 {
		t.Errorf("expected Keys=2, got %d", results[0].Keys)
	}
}
