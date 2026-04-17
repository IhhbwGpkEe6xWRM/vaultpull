package sync_test

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/backup"
	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/output"
	"github.com/your-org/vaultpull/internal/sync"
)

type mockReader struct {
	secrets map[string]string
	err     error
}

func (m *mockReader) ReadSecrets(_ context.Context, _ string) (map[string]string, error) {
	return m.secrets, m.err
}

func baseConfig(t *testing.T) *config.Config {
	t.Helper()
	dir, _ := os.MkdirTemp("", "syncer-test-*")
	t.Cleanup(func() { os.RemoveAll(dir) })
	return &config.Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "tok",
		SecretPath: "secret/app",
		OutputFile: filepath.Join(dir, ".env"),
		Quiet:      true,
	}
}

func TestRun_Success(t *testing.T) {
	cfg := baseConfig(t)
	reader := &mockReader{secrets: map[string]string{"KEY": "val"}}
	w, _ := env.NewWriter(cfg.OutputFile)
	s, _ := sync.NewWithDeps(cfg, reader, w, output.New(true), nil)

	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	data, _ := os.ReadFile(cfg.OutputFile)
	if len(data) == 0 {
		t.Error("expected non-empty output file")
	}
}

func TestRun_ReaderError(t *testing.T) {
	cfg := baseConfig(t)
	reader := &mockReader{err: errors.New("vault down")}
	w, _ := env.NewWriter(cfg.OutputFile)
	s, _ := sync.NewWithDeps(cfg, reader, w, output.New(true), nil)

	if err := s.Run(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestRun_WithBackup(t *testing.T) {
	cfg := baseConfig(t)
	// Pre-create output file so backup has something to copy.
	_ = os.WriteFile(cfg.OutputFile, []byte("OLD=1\n"), 0600)

	bdir, _ := os.MkdirTemp("", "bak-*")
	t.Cleanup(func() { os.RemoveAll(bdir) })
	bk, _ := backup.NewStore(bdir)

	reader := &mockReader{secrets: map[string]string{"NEW": "2"}}
	w, _ := env.NewWriter(cfg.OutputFile)
	s, _ := sync.NewWithDeps(cfg, reader, w, output.New(true), bk)

	if err := s.Run(context.Background()); err != nil {
		t.Fatalf("Run: %v", err)
	}
	files, _ := bk.List(cfg.OutputFile)
	if len(files) != 1 {
		t.Errorf("expected 1 backup, got %d", len(files))
	}
}
