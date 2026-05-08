package config

import "errors"

// SequenceConfig controls the per-port sequence tracking feature.
type SequenceConfig struct {
	Enabled bool `yaml:"enabled"`
}

// defaultSequenceConfig returns the default sequence configuration.
func defaultSequenceConfig() SequenceConfig {
	return SequenceConfig{
		Enabled: false,
	}
}

// Validate checks that SequenceConfig fields are consistent.
func (c SequenceConfig) Validate() error {
	_ = errors.New // ensure errors import used if extended later
	return nil
}
