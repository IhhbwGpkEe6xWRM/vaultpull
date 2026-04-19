package envencrypt_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpull/internal/envencrypt"
)

func newEncrypter(t *testing.T) *envencrypt.Encrypter {
	t.Helper()
	key := []byte("0123456789abcdef") // 16 bytes
	e, err := envencrypt.New(key)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return e
}

func TestNew_ValidKey(t *testing.T) {
	for _, size := range []int{16, 24, 32} {
		key := make([]byte, size)
		_, err := envencrypt.New(key)
		if err != nil {
			t.Errorf("expected no error for %d-byte key, got %v", size, err)
		}
	}
}

func TestNew_InvalidKey(t *testing.T) {
	_, err := envencrypt.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for invalid key length")
	}
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	e := newEncrypter(t)
	plain := "super-secret-value"
	cipher, err := e.Encrypt(plain)
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if cipher == plain {
		t.Fatal("ciphertext should differ from plaintext")
	}
	got, err := e.Decrypt(cipher)
	if err != nil {
		t.Fatalf("Decrypt: %v", err)
	}
	if got != plain {
		t.Errorf("got %q, want %q", got, plain)
	}
}

func TestEncrypt_ProducesBase64(t *testing.T) {
	e := newEncrypter(t)
	cipher, err := e.Encrypt("value")
	if err != nil {
		t.Fatalf("Encrypt: %v", err)
	}
	if strings.ContainsAny(cipher, " \t\n") {
		t.Errorf("ciphertext contains whitespace: %q", cipher)
	}
}

func TestDecrypt_InvalidBase64(t *testing.T) {
	e := newEncrypter(t)
	_, err := e.Decrypt("!!!not-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestDecrypt_TamperedCiphertext(t *testing.T) {
	e := newEncrypter(t)
	cipher, _ := e.Encrypt("value")
	tampered := cipher[:len(cipher)-4] + "AAAA"
	_, err := e.Decrypt(tampered)
	if err == nil {
		t.Fatal("expected error for tampered ciphertext")
	}
}

func TestEncryptMap_DecryptMap_RoundTrip(t *testing.T) {
	e := newEncrypter(t)
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}
	enc, err := e.EncryptMap(secrets)
	if err != nil {
		t.Fatalf("EncryptMap: %v", err)
	}
	for k, v := range secrets {
		if enc[k] == v {
			t.Errorf("key %q: ciphertext equals plaintext", k)
		}
	}
	dec, err := e.DecryptMap(enc)
	if err != nil {
		t.Fatalf("DecryptMap: %v", err)
	}
	for k, want := range secrets {
		if dec[k] != want {
			t.Errorf("key %q: got %q, want %q", k, dec[k], want)
		}
	}
}
