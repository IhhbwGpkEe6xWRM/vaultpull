// Package pagesize provides pagination support for Vault list operations,
// allowing large secret lists to be fetched in bounded chunks.
package pagesize

import "errors"

const (
	DefaultPageSize = 100
	MinPageSize     = 1
	MaxPageSize     = 1000
)

// Config holds pagination settings.
type Config struct {
	PageSize int
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{PageSize: DefaultPageSize}
}

// Paginator splits a slice of keys into pages.
type Paginator struct {
	pageSize int
}

// New returns a Paginator with the given config.
func New(cfg Config) (*Paginator, error) {
	if cfg.PageSize < MinPageSize || cfg.PageSize > MaxPageSize {
		return nil, errors.New("pagesize: page size out of range [1, 1000]")
	}
	return &Paginator{pageSize: cfg.PageSize}, nil
}

// Pages partitions keys into slices of at most PageSize length.
func (p *Paginator) Pages(keys []string) [][]string {
	if len(keys) == 0 {
		return nil
	}
	var pages [][]string
	for i := 0; i < len(keys); i += p.pageSize {
		end := i + p.pageSize
		if end > len(keys) {
			end = len(keys)
		}
		page := make([]string, end-i)
		copy(page, keys[i:end])
		pages = append(pages, page)
	}
	return pages
}

// TotalPages returns the number of pages for a given key count.
func (p *Paginator) TotalPages(total int) int {
	if total <= 0 {
		return 0
	}
	return (total + p.pageSize - 1) / p.pageSize
}
