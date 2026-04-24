package envseal

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
)

// LoadKeyFile reads an HMAC key from path.
// The file may contain raw bytes or a hex-encoded string (64 hex chars = 32 bytes).
func LoadKeyFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envseal: load key file: %w", err)
	}
	raw := strings.TrimSpace(string(data))
	if raw == "" {
		return nil, errors.New("envseal: key file is empty")
	}
	// Attempt hex decode if it looks like a hex string.
	if isHex(raw) {
		decoded, err := hex.DecodeString(raw)
		if err == nil && len(decoded) > 0 {
			return decoded, nil
		}
	}
	return []byte(raw), nil
}

// KeyFromEnv reads the HMAC key from the named environment variable.
func KeyFromEnv(envVar string) ([]byte, error) {
	val := os.Getenv(envVar)
	if val == "" {
		return nil, fmt.Errorf("envseal: env var %q is not set or empty", envVar)
	}
	if isHex(val) {
		decoded, err := hex.DecodeString(val)
		if err == nil && len(decoded) > 0 {
			return decoded, nil
		}
	}
	return []byte(val), nil
}

func isHex(s string) bool {
	if len(s)%2 != 0 || len(s) < 2 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
