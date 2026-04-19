// Package envdiff compares a set of resolved secrets against an existing
// .env file and reports which keys would be added, removed, or changed.
package envdiff

import (
	"bufio"
	"os"
	"strings"
)

// Result holds the categorised keys from a comparison.
type Result struct {
	Added    []string
	Removed  []string
	Modified []string
	Unchanged []string
}

// HasChanges returns true when at least one key differs.
func (r Result) HasChanges() bool {
	return len(r.Added)+len(r.Removed)+len(r.Modified) > 0
}

// Compare reads the existing env file at path and compares it against
// incoming. If the file does not exist it treats all keys as added.
func Compare(path string, incoming map[string]string) (Result, error) {
	existing, err := parseEnvFile(path)
	if err != nil && !os.IsNotExist(err) {
		return Result{}, err
	}

	var res Result
	for k, v := range incoming {
		if old, ok := existing[k]; !ok {
			res.Added = append(res.Added, k)
		} else if old != v {
			res.Modified = append(res.Modified, k)
		} else {
			res.Unchanged = append(res.Unchanged, k)
		}
	}
	for k := range existing {
		if _, ok := incoming[k]; !ok {
			res.Removed = append(res.Removed, k)
		}
	}
	sortStrings(res.Added)
	sortStrings(res.Removed)
	sortStrings(res.Modified)
	sortStrings(res.Unchanged)
	return res, nil
}

func parseEnvFile(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	out := make(map[string]string)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		out[strings.TrimSpace(parts[0])] = strings.Trim(strings.TrimSpace(parts[1]), `"`)
	}
	return out, scanner.Err()
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}
