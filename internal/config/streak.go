package config

import "errors"

// StreakConfig holds configuration for the port streak tracker.
type StreakConfig struct {
	// Enabled controls whether streak tracking is active.
	Enabled bool `yaml:"enabled"`

	// AlertThreshold is the consecutive scan count at which an alert
	// should be raised for a port. Zero means no threshold alerting.
	AlertThreshold int `yaml:"alert_threshold"`
}

// defaultStreakConfig returns a StreakConfig with sensible defaults.
func defaultStreakConfig() StreakConfig {
	return StreakConfig{
		Enabled:        true,
		AlertThreshold: 5,
	}
}

// Validate checks that the StreakConfig fields are consistent.
func (c StreakConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.AlertThreshold < 0 {
		return errors.New("streak: alert_threshold must be >= 0")
	}
	return nil
}
