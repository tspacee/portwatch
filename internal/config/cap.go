package config

import "errors"

// CapConfig controls the Cap primitive used to limit concurrent port tracking.
type CapConfig struct {
	Enabled bool `yaml:"enabled"`
	Max     int  `yaml:"max"`
}

// defaultCapConfig returns a safe, disabled CapConfig.
func defaultCapConfig() CapConfig {
	return CapConfig{
		Enabled: false,
		Max:     256,
	}
}

// Validate checks that the CapConfig fields are consistent.
func (c CapConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Max < 1 {
		return errors.New("cap: max must be at least 1")
	}
	return nil
}
