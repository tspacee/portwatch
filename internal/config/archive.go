package config

import "errors"

// ArchiveConfig controls the in-memory port snapshot archive.
type ArchiveConfig struct {
	Enabled bool `yaml:"enabled"`
	MaxSize int  `yaml:"max_size"`
}

// defaultArchiveConfig returns a sensible default archive configuration.
func defaultArchiveConfig() ArchiveConfig {
	return ArchiveConfig{
		Enabled: true,
		MaxSize: 500,
	}
}

// Validate checks that the ArchiveConfig fields are self-consistent.
func (a ArchiveConfig) Validate() error {
	if !a.Enabled {
		return nil
	}
	if a.MaxSize < 1 {
		return errors.New("archive: max_size must be at least 1 when enabled")
	}
	return nil
}
