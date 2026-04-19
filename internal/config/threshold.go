package config

import "errors"

// ThresholdConfig holds configuration for the per-port threshold tracker.
type ThresholdConfig struct {
	Enabled bool `yaml:"enabled"`
	Limit   int  `yaml:"limit"`
}

// defaultThresholdConfig returns a ThresholdConfig with safe defaults.
func defaultThresholdConfig() ThresholdConfig {
	return ThresholdConfig{
		Enabled: false,
		Limit:   5,
	}
}

// Validate checks that the ThresholdConfig is logically consistent.
func (c ThresholdConfig) Validate() error {
	if c.Enabled && c.Limit < 1 {
		return errors.New("threshold: limit must be >= 1 when enabled")
	}
	return nil
}
