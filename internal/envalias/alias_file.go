package envalias

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// LoadFile reads alias pairs from a file, one "from=to" per line.
// Lines starting with '#' and blank lines are ignored.
func LoadFile(path string) (*Mapper, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envalias: open %s: %w", path, err)
	}
	defer f.Close()

	var pairs []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		pairs = append(pairs, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envalias: read %s: %w", path, err)
	}
	return New(pairs), nil
}
