package envhash_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envhash"
)

func newHasher(t *testing.T) *envhash.Hasher {
	t.Helper()
	return envhash.New()
}

func TestSum_EmptyMap_ReturnsEmptySHA256(t *testing.T) {
	h := newHasher(t)
	got := h.Sum(nil)
	if len(got) != 64 {
		t.Fatalf("expected 64-char hex digest, got %d chars: %s", len(got), got)
	}
}

func TestSum_DeterministicForSameInput(t *testing.T) {
	h := newHasher(t)
	secrets := map[string]string{"FOO": "bar", "BAZ": "qux"}
	if h.Sum(secrets) != h.Sum(secrets) {
		t.Fatal("expected identical hashes for same input")
	}
}

func TestSum_OrderIndependent(t *testing.T) {
	h := newHasher(t)
	a := map[string]string{"A": "1", "B": "2", "C": "3"}
	b := map[string]string{"C": "3", "A": "1", "B": "2"}
	if h.Sum(a) != h.Sum(b) {
		t.Fatalf("expected order-independent hash: %s != %s", h.Sum(a), h.Sum(b))
	}
}

func TestSum_DifferentValues_DifferentHash(t *testing.T) {
	h := newHasher(t)
	a := map[string]string{"KEY": "value1"}
	b := map[string]string{"KEY": "value2"}
	if h.Sum(a) == h.Sum(b) {
		t.Fatal("expected different hashes for different values")
	}
}

func TestSum_DifferentKeys_DifferentHash(t *testing.T) {
	h := newHasher(t)
	a := map[string]string{"KEY_A": "val"}
	b := map[string]string{"KEY_B": "val"}
	if h.Sum(a) == h.Sum(b) {
		t.Fatal("expected different hashes for different keys")
	}
}

func TestEqual_IdenticalMaps_ReturnsTrue(t *testing.T) {
	h := newHasher(t)
	secrets := map[string]string{"X": "y"}
	if !h.Equal(secrets, secrets) {
		t.Fatal("expected Equal to return true for identical maps")
	}
}

func TestEqual_DifferentMaps_ReturnsFalse(t *testing.T) {
	h := newHasher(t)
	a := map[string]string{"X": "y"}
	b := map[string]string{"X": "z"}
	if h.Equal(a, b) {
		t.Fatal("expected Equal to return false for different maps")
	}
}

func TestNewWithSeparator_EmptySep_ReturnsError(t *testing.T) {
	_, err := envhash.NewWithSeparator("")
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestNewWithSeparator_CustomSep_ProducesDifferentHash(t *testing.T) {
	defaultH := envhash.New()
	customH, err := envhash.NewWithSeparator(":")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	secrets := map[string]string{"KEY": "value"}
	if defaultH.Sum(secrets) == customH.Sum(secrets) {
		t.Fatal("expected different hashes for different separators")
	}
}
