package envencrypt_test

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"

	"github.com/your-org/vaultpull/internal/envencrypt"
)

func writeTempKey(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "key.txt")
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("write key file: %v", err)
	}
	return p
}

func TestLoadKeyFile_RawBytes(t *testing.T) {
	key := "0123456789abcdef" // 16 raw bytes
	p := writeTempKey(t, key)
	got, err := envencrypt.LoadKeyFile(p)
	if err != nil {
		t.Fatalf("LoadKeyFile: %v", err)
	}
	if string(got) != key {
		t.Errorf("got %q, want %q", got, key)
	}
}

func TestLoadKeyFile_HexEncoded(t *testing.T) {
	raw := make([]byte, 32)
	for i := range raw {
		raw[i] = byte(i)
	}
	hexStr := hex.EncodeToString(raw)
	p := writeTempKey(t, hexStr+"\n")
	got, err := envencrypt.LoadKeyFile(p)
	if err != nil {
		t.Fatalf("LoadKeyFile: %v", err)
	}
	if len(got) != 32 {
		t.Errorf("got %d bytes, want 32", len(got))
	}
}

func TestLoadKeyFile_Missing(t *testing.T) {
	_, err := envencrypt.LoadKeyFile("/nonexistent/key.txt")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoadKeyFile_Empty(t *testing.T) {
	p := writeTempKey(t, "   \n")
	_, err := envencrypt.LoadKeyFile(p)
	if err == nil {
		t.Fatal("expected error for empty key file")
	}
}

func TestKeyFromEnv_Present(t *testing.T) {
	t.Setenv("VAULT_ENC_KEY", "0123456789abcdef")
	got, err := envencrypt.KeyFromEnv("VAULT_ENC_KEY")
	if err != nil {
		t.Fatalf("KeyFromEnv: %v", err)
	}
	if string(got) != "0123456789abcdef" {
		t.Errorf("unexpected key: %q", got)
	}
}

func TestKeyFromEnv_Missing(t *testing.T) {
	os.Unsetenv("VAULT_ENC_KEY_MISSING")
	_, err := envencrypt.KeyFromEnv("VAULT_ENC_KEY_MISSING")
	if err == nil {
		t.Fatal("expected error for unset env var")
	}
}
