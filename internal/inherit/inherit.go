// Package inherit provides secret inheritance across Vault paths,
// allowing child paths to merge secrets from parent paths with
// child values taking precedence.
package inherit

import (
	"strings"
)

// Resolver merges secrets from ancestor paths into a final map.
type Resolver struct {
	reader func(path string) (map[string]string, error)
	sep    string
}

// New returns a Resolver that uses reader to fetch secrets at each path.
func New(reader func(path string) (map[string]string, error)) *Resolver {
	return &Resolver{reader: reader, sep: "/"}
}

// Resolve walks from the root down to path, merging secrets at each
// level. Values from deeper (more specific) paths override ancestors.
func (r *Resolver) Resolve(path string) (map[string]string, error) {
	ancestors := ancestors(strings.Trim(path, r.sep), r.sep)
	result := make(map[string]string)
	for _, p := range ancestors {
		secrets, err := r.reader(p)
		if err != nil {
			return nil, err
		}
		for k, v := range secrets {
			result[k] = v
		}
	}
	return result, nil
}

// ancestors returns all prefix paths including path itself, from
// shallowest to deepest. E.g. "a/b/c" → ["a", "a/b", "a/b/c"].
func ancestors(path, sep string) []string {
	if path == "" {
		return nil
	}
	parts := strings.Split(path, sep)
	out := make([]string, 0, len(parts))
	for i := range parts {
		out = append(out, strings.Join(parts[:i+1], sep))
	}
	return out
}
