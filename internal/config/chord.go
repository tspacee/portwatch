package config

import "errors"

// ChordConfig defines the configuration for the Chord port-set detector.
type ChordConfig struct {
	Enabled  bool  `yaml:"enabled"`
	Required []int `yaml:"required"`
}

// defaultChordConfig returns a disabled ChordConfig with no required ports.
func defaultChordConfig() ChordConfig {
	return ChordConfig{
		Enabled:  false,
		Required: []int{},
	}
}

// Validate checks that the ChordConfig is well-formed.
func (c ChordConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if len(c.Required) == 0 {
		return errors.New("chord: at least one required port must be specified when enabled")
	}
	for _, p := range c.Required {
		if p < 1 || p > 65535 {
			return errors.New("chord: required port out of valid range 1-65535")
		}
	}
	return nil
}
