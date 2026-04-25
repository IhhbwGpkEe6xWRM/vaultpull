// Package envobfuscate provides key-name obfuscation for env maps.
// Keys are replaced with deterministic aliases derived from a shared
// salt so that the mapping can be reversed by a caller with the same salt.
package envobfuscate

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
)

// Obfuscator replaces env key names with deterministic hex aliases.
type Obfuscator struct {
	salt []byte
}

// New returns an Obfuscator seeded with the given salt.
// An empty salt is accepted but produces weaker aliases.
func New(salt []byte) *Obfuscator {
	s := make([]byte, len(salt))
	copy(s, salt)
	return &Obfuscator{salt: s}
}

// alias returns a short deterministic hex string for key.
func (o *Obfuscator) alias(key string) string {
	mac := hmac.New(sha256.New, o.salt)
	mac.Write([]byte(key))
	return hex.EncodeToString(mac.Sum(nil))[:16]
}

// Apply returns a new map whose keys are replaced by their aliases.
// The original values are preserved unchanged.
func (o *Obfuscator) Apply(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k, v := range src {
		out[o.alias(k)] = v
	}
	return out
}

// Reveal returns the inverse mapping: alias -> original key, for every
// key present in src. It does not recover values.
func (o *Obfuscator) Reveal(src map[string]string) map[string]string {
	out := make(map[string]string, len(src))
	for k := range src {
		out[o.alias(k)] = k
	}
	return out
}

// Keys returns the sorted list of aliases that would be produced for keys.
func (o *Obfuscator) Keys(keys []string) []string {
	aliases := make([]string, len(keys))
	for i, k := range keys {
		aliases[i] = fmt.Sprintf("%s", o.alias(k))
	}
	sort.Strings(aliases)
	return aliases
}
