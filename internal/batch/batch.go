// Package batch provides concurrent secret fetching across multiple Vault paths.
package batch

import (
	"context"
	"fmt"
	"sync"
)

// Reader reads secrets from a single path.
type Reader interface {
	ReadSecrets(ctx context.Context, path string) (map[string]string, error)
}

// Result holds the outcome of reading one path.
type Result struct {
	Path    string
	Secrets map[string]string
	Err     error
}

// Fetcher fetches secrets from multiple paths concurrently.
type Fetcher struct {
	reader      Reader
	concurrency int
}

// New returns a Fetcher with the given reader and concurrency limit.
func New(r Reader, concurrency int) *Fetcher {
	if concurrency <= 0 {
		concurrency = 4
	}
	return &Fetcher{reader: r, concurrency: concurrency}
}

// FetchAll reads all provided paths concurrently and returns results in
// the same order as paths. Cancelled contexts abort in-flight requests.
func (f *Fetcher) FetchAll(ctx context.Context, paths []string) []Result {
	results := make([]Result, len(paths))
	sem := make(chan struct{}, f.concurrency)
	var wg sync.WaitGroup

	for i, p := range paths {
		wg.Add(1)
		go func(idx int, path string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			if ctx.Err() != nil {
				results[idx] = Result{Path: path, Err: fmt.Errorf("context cancelled")}
				return
			}
			secrets, err := f.reader.ReadSecrets(ctx, path)
			results[idx] = Result{Path: path, Secrets: secrets, Err: err}
		}(i, p)
	}

	wg.Wait()
	return results
}

// Merge combines all successful results into a single map.
// Later paths overwrite earlier ones on key collision.
func Merge(results []Result) (map[string]string, error) {
	out := make(map[string]string)
	for _, r := range results {
		if r.Err != nil {
			return nil, fmt.Errorf("path %q: %w", r.Path, r.Err)
		}
		for k, v := range r.Secrets {
			out[k] = v
		}
	}
	return out, nil
}
