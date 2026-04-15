// Package config loads and validates vaultpull runtime configuration.
package config

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// Config holds all configuration required for a vaultpull run.
type Config struct {
	VaultAddr  string
	VaultToken string
	SecretPath string
	Namespace  string
	OutputFile string
	Quiet      bool
	AuditLog   string
}

// Load reads configuration from environment variables, applying any
// overrides supplied via the opts map (key = env var name, value = override).
func Load(overrides map[string]string) (*Config, error) {
	get := func(key string) string {
		return getEnvOrOverride(key, overrides)
	}

	cfg := &Config{
		VaultAddr:  get("VAULT_ADDR"),
		VaultToken: get("VAULT_TOKEN"),
		SecretPath: get("VAULTPULL_SECRET_PATH"),
		Namespace:  strings.Trim(get("VAULTPULL_NAMESPACE"), "/"),
		OutputFile: get("VAULTPULL_OUTPUT_FILE"),
		AuditLog:   get("VAULTPULL_AUDIT_LOG"),
		Quiet:      get("VAULTPULL_QUIET") == "true",
	}

	if cfg.VaultAddr == "" {
		return nil, errors.New("config: VAULT_ADDR is required")
	}
	if cfg.VaultToken == "" {
		return nil, errors.New("config: VAULT_TOKEN is required")
	}
	if cfg.SecretPath == "" {
		return nil, errors.New("config: VAULTPULL_SECRET_PATH is required")
	}
	if cfg.OutputFile == "" {
		cfg.OutputFile = ".env"
	}

	return cfg, nil
}

// Validate performs additional semantic validation beyond presence checks.
func (c *Config) Validate() error {
	if !strings.HasPrefix(c.VaultAddr, "http://") && !strings.HasPrefix(c.VaultAddr, "https://") {
		return fmt.Errorf("config: VAULT_ADDR must start with http:// or https://, got %q", c.VaultAddr)
	}
	return nil
}

func getEnvOrOverride(key string, overrides map[string]string) string {
	if v, ok := overrides[key]; ok {
		return v
	}
	return os.Getenv(key)
}
