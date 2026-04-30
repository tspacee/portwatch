package config

import "errors"

// RosterConfig holds configuration for the port Roster feature.
type RosterConfig struct {
	// Enabled controls whether roster tracking is active.
	Enabled bool `yaml:"enabled"`

	// InitialPorts is the list of ports pre-enrolled at startup.
	InitialPorts []int `yaml:"initial_ports"`
}

// defaultRosterConfig returns a RosterConfig with sensible defaults.
func defaultRosterConfig() RosterConfig {
	return RosterConfig{
		Enabled:      true,
		InitialPorts: []int{},
	}
}

// Validate checks that all initial ports are within the valid range.
func (c RosterConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	for _, p := range c.InitialPorts {
		if p < 1 || p > 65535 {
			return errors.New("roster: initial port out of range [1, 65535]")
		}
	}
	return nil
}
