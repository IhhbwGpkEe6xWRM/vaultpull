package envsubst

import (
	"testing"
)

func newSubstitutor(opts ...Option) *Substitutor {
	return New(opts...)
}

func TestApply_NoReferences(t *testing.T) {
	s := newSubstitutor()
	src := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := s.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected unchanged values, got %v", out)
	}
}

func TestApply_BraceStyle(t *testing.T) {
	s := newSubstitutor()
	src := map[string]string{"HOST": "localhost", "ADDR": "${HOST}:5432"}
	out, err := s.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["ADDR"] != "localhost:5432" {
		t.Errorf("expected 'localhost:5432', got %q", out["ADDR"])
	}
}

func TestApply_NoBraceStyle(t *testing.T) {
	s := newSubstitutor()
	src := map[string]string{"PORT": "8080", "URL": "http://host:$PORT"}
	out, err := s.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["URL"] != "http://host:8080" {
		t.Errorf("expected 'http://host:8080', got %q", out["URL"])
	}
}

func TestApply_ChainedReferences(t *testing.T) {
	s := newSubstitutor()
	src := map[string]string{"A": "hello", "B": "${A} world", "C": "${B}!"}
	out, err := s.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["C"] != "hello world!" {
		t.Errorf("expected 'hello world!', got %q", out["C"])
	}
}

func TestApply_MissingKey_BecomesEmpty(t *testing.T) {
	s := newSubstitutor()
	src := map[string]string{"VAL": "${MISSING}"}
	out, err := s.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["VAL"] != "" {
		t.Errorf("expected empty string, got %q", out["VAL"])
	}
}

func TestApply_WithOSFallback(t *testing.T) {
	os := map[string]string{"REGION": "us-east-1"}
	s := newSubstitutor(WithOSEnv(os))
	src := map[string]string{"BUCKET": "my-bucket-${REGION}"}
	out, err := s.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["BUCKET"] != "my-bucket-us-east-1" {
		t.Errorf("expected 'my-bucket-us-east-1', got %q", out["BUCKET"])
	}
}

func TestApply_SrcTakesPrecedenceOverOS(t *testing.T) {
	os := map[string]string{"ENV": "production"}
	s := newSubstitutor(WithOSEnv(os))
	src := map[string]string{"ENV": "staging", "TAG": "${ENV}"}
	out, err := s.Apply(src)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["TAG"] != "staging" {
		t.Errorf("expected 'staging', got %q", out["TAG"])
	}
}

func TestContainsReferences_True(t *testing.T) {
	if !ContainsReferences("${FOO}") {
		t.Error("expected true for '${FOO}'")
	}
	if !ContainsReferences("$BAR") {
		t.Error("expected true for '$BAR'")
	}
}

func TestContainsReferences_False(t *testing.T) {
	if ContainsReferences("plain value") {
		t.Error("expected false for plain value")
	}
}
