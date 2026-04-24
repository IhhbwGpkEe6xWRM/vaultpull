// Package envseal provides tamper-detection for .env files using HMAC signatures.
// A seal is computed over the sorted key=value pairs and stored alongside the file.
// Verification fails if any key, value, or ordering has changed since sealing.
package envseal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

// ErrTampered is returned when the seal does not match the current file contents.
var ErrTampered = errors.New("envseal: seal mismatch — file may have been tampered with")

// ErrMissingSeal is returned when no seal file exists for the given env file.
var ErrMissingSeal = errors.New("envseal: no seal found for file")

// Sealer signs and verifies env secret maps.
type Sealer struct {
	key []byte
}

// New creates a Sealer using the provided HMAC key.
func New(key []byte) (*Sealer, error) {
	if len(key) == 0 {
		return nil, errors.New("envseal: key must not be empty")
	}
	return &Sealer{key: key}, nil
}

// Sign computes an HMAC-SHA256 signature over the sorted key=value pairs.
func (s *Sealer) Sign(secrets map[string]string) string {
	h := hmac.New(sha256.New, s.key)
	for _, line := range canonical(secrets) {
		h.Write([]byte(line))
	}
	return hex.EncodeToString(h.Sum(nil))
}

// Verify returns nil if the signature matches the current secrets, ErrTampered otherwise.
func (s *Sealer) Verify(secrets map[string]string, sig string) error {
	expected := s.Sign(secrets)
	if !hmac.Equal([]byte(expected), []byte(sig)) {
		return ErrTampered
	}
	return nil
}

// WriteSealFile writes the signature for secrets to sealPath.
func (s *Sealer) WriteSealFile(secrets map[string]string, sealPath string) error {
	sig := s.Sign(secrets)
	return os.WriteFile(sealPath, []byte(sig), 0600)
}

// VerifySealFile reads the signature from sealPath and verifies it against secrets.
func (s *Sealer) VerifySealFile(secrets map[string]string, sealPath string) error {
	data, err := os.ReadFile(sealPath)
	if errors.Is(err, os.ErrNotExist) {
		return ErrMissingSeal
	}
	if err != nil {
		return fmt.Errorf("envseal: read seal file: %w", err)
	}
	return s.Verify(secrets, strings.TrimSpace(string(data)))
}

// canonical returns sorted "key=value" lines for deterministic signing.
func canonical(secrets map[string]string) []string {
	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	lines := make([]string, len(keys))
	for i, k := range keys {
		lines[i] = k + "=" + secrets[k]
	}
	return lines
}
