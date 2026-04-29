package config

import (
	"errors"
	"time"
)

// VelocityConfig holds configuration for the port-change velocity tracker.
type VelocityConfig struct {
	Enabled bool          `yaml:"enabled"`
	Window  time.Duration `yaml:"window"`
}

// defaultVelocityConfig returns a VelocityConfig with sensible defaults.
func defaultVelocityConfig() VelocityConfig {
	return VelocityConfig{
		Enabled: true,
		Window:  30 * time.Second,
	}
}

// Validate checks that the VelocityConfig is well-formed.
func (c VelocityConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return errors.New("velocity: window must be positive")
	}
	return nil
}
