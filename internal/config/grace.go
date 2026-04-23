package config

import (
	"errors"
	"time"
)

// GraceConfig controls the per-port grace period feature.
type GraceConfig struct {
	Enabled bool          `yaml:"enabled"`
	Window  time.Duration `yaml:"window"`
}

// defaultGraceConfig returns a GraceConfig with sensible defaults.
// Grace is enabled with a 5-second window by default.
func defaultGraceConfig() GraceConfig {
	return GraceConfig{
		Enabled: true,
		Window:  5 * time.Second,
	}
}

// Validate checks that the GraceConfig is consistent.
func (c GraceConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return errors.New("grace: window must be positive when enabled")
	}
	return nil
}
