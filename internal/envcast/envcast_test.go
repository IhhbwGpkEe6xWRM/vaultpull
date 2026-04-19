package envcast_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/envcast"
)

func newCaster() *envcast.Caster {
	return envcast.New()
}

func TestString_Present(t *testing.T) {
	c := newCaster()
	m := map[string]string{"KEY": "hello"}
	v, err := c.String(m, "KEY")
	if err != nil || v != "hello" {
		t.Fatalf("expected hello, got %q %v", v, err)
	}
}

func TestString_Missing(t *testing.T) {
	c := newCaster()
	_, err := c.String(map[string]string{}, "MISSING")
	if err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestInt_Valid(t *testing.T) {
	c := newCaster()
	m := map[string]string{"PORT": "8080"}
	v, err := c.Int(m, "PORT")
	if err != nil || v != 8080 {
		t.Fatalf("expected 8080, got %d %v", v, err)
	}
}

func TestInt_Invalid(t *testing.T) {
	c := newCaster()
	m := map[string]string{"PORT": "abc"}
	_, err := c.Int(m, "PORT")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestBool_TrueVariants(t *testing.T) {
	c := newCaster()
	for _, val := range []string{"true", "1", "yes", "TRUE", "Yes"} {
		m := map[string]string{"FLAG": val}
		v, err := c.Bool(m, "FLAG")
		if err != nil || !v {
			t.Errorf("expected true for %q, got %v %v", val, v, err)
		}
	}
}

func TestBool_FalseVariants(t *testing.T) {
	c := newCaster()
	for _, val := range []string{"false", "0", "no", "FALSE"} {
		m := map[string]string{"FLAG": val}
		v, err := c.Bool(m, "FLAG")
		if err != nil || v {
			t.Errorf("expected false for %q, got %v %v", val, v, err)
		}
	}
}

func TestBool_Invalid(t *testing.T) {
	c := newCaster()
	m := map[string]string{"FLAG": "maybe"}
	_, err := c.Bool(m, "FLAG")
	if err == nil {
		t.Fatal("expected error for invalid bool")
	}
}

func TestFloat_Valid(t *testing.T) {
	c := newCaster()
	m := map[string]string{"RATIO": "3.14"}
	v, err := c.Float(m, "RATIO")
	if err != nil || v != 3.14 {
		t.Fatalf("expected 3.14, got %f %v", v, err)
	}
}

func TestFloat_Invalid(t *testing.T) {
	c := newCaster()
	m := map[string]string{"RATIO": "not-a-float"}
	_, err := c.Float(m, "RATIO")
	if err == nil {
		t.Fatal("expected parse error")
	}
}
