package config

import (
	"errors"
	"time"
)

// FlapConfig controls the flap-detection subsystem.
type FlapConfig struct {
	Enabled   bool          `yaml:"enabled"`
	Window    time.Duration `yaml:"window"`
	Threshold int           `yaml:"threshold"`
}

// defaultFlapConfig returns sensible defaults: enabled, 30-second window,
// threshold of 4 state changes.
func defaultFlapConfig() FlapConfig {
	return FlapConfig{
		Enabled:   true,
		Window:    30 * time.Second,
		Threshold: 4,
	}
}

// Validate returns an error if the configuration is invalid.
func (c FlapConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return errors.New("flap: window must be positive")
	}
	if c.Threshold < 2 {
		return errors.New("flap: threshold must be at least 2")
	}
	return nil
}
