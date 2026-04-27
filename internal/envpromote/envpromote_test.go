package envpromote

import (
	"testing"
)

func TestPromote_AllKeys(t *testing.T) {
	p := New()
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}
	res, err := p.Promote(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 2 {
		t.Fatalf("expected 2 promoted, got %d", len(res.Promoted))
	}
	if dst["A"] != "1" || dst["B"] != "2" {
		t.Errorf("destination not updated correctly: %v", dst)
	}
}

func TestPromote_WithAllowList(t *testing.T) {
	p := New(WithAllowList([]string{"A"}))
	src := map[string]string{"A": "1", "B": "2"}
	dst := map[string]string{}
	res, err := p.Promote(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "A" {
		t.Errorf("expected only A promoted, got %v", res.Promoted)
	}
	if len(res.Skipped) != 1 || res.Skipped[0] != "B" {
		t.Errorf("expected B skipped, got %v", res.Skipped)
	}
	if _, ok := dst["B"]; ok {
		t.Error("B should not be in destination")
	}
}

func TestPromote_DryRun_DoesNotMutateDst(t *testing.T) {
	p := New(WithDryRun())
	src := map[string]string{"X": "secret"}
	dst := map[string]string{}
	res, err := p.Promote(src, dst)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.DryRun {
		t.Error("expected DryRun flag set")
	}
	if len(dst) != 0 {
		t.Errorf("dry-run must not mutate destination, got %v", dst)
	}
	if len(res.Promoted) != 1 || res.Promoted[0] != "X" {
		t.Errorf("expected X in promoted list, got %v", res.Promoted)
	}
}

func TestPromote_NilSource_ReturnsError(t *testing.T) {
	p := New()
	_, err := p.Promote(nil, map[string]string{})
	if err == nil {
		t.Fatal("expected error for nil source")
	}
}

func TestPromote_NilDest_ReturnsError(t *testing.T) {
	p := New()
	_, err := p.Promote(map[string]string{"A": "1"}, nil)
	if err == nil {
		t.Fatal("expected error for nil destination")
	}
}

func TestPromote_NilDest_DryRunOK(t *testing.T) {
	p := New(WithDryRun())
	_, err := p.Promote(map[string]string{"A": "1"}, nil)
	if err != nil {
		t.Fatalf("dry-run with nil dst should not error: %v", err)
	}
}

func TestPromote_PromotedKeysSorted(t *testing.T) {
	p := New()
	src := map[string]string{"Z": "z", "A": "a", "M": "m"}
	res, err := p.Promote(src, map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := []string{"A", "M", "Z"}
	for i, k := range expected {
		if res.Promoted[i] != k {
			t.Errorf("index %d: want %s got %s", i, k, res.Promoted[i])
		}
	}
}
