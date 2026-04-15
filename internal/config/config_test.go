package config_test

import (
	"testing"

	"github.com/yourusername/vaultpull/internal/config"
)

var baseOverrides = map[string]string{
	"VAULT_ADDR":             "https://vault.example.com",
	"VAULT_TOKEN":            "s.abc123",
	"VAULTPULL_SECRET_PATH": "secret/myapp",
}

func TestLoad_ValidConfig(t *testing.T) {
	cfg, err := config.Load(baseOverrides)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.VaultAddr != "https://vault.example.com" {
		t.Errorf("unexpected VaultAddr: %q", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default OutputFile '.env', got %q", cfg.OutputFile)
	}
}

func TestLoad_NamespaceTrimsSlashes(t *testing.T) {
	ov := copyMap(baseOverrides)
	ov["VAULTPULL_NAMESPACE"] = "/team/backend/"
	cfg, err := config.Load(ov)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Namespace != "team/backend" {
		t.Errorf("expected trimmed namespace, got %q", cfg.Namespace)
	}
}

func TestLoad_MissingVaultAddr(t *testing.T) {
	ov := copyMap(baseOverrides)
	delete(ov, "VAULT_ADDR")
	_, err := config.Load(ov)
	if err == nil {
		t.Fatal("expected error for missing VAULT_ADDR")
	}
}

func TestLoad_MissingVaultToken(t *testing.T) {
	ov := copyMap(baseOverrides)
	delete(ov, "VAULT_TOKEN")
	_, err := config.Load(ov)
	if err == nil {
		t.Fatal("expected error for missing VAULT_TOKEN")
	}
}

func TestLoad_MissingSecretPath(t *testing.T) {
	ov := copyMap(baseOverrides)
	delete(ov, "VAULTPULL_SECRET_PATH")
	_, err := config.Load(ov)
	if err == nil {
		t.Fatal("expected error for missing VAULTPULL_SECRET_PATH")
	}
}

func TestLoad_AuditLogPassedThrough(t *testing.T) {
	ov := copyMap(baseOverrides)
	ov["VAULTPULL_AUDIT_LOG"] = "/tmp/audit.log"
	cfg, err := config.Load(ov)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.AuditLog != "/tmp/audit.log" {
		t.Errorf("expected AuditLog '/tmp/audit.log', got %q", cfg.AuditLog)
	}
}

func TestValidate_InvalidScheme(t *testing.T) {
	ov := copyMap(baseOverrides)
	ov["VAULT_ADDR"] = "ftp://vault.example.com"
	cfg, err := config.Load(ov)
	if err != nil {
		t.Fatalf("unexpected load error: %v", err)
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected validation error for ftp scheme")
	}
}

func copyMap(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
