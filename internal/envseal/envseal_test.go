package envseal_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envseal"
)

func newSealer(t *testing.T) *envseal.Sealer {
	t.Helper()
	s, err := envseal.New([]byte("test-hmac-key-32bytes-padded!!"))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func TestNew_EmptyKey_ReturnsError(t *testing.T) {
	_, err := envseal.New(nil)
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestSign_DeterministicOutput(t *testing.T) {
	s := newSealer(t)
	secrets := map[string]string{"DB_HOST": "localhost", "DB_PORT": "5432"}
	if s.Sign(secrets) != s.Sign(secrets) {
		t.Fatal("Sign must be deterministic")
	}
}

func TestSign_OrderIndependent(t *testing.T) {
	s := newSealer(t)
	a := map[string]string{"A": "1", "B": "2"}
	b := map[string]string{"B": "2", "A": "1"}
	if s.Sign(a) != s.Sign(b) {
		t.Fatal("Sign must be order-independent")
	}
}

func TestVerify_ValidSignature(t *testing.T) {
	s := newSealer(t)
	secrets := map[string]string{"TOKEN": "abc123"}
	sig := s.Sign(secrets)
	if err := s.Verify(secrets, sig); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestVerify_TamperedValue(t *testing.T) {
	s := newSealer(t)
	secrets := map[string]string{"TOKEN": "abc123"}
	sig := s.Sign(secrets)
	secrets["TOKEN"] = "hacked"
	if err := s.Verify(secrets, sig); err != envseal.ErrTampered {
		t.Fatalf("expected ErrTampered, got %v", err)
	}
}

func TestVerify_AddedKey(t *testing.T) {
	s := newSealer(t)
	secrets := map[string]string{"A": "1"}
	sig := s.Sign(secrets)
	secrets["EXTRA"] = "injected"
	if err := s.Verify(secrets, sig); err != envseal.ErrTampered {
		t.Fatalf("expected ErrTampered, got %v", err)
	}
}

func TestWriteAndVerifySealFile_RoundTrip(t *testing.T) {
	s := newSealer(t)
	secrets := map[string]string{"KEY": "value"}
	path := filepath.Join(t.TempDir(), ".env.seal")
	if err := s.WriteSealFile(secrets, path); err != nil {
		t.Fatalf("WriteSealFile: %v", err)
	}
	if err := s.VerifySealFile(secrets, path); err != nil {
		t.Fatalf("VerifySealFile: %v", err)
	}
}

func TestVerifySealFile_MissingFile(t *testing.T) {
	s := newSealer(t)
	secrets := map[string]string{"K": "v"}
	err := s.VerifySealFile(secrets, "/nonexistent/.env.seal")
	if err != envseal.ErrMissingSeal {
		t.Fatalf("expected ErrMissingSeal, got %v", err)
	}
}

func TestVerifySealFile_TamperedAfterWrite(t *testing.T) {
	s := newSealer(t)
	secrets := map[string]string{"SECRET": "original"}
	path := filepath.Join(t.TempDir(), ".env.seal")
	_ = s.WriteSealFile(secrets, path)
	tampered := map[string]string{"SECRET": "changed"}
	if err := s.VerifySealFile(tampered, path); err != envseal.ErrTampered {
		t.Fatalf("expected ErrTampered, got %v", err)
	}
}

func TestWriteSealFile_PermissionsAreRestrictive(t *testing.T) {
	s := newSealer(t)
	path := filepath.Join(t.TempDir(), ".env.seal")
	_ = s.WriteSealFile(map[string]string{"X": "y"}, path)
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Fatalf("expected 0600, got %v", info.Mode().Perm())
	}
}
