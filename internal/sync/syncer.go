package sync

import (
	"fmt"
	"log"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/env"
	"github.com/user/vaultpull/internal/vault"
)

// SecretReader abstracts vault.Client for testability.
type SecretReader interface {
	ReadSecrets(path string) (map[string]string, error)
}

// EnvWriter abstracts env.Writer for testability.
type EnvWriter interface {
	Write(secrets map[string]string) error
}

// Syncer orchestrates pulling secrets from Vault and writing them to a .env file.
type Syncer struct {
	cfg    *config.Config
	reader SecretReader
	writer EnvWriter
}

// New creates a Syncer wired with a real Vault client and env writer.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create vault client: %w", err)
	}

	writer, err := env.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create env writer: %w", err)
	}

	return &Syncer{
		cfg:    cfg,
		reader: client,
		writer: writer,
	}, nil
}

// NewWithDeps creates a Syncer with injected dependencies (useful for testing).
func NewWithDeps(cfg *config.Config, reader SecretReader, writer EnvWriter) *Syncer {
	return &Syncer{cfg: cfg, reader: reader, writer: writer}
}

// Run fetches secrets from Vault and writes them to the configured output file.
func (s *Syncer) Run() error {
	log.Printf("syncer: reading secrets from path %q (namespace: %q)", s.cfg.SecretPath, s.cfg.Namespace)

	secrets, err := s.reader.ReadSecrets(s.cfg.SecretPath)
	if err != nil {
		return fmt.Errorf("syncer: failed to read secrets: %w", err)
	}

	if len(secrets) == 0 {
		log.Println("syncer: no secrets found at path")
		return nil
	}

	log.Printf("syncer: writing %d secrets to %q", len(secrets), s.cfg.OutputFile)

	if err := s.writer.Write(secrets); err != nil {
		return fmt.Errorf("syncer: failed to write env file: %w", err)
	}

	log.Println("syncer: sync complete")
	return nil
}
