// Package envchain resolves secrets through an ordered chain of sources,
// returning the first successful result for each key.
package envchain

import "context"

// Reader is the interface for any secret source.
type Reader interface {
	Read(ctx context.Context, path string) (map[string]string, error)
}

// Chain holds an ordered list of readers. For each key, the first reader
// that returns a non-empty value wins.
type Chain struct {
	readers []Reader
}

// New creates a Chain from the provided readers. Readers are consulted in
// order; earlier readers take precedence.
func New(readers ...Reader) *Chain {
	return &Chain{readers: readers}
}

// Resolve queries every reader for the given path and merges the results.
// Later readers fill in keys that earlier readers did not provide; earlier
// readers always win on conflicts.
func (c *Chain) Resolve(ctx context.Context, path string) (map[string]string, error) {
	result := make(map[string]string)

	// Iterate in reverse so that earlier readers overwrite later ones.
	for i := len(c.readers) - 1; i >= 0; i-- {
		secrets, err := c.readers[i].Read(ctx, path)
		if err != nil {
			return nil, err
		}
		for k, v := range secrets {
			if v != "" {
				result[k] = v
			}
		}
	}

	return result, nil
}

// Len returns the number of readers in the chain.
func (c *Chain) Len() int {
	return len(c.readers)
}
