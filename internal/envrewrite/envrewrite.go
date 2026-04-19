// Package envrewrite provides key rewriting rules for secret maps,
// allowing keys to be renamed, aliased, or dropped before writing to .env files.
package envrewrite

import "fmt"

// Rule defines a single rewrite operation.
type Rule struct {
	From string
	To   string // empty means drop the key
}

// Rewriter applies a set of rewrite rules to a secret map.
type Rewriter struct {
	rules []Rule
}

// New returns a Rewriter with the given rules.
func New(rules []Rule) (*Rewriter, error) {
	for i, r := range rules {
		if r.From == "" {
			return nil, fmt.Errorf("rule %d: From must not be empty", i)
		}
	}
	return &Rewriter{rules: rules}, nil
}

// Apply returns a new map with rewrite rules applied.
// Keys matched by a rule are renamed to the To value.
// If To is empty the key is dropped.
// Unmatched keys are passed through unchanged.
func (r *Rewriter) Apply(secrets map[string]string) map[string]string {
	index := make(map[string]Rule, len(r.rules))
	for _, rule := range r.rules {
		index[rule.From] = rule
	}

	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		rule, matched := index[k]
		switch {
		case !matched:
			out[k] = v
		case rule.To == "":
			// drop
		default:
			out[rule.To] = v
		}
	}
	return out
}

// Keys returns the list of From keys covered by the rewriter's rules.
func (r *Rewriter) Keys() []string {
	keys := make([]string, 0, len(r.rules))
	for _, rule := range r.rules {
		keys = append(keys, rule.From)
	}
	return keys
}
