package config

import "errors"

// NomadConfig controls the Nomad port-movement tracker.
type NomadConfig struct {
	Enabled  bool `yaml:"enabled"`
	MaxMove  int  `yaml:"max_move"`
}

// defaultNomadConfig returns a conservative default configuration.
func defaultNomadConfig() NomadConfig {
	return NomadConfig{
		Enabled: false,
		MaxMove: 3,
	}
}

// Validate checks that the NomadConfig fields are self-consistent.
func (c NomadConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MaxMove < 1 {
		return errors.New("nomad: max_move must be at least 1")
	}
	return nil
}
