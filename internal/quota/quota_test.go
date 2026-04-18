package quota_test

import (
	"strings"
	"testing"

	"github.com/yourusername/vaultpull/internal/quota"
)

func makeSecrets(n int) map[string]string {
	m := make(map[string]string, n)
	for i := 0; i < n; i++ {
		key := fmt.Sprintf("KEY_%d", i)
		m[key] = "value"
	}
	return m
}

import "fmt"

func TestCheckPath_UnderLimit(t *testing.T) {
	e := quota.New(quota.DefaultConfig())
	if err := e.CheckPath("secret/app", makeSecrets(10)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckPath_ExceedsLimit(t *testing.T) {
	cfg := quota.Config{MaxSecretsPerPath: 5, MaxTotalSecrets: 1000}
	e := quota.New(cfg)
	err := e.CheckPath("secret/app", makeSecrets(6))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "secret/app") {
		t.Errorf("error should mention path, got: %v", err)
	}
}

func TestCheckPath_ExactLimit_OK(t *testing.T) {
	cfg := quota.Config{MaxSecretsPerPath: 5, MaxTotalSecrets: 1000}
	e := quota.New(cfg)
	if err := e.CheckPath("secret/app", makeSecrets(5)); err != nil {
		t.Fatalf("unexpected error at exact limit: %v", err)
	}
}

func TestCheckTotal_UnderLimit(t *testing.T) {
	e := quota.New(quota.DefaultConfig())
	all := map[string]map[string]string{
		"path/a": makeSecrets(10),
		"path/b": makeSecrets(20),
	}
	if err := e.CheckTotal(all); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCheckTotal_ExceedsLimit(t *testing.T) {
	cfg := quota.Config{MaxSecretsPerPath: 1000, MaxTotalSecrets: 15}
	e := quota.New(cfg)
	all := map[string]map[string]string{
		"path/a": makeSecrets(10),
		"path/b": makeSecrets(10),
	}
	err := e.CheckTotal(all)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "20") {
		t.Errorf("error should mention total count, got: %v", err)
	}
}

func TestNew_ZeroLimits_UsesDefaults(t *testing.T) {
	e := quota.New(quota.Config{})
	if err := e.CheckPath("p", makeSecrets(1)); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
