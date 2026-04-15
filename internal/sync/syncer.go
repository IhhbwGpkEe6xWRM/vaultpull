package sync

import (
	"fmt"

	"github.com/example/vaultpull/internal/audit"
	"github.com/example/vaultpull/internal/cache"
	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/env"
	"github.com/example/vaultpull/internal/filter"
	"github.com/example/vaultpull/internal/output"
)

// SecretReader fetches secrets from a remote source.
type SecretReader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// Syncer orchestrates reading secrets and writing .env files.
type Syncer struct {
	cfg     *config.Config
	reader  SecretReader
	writer  *env.Writer
	matcher *filter.Matcher
	fmt     *output.Formatter
	audit   *audit.Logger
	cache   *cache.Store
}

// New builds a Syncer from a validated Config using real implementations.
func New(cfg *config.Config) (*Syncer, error) {
	return NewWithDeps(cfg, nil, nil, nil, nil, nil)
}

// NewWithDeps constructs a Syncer with injectable dependencies (useful in tests).
func NewWithDeps(
	cfg *config.Config,
	reader SecretReader,
	writer *env.Writer,
	fmt *output.Formatter,
	log *audit.Logger,
	store *cache.Store,
) (*Syncer, error) {
	matcher, err := filter.NewMatcher(cfg.Namespace)
	if err != nil {
		return nil, fmt.Errorf("filter: %w", err)
	}
	return &Syncer{
		cfg:     cfg,
		reader:  reader,
		writer:  writer,
		matcher: matcher,
		fmt:     fmt,
		audit:   log,
		cache:   store,
	}, nil
}

// Run fetches secrets and writes the .env file, skipping if the cache is fresh.
func (s *Syncer) Run() error {
	if !s.matcher.Match(s.cfg.SecretPath) {
		s.fmt.Info("skipping %s (namespace filter)", s.cfg.SecretPath)
		return nil
	}

	secrets, err := s.reader.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("read secrets: %w", err)
	}

	if s.cache != nil && s.cache.IsFresh(s.cfg.SecretPath, secrets) {
		s.fmt.Info("cache hit for %s — skipping write", s.cfg.SecretPath)
		return nil
	}

	if err := s.writer.Write(s.cfg.OutputFile, secrets); err != nil {
		return fmt.Errorf("write env: %w", err)
	}

	if s.cache != nil {
		s.cache.Set(s.cfg.SecretPath, secrets)
		_ = s.cache.Save()
	}

	s.audit.Record(s.cfg.SecretPath, "synced", len(secrets))
	s.fmt.Success("wrote %d secrets to %s", len(secrets), s.cfg.OutputFile)
	return nil
}
