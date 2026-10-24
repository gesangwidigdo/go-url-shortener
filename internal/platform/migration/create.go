package migration

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var nonAlnum = regexp.MustCompile(`^[a-z0-9]+`)

func Create(dir, name string) ([]string, error) {
	// name = strings.Trim(nonAlnum.ReplaceAllString(strings.ToLower(name), "_"), "_")
	if name == "" {
		return nil, errors.New("migration name is required")
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("read migrations dir: %w", err)
	}

	highest := 0
	for _, e := range entries {
		prefix, _, ok := strings.Cut(e.Name(), "_")
		if !ok {
			continue
		}

		if n, err := strconv.Atoi(prefix); err == nil && n > highest {
			highest = n
		}
	}

	var created []string
	for _, direction := range []string{"up", "down"} {
		path := filepath.Join(dir, fmt.Sprintf("%06d_%s.%s.sql", highest+1, name, direction))
		if err := os.WriteFile(path, nil, 0o644); err != nil {
			return nil, fmt.Errorf("create %s: %w", path, err)
		}
		created = append(created, path)
	}

	return created, nil
}
