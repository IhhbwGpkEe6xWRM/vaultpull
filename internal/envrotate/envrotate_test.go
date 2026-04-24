package envrotate

import (
	"errors"
	"strings"
	"testing"
)

func identity(_, v string) (string, error) { return v, nil }

func upper(_, v string) (string, error) { return strings.ToUpper(v), nil }

func failOn(target string) RotateFunc {
	return func(k, v string) (string, error) {
		if k == target {
			return "", errors.New("forced error")
		}
		return v, nil
	}
}

func TestNew_NilFunc_ReturnsError(t *testing.T) {
	_, err := New(nil)
	if err == nil {
		t.Fatal("expected error for nil func")
	}
}

func TestApply_NilMap_ReturnsError(t *testing.T) {
	r, _ := New(identity)
	_, _, err := r.Apply(nil)
	if err == nil {
		t.Fatal("expected error for nil map")
	}
}

func TestApply_NoChanges_RotatedFalse(t *testing.T) {
	r, _ := New(identity)
	secrets := map[string]string{"KEY": "value"}
	out, results, err := r.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected value unchanged, got %q", out["KEY"])
	}
	if len(results) != 1 || results[0].Rotated {
		t.Errorf("expected Rotated=false, got %+v", results)
	}
}

func TestApply_ValueChanged_RotatedTrue(t *testing.T) {
	r, _ := New(upper)
	secrets := map[string]string{"DB_PASS": "secret"}
	out, results, err := r.Apply(secrets)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_PASS"] != "SECRET" {
		t.Errorf("expected SECRET, got %q", out["DB_PASS"])
	}
	if len(results) != 1 || !results[0].Rotated {
		t.Errorf("expected Rotated=true, got %+v", results)
	}
}

func TestApply_FuncError_Propagates(t *testing.T) {
	r, _ := New(failOn("BAD_KEY"))
	secrets := map[string]string{"BAD_KEY": "v"}
	_, _, err := r.Apply(secrets)
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "BAD_KEY") {
		t.Errorf("expected key in error, got %v", err)
	}
}

func TestApply_HashDiffersWhenValueChanges(t *testing.T) {
	r, _ := New(upper)
	secrets := map[string]string{"X": "abc"}
	_, results, _ := r.Apply(secrets)
	if results[0].OldHash == results[0].NewHash {
		t.Error("expected old and new hash to differ after rotation")
	}
}

func TestApply_EmptyMap_ReturnsEmptyResults(t *testing.T) {
	r, _ := New(identity)
	out, results, err := r.Apply(map[string]string{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 || len(results) != 0 {
		t.Error("expected empty output for empty input")
	}
}
