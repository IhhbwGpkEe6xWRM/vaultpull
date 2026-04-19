package envencrypt

import (
	"encoding/hex"
	"errors"
	"os"
	"strings"
)

// ErrEmptyKeyFile is returned when the key file is empty.
var ErrEmptyKeyFile = errors.New("envencrypt: key file is empty")

// LoadKeyFile reads an AES key from a file.
// The file may contain raw bytes or a hex-encoded string.
// Whitespace is trimmed before decoding.
func LoadKeyFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := strings.TrimSpace(string(data))
	if s == "" {
		return nil, ErrEmptyKeyFile
	}
	// Attempt hex decode first.
	if decoded, err := hex.DecodeString(s); err == nil {
		return decoded, nil
	}
	// Fall back to raw bytes.
	return []byte(s), nil
}

// KeyFromEnv reads an AES key from an environment variable.
// The value may be raw bytes or hex-encoded.
func KeyFromEnv(varName string) ([]byte, error) {
	val := os.Getenv(varName)
	if val == "" {
		return nil, errors.New("envencrypt: env var " + varName + " is not set")
	}
	if decoded, err := hex.DecodeString(val); err == nil {
		return decoded, nil
	}
	return []byte(val), nil
}
