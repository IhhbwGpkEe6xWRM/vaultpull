// Package quota enforces per-path secret count limits during sync.
package quota

import "fmt"

// DefaultMaxSecrets is the default maximum number of secrets allowed per path.
const DefaultMaxSecrets = 500

// Config holds quota enforcement settings.
type Config struct {
	MaxSecretsPerPath int
	MaxTotalSecrets   int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		MaxSecretsPerPath: DefaultMaxSecrets,
		MaxTotalSecrets:   2000,
	}
}

// Enforcer checks secret maps against configured limits.
type Enforcer struct {
	cfg Config
}

// New returns an Enforcer using the provided Config.
func New(cfg Config) *Enforcer {
	if cfg.MaxSecretsPerPath <= 0 {
		cfg.MaxSecretsPerPath = DefaultMaxSecrets
	}
	if cfg.MaxTotalSecrets <= 0 {
		cfg.MaxTotalSecrets = 2000
	}
	return &Enforcer{cfg: cfg}
}

// CheckPath returns an error if the secret map for a single path exceeds the per-path limit.
func (e *Enforcer) CheckPath(path string, secrets map[string]string) error {
	if len(secrets) > e.cfg.MaxSecretsPerPath {
		return fmt.Errorf("quota: path %q has %d secrets, exceeds limit of %d",
			path, len(secrets), e.cfg.MaxSecretsPerPath)
	}
	return nil
}

// CheckTotal returns an error if the combined secret count exceeds the total limit.
func (e *Enforcer) CheckTotal(all map[string]map[string]string) error {
	total := 0
	for _, secrets := range all {
		total += len(secrets)
	}
	if total > e.cfg.MaxTotalSecrets {
		return fmt.Errorf("quota: total secret count %d exceeds limit of %d",
			total, e.cfg.MaxTotalSecrets)
	}
	return nil
}
