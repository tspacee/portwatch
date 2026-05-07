package config

import (
	"errors"
	"time"
)

// CoolMapConfig holds configuration for the per-port CoolMap.
type CoolMapConfig struct {
	Enabled bool          `yaml:"enabled"`
	Window  time.Duration `yaml:"window"`
}

// defaultCoolMapConfig returns a CoolMapConfig with sensible defaults.
func defaultCoolMapConfig() CoolMapConfig {
	return CoolMapConfig{
		Enabled: true,
		Window:  10 * time.Second,
	}
}

// Validate checks that the CoolMapConfig is consistent.
func (c CoolMapConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return errors.New("coolmap: window must be positive when enabled")
	}
	return nil
}
