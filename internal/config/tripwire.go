package config

import "errors"

// TripwireConfig holds configuration for the Tripwire feature.
type TripwireConfig struct {
	Enabled bool  `yaml:"enabled"`
	Ports   []int `yaml:"ports"`
}

// defaultTripwireConfig returns a safe, disabled TripwireConfig.
func defaultTripwireConfig() TripwireConfig {
	return TripwireConfig{
		Enabled: false,
		Ports:   []int{},
	}
}

// Validate checks that the TripwireConfig is internally consistent.
func (c TripwireConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if len(c.Ports) == 0 {
		return errors.New("tripwire: at least one port must be specified when enabled")
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return errors.New("tripwire: port out of valid range (1-65535)")
		}
	}
	return nil
}
