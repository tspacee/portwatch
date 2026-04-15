package config

import (
	"fmt"
	"os"
)

// LoadOrDefault attempts to load config from path.
// If path is empty or the file does not exist, Default() is returned.
func LoadOrDefault(path string) (*Config, error) {
	if path == "" {
		return Default(), nil
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "portwatch: config file %q not found, using defaults\n", path)
		return Default(), nil
	}

	return Load(path)
}
