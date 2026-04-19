// Package flatten provides utilities for flattening nested Vault secret
// maps into a single-level key=value map suitable for .env files.
package flatten

import "fmt"

// Separator is the default delimiter used when joining nested keys.
const Separator = "_"

// Flattener converts nested map[string]any structures into a flat
// map[string]string, joining key segments with a configurable separator.
type Flattener struct {
	sep string
}

// New returns a Flattener using the default separator.
func New() *Flattener {
	return &Flattener{sep: Separator}
}

// NewWithSeparator returns a Flattener that joins keys with sep.
func NewWithSeparator(sep string) *Flattener {
	if sep == "" {
		sep = Separator
	}
	return &Flattener{sep: sep}
}

// Flatten converts a potentially nested map into a flat map[string]string.
// Nested maps are recursively expanded; all other values are formatted with
// fmt.Sprintf("%v", value).
func (f *Flattener) Flatten(input map[string]any) map[string]string {
	out := make(map[string]string)
	f.flatten(input, "", out)
	return out
}

func (f *Flattener) flatten(m map[string]any, prefix string, out map[string]string) {
	for k, v := range m {
		key := k
		if prefix != "" {
			key = prefix + f.sep + k
		}
		switch val := v.(type) {
		case map[string]any:
			f.flatten(val, key, out)
		case map[string]string:
			for sk, sv := range val {
				out[key+f.sep+sk] = sv
			}
		default:
			out[key] = fmt.Sprintf("%v", val)
		}
	}
}
