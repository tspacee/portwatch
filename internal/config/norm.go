package config

import "errors"

// NormConfig holds configuration for the Norm baseline tracker.
type NormConfig struct {
	// Enabled controls whether norm deviation detection is active.
	Enabled bool `yaml:"enabled"`

	// Ports is the list of ports considered normal/expected.
	Ports []int `yaml:"ports"`

	// AutoLearn, when true, allows the norm to learn from the first scan
	// before freezing for subsequent comparisons.
	AutoLearn bool `yaml:"auto_learn"`
}

// defaultNormConfig returns a NormConfig with sensible defaults.
func defaultNormConfig() NormConfig {
	return NormConfig{
		Enabled:   false,
		Ports:     []int{},
		AutoLearn: false,
	}
}

// Validate checks that the NormConfig is well-formed.
func (c NormConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return errors.New("norm: port out of range")
		}
	}
	return nil
}
