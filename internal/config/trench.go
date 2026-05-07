package config

import (
	"errors"
	"time"
)

// TrenchConfig controls the Trench absence-tracking component.
type TrenchConfig struct {
	Enabled bool          `yaml:"enabled"`
	Window  time.Duration `yaml:"window"`
}

// defaultTrenchConfig returns a sensible default: disabled, 30-minute window.
func defaultTrenchConfig() TrenchConfig {
	return TrenchConfig{
		Enabled: false,
		Window:  30 * time.Minute,
	}
}

// Validate checks TrenchConfig for correctness.
func (c TrenchConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return errors.New("trench: window must be positive when enabled")
	}
	return nil
}
