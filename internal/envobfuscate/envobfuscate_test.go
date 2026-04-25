package envobfuscate

import (
	"strings"
	"testing"
)

func newObfuscator(salt string) *Obfuscator {
	return New([]byte(salt))
}

func TestApply_ReplacesKeys(t *testing.T) {
	o := newObfuscator("testsalt")
	src := map[string]string{"SECRET_KEY": "abc123"}
	out := o.Apply(src)
	if _, ok := out["SECRET_KEY"]; ok {
		t.Fatal("original key should not appear in output")
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(out))
	}
}

func TestApply_PreservesValues(t *testing.T) {
	o := newObfuscator("testsalt")
	src := map[string]string{"DB_PASS": "hunter2"}
	out := o.Apply(src)
	for _, v := range out {
		if v != "hunter2" {
			t.Fatalf("expected value %q, got %q", "hunter2", v)
		}
	}
}

func TestApply_DeterministicAlias(t *testing.T) {
	o := newObfuscator("salt1")
	out1 := o.Apply(map[string]string{"KEY": "v"})
	out2 := o.Apply(map[string]string{"KEY": "v"})
	var alias1, alias2 string
	for k := range out1 {
		alias1 = k
	}
	for k := range out2 {
		alias2 = k
	}
	if alias1 != alias2 {
		t.Fatalf("aliases differ: %q vs %q", alias1, alias2)
	}
}

func TestApply_DifferentSalts_DifferentAliases(t *testing.T) {
	o1 := newObfuscator("saltA")
	o2 := newObfuscator("saltB")
	var a1, a2 string
	for k := range o1.Apply(map[string]string{"KEY": "v"}) {
		a1 = k
	}
	for k := range o2.Apply(map[string]string{"KEY": "v"}) {
		a2 = k
	}
	if a1 == a2 {
		t.Fatal("different salts should produce different aliases")
	}
}

func TestReveal_MapsAliasToOriginal(t *testing.T) {
	o := newObfuscator("testsalt")
	src := map[string]string{"API_TOKEN": "tok"}
	rev := o.Reveal(src)
	for _, orig := range rev {
		if orig != "API_TOKEN" {
			t.Fatalf("expected API_TOKEN, got %q", orig)
		}
	}
}

func TestKeys_ReturnsSorted(t *testing.T) {
	o := newObfuscator("s")
	keys := o.Keys([]string{"Z_KEY", "A_KEY", "M_KEY"})
	for i := 1; i < len(keys); i++ {
		if keys[i] < keys[i-1] {
			t.Fatal("keys not sorted")
		}
	}
}

func TestAlias_IsHex(t *testing.T) {
	o := newObfuscator("hex-test")
	out := o.Apply(map[string]string{"SOME_KEY": "val"})
	for k := range out {
		if strings.ContainsAny(k, "ghijklmnopqrstuvwxyzGHIJKLMNOPQRSTUVWXYZ") {
			t.Fatalf("alias %q contains non-hex characters", k)
		}
		if len(k) != 16 {
			t.Fatalf("expected alias length 16, got %d", len(k))
		}
	}
}

func TestApply_EmptyMap_ReturnsEmpty(t *testing.T) {
	o := newObfuscator("s")
	out := o.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty map, got %d entries", len(out))
	}
}
