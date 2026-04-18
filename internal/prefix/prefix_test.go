package prefix

import (
	"testing"
)

func TestNew_TrimsSlashes(t *testing.T) {
	s := New("/myapp/")
	if s.prefix != "myapp" {
		t.Fatalf("expected 'myapp', got %q", s.prefix)
	}
}

func TestStrip_EmptyPrefix_Unchanged(t *testing.T) {
	s := New("")
	input := map[string]string{"myapp/DB_HOST": "localhost"}
	out := s.Strip(input)
	if out["myapp/DB_HOST"] != "localhost" {
		t.Fatal("expected key to be unchanged")
	}
}

func TestStrip_RemovesPrefixWithSlash(t *testing.T) {
	s := New("myapp")
	input := map[string]string{"myapp/DB_HOST": "localhost", "myapp/DB_PORT": "5432"}
	out := s.Strip(input)
	if _, ok := out["DB_HOST"]; !ok {
		t.Fatal("expected DB_HOST after stripping prefix")
	}
	if _, ok := out["DB_PORT"]; !ok {
		t.Fatal("expected DB_PORT after stripping prefix")
	}
}

func TestStrip_KeyWithoutPrefix_Unchanged(t *testing.T) {
	s := New("myapp")
	input := map[string]string{"other/KEY": "val"}
	out := s.Strip(input)
	if out["other/KEY"] != "val" {
		t.Fatal("expected unrelated key to remain unchanged")
	}
}

func TestStrip_PreservesValues(t *testing.T) {
	s := New("svc")
	input := map[string]string{"svc/SECRET": "s3cr3t"}
	out := s.Strip(input)
	if out["SECRET"] != "s3cr3t" {
		t.Fatalf("expected value 's3cr3t', got %q", out["SECRET"])
	}
}

func TestHasPrefix_True(t *testing.T) {
	s := New("myapp")
	input := map[string]string{"myapp/KEY": "val"}
	if !s.HasPrefix(input) {
		t.Fatal("expected HasPrefix to return true")
	}
}

func TestHasPrefix_False(t *testing.T) {
	s := New("myapp")
	input := map[string]string{"other/KEY": "val"}
	if s.HasPrefix(input) {
		t.Fatal("expected HasPrefix to return false")
	}
}

func TestHasPrefix_EmptyPrefix_AlwaysFalse(t *testing.T) {
	s := New("")
	input := map[string]string{"anything": "val"}
	if s.HasPrefix(input) {
		t.Fatal("expected HasPrefix to return false for empty prefix")
	}
}
