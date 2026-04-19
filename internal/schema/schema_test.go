package schema_test

import (
	"testing"

	"github.com/your-org/vaultpull/internal/schema"
)

func TestValidate_AllValid(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "PORT", Type: schema.TypeInt, Required: true},
		{Key: "DEBUG", Type: schema.TypeBool, Required: false},
		{Key: "NAME", Type: schema.TypeString, Required: true},
	})
	secrets := map[string]string{"PORT": "8080", "DEBUG": "true", "NAME": "app"}
	if v := s.Validate(secrets); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "TOKEN", Type: schema.TypeString, Required: true},
	})
	v := s.Validate(map[string]string{})
	if len(v) != 1 || v[0].Key != "TOKEN" {
		t.Fatalf("expected TOKEN violation, got %v", v)
	}
}

func TestValidate_InvalidInt(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "PORT", Type: schema.TypeInt, Required: true},
	})
	v := s.Validate(map[string]string{"PORT": "not-a-number"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %v", v)
	}
}

func TestValidate_InvalidBool(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "FLAG", Type: schema.TypeBool, Required: true},
	})
	v := s.Validate(map[string]string{"FLAG": "maybe"})
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %v", v)
	}
}

func TestValidate_PatternMatch(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "ENV", Type: schema.TypeString, Required: true, Pattern: "^(dev|staging|prod)$"},
	})
	if v := s.Validate(map[string]string{"ENV": "prod"}); len(v) != 0 {
		t.Fatalf("expected no violations, got %v", v)
	}
	if v := s.Validate(map[string]string{"ENV": "local"}); len(v) != 1 {
		t.Fatalf("expected pattern violation, got %v", v)
	}
}

func TestValidate_OptionalMissingIsOK(t *testing.T) {
	s := schema.New([]schema.Field{
		{Key: "OPTIONAL", Type: schema.TypeString, Required: false},
	})
	if v := s.Validate(map[string]string{}); len(v) != 0 {
		t.Fatalf("expected no violations for optional missing key, got %v", v)
	}
}

func TestViolation_Error(t *testing.T) {
	v := schema.Violation{Key: "FOO", Message: "required key is missing or empty"}
	if v.Error() == "" {
		t.Fatal("expected non-empty error string")
	}
}
