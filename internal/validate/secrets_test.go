package validate

import (
	"strings"
	"testing"
)

func TestSecrets_ValidMap(t *testing.T) {
	secrets := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
		"SECRET_KEY": "abc123",
	}
	result := Secrets(secrets)
	if !result.OK() {
		t.Fatalf("expected no issues, got: %v", result.Err())
	}
}

func TestSecrets_EmptyKey(t *testing.T) {
	secrets := map[string]string{
		"": "value",
	}
	result := Secrets(secrets)
	if result.OK() {
		t.Fatal("expected issue for empty key")
	}
	if !strings.Contains(result.Err().Error(), "must not be empty") {
		t.Errorf("unexpected error message: %v", result.Err())
	}
}

func TestSecrets_InvalidKeyCharacters(t *testing.T) {
	cases := []string{"my-key", "my key", "key.name", "key@host"}
	for _, k := range cases {
		t.Run(k, func(t *testing.T) {
			result := Secrets(map[string]string{k: "val"})
			if result.OK() {
				t.Fatalf("expected issue for key %q", k)
			}
			if !strings.Contains(result.Err().Error(), "invalid characters") {
				t.Errorf("unexpected message: %v", result.Err())
			}
		})
	}
}

func TestSecrets_ValueTooLong(t *testing.T) {
	secrets := map[string]string{
		"BIG_SECRET": strings.Repeat("x", maxValueLen+1),
	}
	result := Secrets(secrets)
	if result.OK() {
		t.Fatal("expected issue for oversized value")
	}
	if !strings.Contains(result.Err().Error(), "exceeds maximum length") {
		t.Errorf("unexpected message: %v", result.Err())
	}
}

func TestSecrets_MultipleIssues(t *testing.T) {
	secrets := map[string]string{
		"bad-key":    "fine",
		"ALSO_BAD-1": "fine",
	}
	result := Secrets(secrets)
	if len(result.Issues) < 2 {
		t.Fatalf("expected at least 2 issues, got %d", len(result.Issues))
	}
}

func TestResult_ErrNilWhenOK(t *testing.T) {
	result := Secrets(map[string]string{"VALID": "value"})
	if result.Err() != nil {
		t.Fatalf("expected nil error, got %v", result.Err())
	}
}

func TestIsValidKey(t *testing.T) {
	valid := []string{"A", "z", "_", "ABC_123", "a1_B2"}
	for _, k := range valid {
		if !isValidKey(k) {
			t.Errorf("expected %q to be valid", k)
		}
	}
	invalid := []string{"-", "a b", "a.b", "a@b"}
	for _, k := range invalid {
		if isValidKey(k) {
			t.Errorf("expected %q to be invalid", k)
		}
	}
}
