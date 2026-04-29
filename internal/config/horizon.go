package config

import (
	"errors"
	"time"
)

// HorizonConfig configures the Horizon tracker.
type HorizonConfig struct {
	Enabled bool          `yaml:"enabled"`
	Cutoff  time.Duration `yaml:"cutoff"`
}

// defaultHorizonConfig returns a HorizonConfig with sensible defaults.
func defaultHorizonConfig() HorizonConfig {
	return HorizonConfig{
		Enabled: false,
		Cutoff:  24 * time.Hour,
	}
}

// Validate checks that the HorizonConfig fields are valid.
func (h HorizonConfig) Validate() error {
	if !h.Enabled {
		return nil
	}
	if h.Cutoff <= 0 {
		return errors.New("horizon: cutoff must be a positive duration")
	}
	return nil
}
