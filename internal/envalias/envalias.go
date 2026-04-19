// Package envalias maps vault secret keys to alternative environment variable names.
package envalias

import "fmt"

// Alias represents a single key mapping: From is the vault key, To is the env var name.
type Alias struct {
	From string
	To   string
}

// Mapper applies a set of aliases to a secrets map.
type Mapper struct {
	aliases []Alias
}

// New creates a Mapper from a slice of raw "from=to" pair strings.
// Malformed pairs are silently ignored.
func New(pairs []string) *Mapper {
	var aliases []Alias
	for _, p := range pairs {
		a, err := parsePair(p)
		if err == nil {
			aliases = append(aliases, a)
		}
	}
	return &Mapper{aliases: aliases}
}

// Apply returns a copy of secrets with aliased keys renamed.
// If a From key is not present the alias is skipped.
// If To already exists it is not overwritten.
func (m *Mapper) Apply(secrets map[string]string) map[string]string {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		out[k] = v
	}
	for _, a := range m.aliases {
		v, ok := out[a.From]
		if !ok {
			continue
		}
		if _, exists := out[a.To]; !exists {
			out[a.To] = v
		}
		delete(out, a.From)
	}
	return out
}

// Pairs returns the parsed aliases.
func (m *Mapper) Pairs() []Alias { return m.aliases }

func parsePair(s string) (Alias, error) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			from, to := s[:i], s[i+1:]
			if from == "" || to == "" {
				return Alias{}, fmt.Errorf("envalias: empty from or to in %q", s)
			}
			return Alias{From: from, To: to}, nil
		}
	}
	return Alias{}, fmt.Errorf("envalias: no '=' in pair %q", s)
}
