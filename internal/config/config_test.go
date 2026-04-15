package config

import (
	"testing"
)

func TestLoad_ValidConfig(t *testing.T) {
	overrides := map[string]string{
		"VAULT_ADDR":            "http://127.0.0.1:8200",
		"VAULT_TOKEN":           "s.testtoken",
		"VAULTPULL_SECRET_PATH": "secret/data/myapp",
		"VAULT_NAMESPACE":       "team/backend",
	}

	cfg, err := Load(overrides)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if cfg.VaultAddr != "http://127.0.0.1:8200" {
		t.Errorf("unexpected VaultAddr: %s", cfg.VaultAddr)
	}
	if cfg.OutputFile != ".env" {
		t.Errorf("expected default OutputFile '.env', got: %s", cfg.OutputFile)
	}
	if cfg.Namespace != "team/backend" {
		t.Errorf("unexpected Namespace: %s", cfg.Namespace)
	}
}

func TestLoad_NamespaceTrimsSlashes(t *testing.T) {
	overrides := map[string]string{
		"VAULT_ADDR":            "http://127.0.0.1:8200",
		"VAULT_TOKEN":           "s.testtoken",
		"VAULTPULL_SECRET_PATH": "secret/data/myapp",
		"VAULT_NAMESPACE":       "/team/backend/",
	}

	cfg, err := Load(overrides)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.Namespace != "team/backend" {
		t.Errorf("expected trimmed namespace, got: %s", cfg.Namespace)
	}
}

func TestLoad_MissingVaultAddr(t *testing.T) {
	overrides := map[string]string{
		"VAULT_TOKEN":           "s.testtoken",
		"VAULTPULL_SECRET_PATH": "secret/data/myapp",
	}

	_, err := Load(overrides)
	if err == nil {
		t.Fatal("expected error for missing VAULT_ADDR")
	}
}

func TestLoad_MissingVaultToken(t *testing.T) {
	overrides := map[string]string{
		"VAULT_ADDR":            "http://127.0.0.1:8200",
		"VAULTPULL_SECRET_PATH": "secret/data/myapp",
	}

	_, err := Load(overrides)
	if err == nil {
		t.Fatal("expected error for missing VAULT_TOKEN")
	}
}

func TestLoad_MissingSecretPath(t *testing.T) {
	overrides := map[string]string{
		"VAULT_ADDR":  "http://127.0.0.1:8200",
		"VAULT_TOKEN": "s.testtoken",
	}

	_, err := Load(overrides)
	if err == nil {
		t.Fatal("expected error for missing VAULTPULL_SECRET_PATH")
	}
}

func TestLoad_CustomOutputFile(t *testing.T) {
	overrides := map[string]string{
		"VAULT_ADDR":            "http://127.0.0.1:8200",
		"VAULT_TOKEN":           "s.testtoken",
		"VAULTPULL_SECRET_PATH": "secret/data/myapp",
		"VAULTPULL_OUTPUT":      ".env.local",
	}

	cfg, err := Load(overrides)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if cfg.OutputFile != ".env.local" {
		t.Errorf("expected OutputFile '.env.local', got: %s", cfg.OutputFile)
	}
}
