package sync_test

import (
	"errors"
	"testing"

	"github.com/user/vaultpull/internal/config"
	"github.com/user/vaultpull/internal/sync"
)

// mockReader implements SecretReader.
type mockReader struct {
	secrets map[string]string
	err     error
}

func (m *mockReader) ReadSecrets(_ string) (map[string]string, error) {
	return m.secrets, m.err
}

// mockWriter implements EnvWriter.
type mockWriter struct {
	written map[string]string
	err     error
}

func (m *mockWriter) Write(secrets map[string]string) error {
	if m.err != nil {
		return m.err
	}
	m.written = secrets
	return nil
}

func baseConfig() *config.Config {
	return &config.Config{
		VaultAddr:  "http://127.0.0.1:8200",
		VaultToken: "test-token",
		SecretPath: "secret/myapp",
		OutputFile: ".env",
	}
}

func TestRun_Success(t *testing.T) {
	reader := &mockReader{secrets: map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}}
	writer := &mockWriter{}
	s := sync.NewWithDeps(baseConfig(), reader, writer)

	if err := s.Run(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(writer.written) != 2 {
		t.Errorf("expected 2 secrets written, got %d", len(writer.written))
	}
}

func TestRun_ReaderError(t *testing.T) {
	reader := &mockReader{err: errors.New("vault unavailable")}
	writer := &mockWriter{}
	s := sync.NewWithDeps(baseConfig(), reader, writer)

	err := s.Run()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if writer.written != nil {
		t.Error("expected writer not to be called on reader error")
	}
}

func TestRun_WriterError(t *testing.T) {
	reader := &mockReader{secrets: map[string]string{"KEY": "value"}}
	writer := &mockWriter{err: errors.New("disk full")}
	s := sync.NewWithDeps(baseConfig(), reader, writer)

	if err := s.Run(); err == nil {
		t.Fatal("expected error from writer, got nil")
	}
}

func TestRun_EmptySecrets(t *testing.T) {
	reader := &mockReader{secrets: map[string]string{}}
	writer := &mockWriter{}
	s := sync.NewWithDeps(baseConfig(), reader, writer)

	if err := s.Run(); err != nil {
		t.Fatalf("expected no error for empty secrets, got %v", err)
	}
	if writer.written != nil {
		t.Error("expected writer not to be called when no secrets returned")
	}
}
