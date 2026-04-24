// Package envhash provides content-based hashing for secret maps,
// allowing callers to detect whether a set of secrets has changed
// since the last sync without comparing values directly.
package envhash

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

// Hasher computes a deterministic hash over a map of secret key/value pairs.
type Hasher struct {
	separator string
}

// New returns a Hasher with the default key=value separator.
func New() *Hasher {
	return &Hasher{separator: "="}
}

// NewWithSeparator returns a Hasher that uses sep between key and value
// when building the hash input. sep must not be empty.
func NewWithSeparator(sep string) (*Hasher, error) {
	if sep == "" {
		return nil, fmt.Errorf("envhash: separator must not be empty")
	}
	return &Hasher{separator: sep}, nil
}

// Sum returns a lowercase hex-encoded SHA-256 digest of the supplied secrets.
// The digest is order-independent: keys are sorted before hashing so that
// two maps with identical contents always produce the same hash.
func (h *Hasher) Sum(secrets map[string]string) string {
	if len(secrets) == 0 {
		return emptySHA256()
	}

	keys := make([]string, 0, len(secrets))
	for k := range secrets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString(h.separator)
		sb.WriteString(secrets[k])
		sb.WriteByte('\n')
	}

	sum := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(sum[:])
}

// Equal reports whether two secret maps produce the same hash.
func (h *Hasher) Equal(a, b map[string]string) bool {
	return h.Sum(a) == h.Sum(b)
}

func emptySHA256() string {
	sum := sha256.Sum256(nil)
	return hex.EncodeToString(sum[:])
}
