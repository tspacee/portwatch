package config

import (
	"errors"
	"time"
)

// BurstConfig holds configuration for the burst detector.
type BurstConfig struct {
	Enabled   bool          `yaml:"enabled"`
	Window    time.Duration `yaml:"window"`
	Threshold int           `yaml:"threshold"`
}

// defaultBurstConfig returns a BurstConfig with sensible defaults.
func defaultBurstConfig() BurstConfig {
	return BurstConfig{
		Enabled:   false,
		Window:    10 * time.Second,
		Threshold: 10,
	}
}

// Validate checks the BurstConfig for correctness.
func (c BurstConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window <= 0 {
		return errors.New("burst: window must be greater than zero")
	}
	if c.Threshold < 1 {
		return errors.New("burst: threshold must be at least 1")
	}
	return nil
}
