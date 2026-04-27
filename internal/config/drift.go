package config

import "errors"

// DriftConfig controls the port-drift tracking feature.
type DriftConfig struct {
	Enabled bool    `yaml:"enabled"`
	Decay   float64 `yaml:"decay"`
}

// defaultDriftConfig returns a safe default configuration for drift tracking.
func defaultDriftConfig() DriftConfig {
	return DriftConfig{
		Enabled: false,
		Decay:   0.3,
	}
}

// Validate checks that the DriftConfig fields are within acceptable bounds.
func (c DriftConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Decay <= 0 || c.Decay > 1 {
		return errors.New("drift: decay must be in range (0, 1]")
	}
	return nil
}
