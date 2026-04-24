package config

import (
	"errors"
	"time"
)

// SuppressConfig controls the port suppression (mute) feature.
type SuppressConfig struct {
	Enabled bool          `yaml:"enabled"`
	Window  time.Duration `yaml:"window"`
}

// defaultSuppressConfig returns a SuppressConfig with sensible defaults.
func defaultSuppressConfig() SuppressConfig {
	return SuppressConfig{
		Enabled: false,
		Window:  5 * time.Minute,
	}
}

// Validate checks that the SuppressConfig is consistent.
func (c SuppressConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return errors.New("suppress: window must be positive when enabled")
	}
	return nil
}
