package sync

import (
	"context"
	"fmt"

	"github.com/your-org/vaultpull/internal/backup"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/filter"
	"github.com/your-org/vaultpull/internal/output"
	"github.com/your-org/vaultpull/internal/validate"
)

// SecretsReader fetches secrets from Vault.
type SecretsReader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Syncer orchestrates pulling secrets from Vault and writing them locally.
type Syncer struct {
	cfg     *config.Config
	reader  SecretsReader
	writer  *env.Writer
	out     *output.Formatter
	matcher *filter.Matcher
	backups *backup.Store
}

// New creates a Syncer wired to real dependencies.
func New(cfg *config.Config, reader SecretsReader) (*Syncer, error) {
	w, err := env.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, err
	}
	return NewWithDeps(cfg, reader, w, output.New(cfg.Quiet), nil)
}

// NewWithDeps creates a Syncer with explicit dependencies (useful in tests).
func NewWithDeps(cfg *config.Config, reader SecretsReader, writer *env.Writer, out *output.Formatter, bk *backup.Store) (*Syncer, error) {
	m, err := filter.NewMatcher(cfg.Namespace)
	if err != nil {
		return nil, fmt.Errorf("syncer: namespace matcher: %w", err)
	}
	return &Syncer{cfg: cfg, reader: reader, writer: writer, out: out, matcher: m, backups: bk}, nil
}

// Run executes the sync: read → filter → validate → backup → write.
func (s *Syncer) Run(ctx context.Context) error {
	s.out.Info("Fetching secrets from %s", s.cfg.SecretPath)

	secrets, err := s.reader.ReadSecrets(ctx, s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("read secrets: %w", err)
	}

	filtered := make(map[string]string)
	for k, v := range secrets {
		if s.matcher.Match(k) {
			filtered[k] = v
		}
	}

	if issues := validate.Secrets(filtered); len(issues) > 0 {
		for _, iss := range issues {
			s.out.Warn("%s", iss)
		}
	}

	if s.backups != nil {
		if bak, berr := s.backups.Save(s.cfg.OutputFile); berr != nil {
			s.out.Warn("backup failed: %v", berr)
		} else if bak != "" {
			s.out.Info("Backup saved to %s", bak)
		}
	}

	if err := s.writer.Write(filtered); err != nil {
		return fmt.Errorf("write env: %w", err)
	}

	s.out.Success("Wrote %d secrets to %s", len(filtered), s.cfg.OutputFile)
	return nil
}
