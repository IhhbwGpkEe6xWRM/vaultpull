package masker_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/masker"
)

func TestMask_FullMasker_ReplacesValue(t *testing.T) {
	m := masker.New()
	if got := m.Mask("supersecret"); got != masker.DefaultMask {
		t.Errorf("expected %q, got %q", masker.DefaultMask, got)
	}
}

func TestMask_FullMasker_EmptyValue(t *testing.T) {
	m := masker.New()
	if got := m.Mask(""); got != "" {
		t.Errorf("expected empty string, got %q", got)
	}
}

func TestMask_PartialMasker_ShortValue(t *testing.T) {
	m := masker.NewPartial()
	// value shorter than 2*visibleChars should be fully masked
	if got := m.Mask("tiny"); got != masker.DefaultMask {
		t.Errorf("expected full mask for short value, got %q", got)
	}
}

func TestMask_PartialMasker_LongValue(t *testing.T) {
	m := masker.NewPartial()
	value := "abcd1234efgh5678"
	got := m.Mask(value)

	if len(got) == 0 {
		t.Fatal("expected non-empty masked value")
	}
	if got == value {
		t.Error("expected value to be partially masked")
	}
	// Should start with first 4 chars and end with last 4 chars
	if got[:4] != "abcd" {
		t.Errorf("expected prefix %q, got %q", "abcd", got[:4])
	}
	if got[len(got)-4:] != "5678" {
		t.Errorf("expected suffix %q, got %q", "5678", got[len(got)-4:])
	}
}

func TestMaskMap_RedactsAllValues(t *testing.T) {
	m := masker.New()
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}
	masked := m.MaskMap(secrets)

	for k, v := range masked {
		if v != masker.DefaultMask {
			t.Errorf("key %q: expected %q, got %q", k, masker.DefaultMask, v)
		}
	}
	if len(masked) != len(secrets) {
		t.Errorf("expected %d entries, got %d", len(secrets), len(masked))
	}
}

func TestMaskMap_DoesNotMutateOriginal(t *testing.T) {
	m := masker.New()
	secrets := map[string]string{"KEY": "value"}
	_ = m.MaskMap(secrets)
	if secrets["KEY"] != "value" {
		t.Error("original map was mutated")
	}
}

func TestContainsSensitive_DetectsSecret(t *testing.T) {
	secrets := map[string]string{"TOKEN": "s3cr3t"}
	if !masker.ContainsSensitive("the token is s3cr3t here", secrets) {
		t.Error("expected sensitive content to be detected")
	}
}

func TestContainsSensitive_NoMatch(t *testing.T) {
	secrets := map[string]string{"TOKEN": "s3cr3t"}
	if masker.ContainsSensitive("nothing sensitive here", secrets) {
		t.Error("expected no sensitive content")
	}
}

func TestContainsSensitive_EmptyValue_NoFalsePositive(t *testing.T) {
	secrets := map[string]string{"EMPTY": ""}
	if masker.ContainsSensitive("anything", secrets) {
		t.Error("empty secret value should not trigger a match")
	}
}
