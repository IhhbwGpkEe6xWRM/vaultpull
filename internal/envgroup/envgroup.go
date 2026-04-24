// Package envgroup groups environment variables by a shared key prefix,
// returning named sub-maps that can be processed independently.
package envgroup

import (
	"fmt"
	"strings"
)

// Group holds a named collection of key-value pairs extracted from a larger
// secrets map. Keys inside the group have their shared prefix stripped.
type Group struct {
	Name   string
	Values map[string]string
}

// Grouper partitions a flat secrets map into named groups based on prefix
// rules registered at construction time.
type Grouper struct {
	rules []rule
}

type rule struct {
	name   string
	prefix string // normalised: no leading/trailing underscore
}

// New creates a Grouper from a slice of "name=PREFIX" pair strings.
// Malformed pairs or pairs with empty names or prefixes are rejected.
func New(pairs []string) (*Grouper, error) {
	g := &Grouper{}
	for _, p := range pairs {
		idx := strings.IndexByte(p, '=')
		if idx < 1 {
			return nil, fmt.Errorf("envgroup: malformed pair %q: want name=PREFIX", p)
		}
		name := p[:idx]
		prefix := strings.Trim(p[idx+1:], "_")
		if prefix == "" {
			return nil, fmt.Errorf("envgroup: empty prefix for group %q", name)
		}
		g.rules = append(g.rules, rule{name: name, prefix: strings.ToUpper(prefix)})
	}
	return g, nil
}

// Split partitions secrets into groups. Keys that do not match any rule are
// collected in a group named "" (the default group).
func (g *Grouper) Split(secrets map[string]string) []Group {
	buckets := make(map[string]map[string]string)

	for k, v := range secrets {
		matched := false
		for _, r := range g.rules {
			upper := strings.ToUpper(k)
			if upper == r.prefix || strings.HasPrefix(upper, r.prefix+"_") {
				if buckets[r.name] == nil {
					buckets[r.name] = make(map[string]string)
				}
				stripped := k[len(r.prefix):]
				stripped = strings.TrimPrefix(stripped, "_")
				if stripped == "" {
					stripped = k
				}
				buckets[r.name][stripped] = v
				matched = true
				break
			}
		}
		if !matched {
			if buckets[""] == nil {
				buckets[""] = make(map[string]string)
			}
			buckets[""][k] = v
		}
	}

	result := make([]Group, 0, len(buckets))
	for name, vals := range buckets {
		result = append(result, Group{Name: name, Values: vals})
	}
	return result
}
