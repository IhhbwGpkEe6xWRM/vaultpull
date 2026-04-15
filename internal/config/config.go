package config

import (
	"errors"
	"os"
	"strings"
)

// Config holds all configuration for vaultpull.
type Config struct {
	VaultAddr  string
	VaultToken string
	Namespace  string
	OutputFile string
	SecretPath string
}

// Load reads configuration from environment variables with optional overrides.
func Load(overrides map[string]string) (*Config, error) {
	cfg := &Config{
		VaultAddr:  getEnvOrOverride("VAULT_ADDR", overrides),
		VaultToken: getEnvOrOverride("VAULT_TOKEN", overrides),
		Namespace:  getEnvOrOverride("VAULT_NAMESPACE", overrides),
		OutputFile: getEnvOrOverride("VAULTPULL_OUTPUT", overrides),
		SecretPath: getEnvOrOverride("VAULTPULL_SECRET_PATH", overrides),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	if cfg.OutputFile == "" {
		cfg.OutputFile = ".env"
	}

	cfg.Namespace = strings.Trim(cfg.Namespace, "/")

	return cfg, nil
}

func (c *Config) validate() error {
	if c.VaultAddr == "" {
		return errors.New("VAULT_ADDR is required")
	}
	if c.VaultToken == "" {
		return errors.New("VAULT_TOKEN is required")
	}
	if c.SecretPath == "" {
		return errors.New("VAULTPULL_SECRET_PATH is required")
	}
	return nil
}

func getEnvOrOverride(key string, overrides map[string]string) string {
	if overrides != nil {
		if val, ok := overrides[key]; ok {
			return val
		}
	}
	return os.Getenv(key)
}
