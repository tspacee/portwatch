package config

import "fmt"

// ClampConfig holds configuration for the port count clamp feature.
type ClampConfig struct {
	Enabled bool `yaml:"enabled"`
	Min     int  `yaml:"min"`
	Max     int  `yaml:"max"`
}

// defaultClampConfig returns a safe default ClampConfig.
func defaultClampConfig() ClampConfig {
	return ClampConfig{
		Enabled: false,
		Min:     0,
		Max:     1024,
	}
}

// Validate checks that the ClampConfig fields are internally consistent.
func (c ClampConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Min < 0 {
		return fmt.Errorf("clamp: min must be non-negative, got %d", c.Min)
	}
	if c.Max < c.Min {
		return fmt.Errorf("clamp: max (%d) must be >= min (%d)", c.Max, c.Min)
	}
	return nil
}
