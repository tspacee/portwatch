package config

import "errors"

// CanaryConfig holds configuration for the canary port monitor.
type CanaryConfig struct {
	Enabled bool  `yaml:"enabled"`
	Ports   []int `yaml:"ports"`
}

// defaultCanaryConfig returns a CanaryConfig with safe defaults.
func defaultCanaryConfig() CanaryConfig {
	return CanaryConfig{
		Enabled: false,
		Ports:   []int{},
	}
}

// Validate checks that all configured canary ports are within the valid range.
func (c CanaryConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if len(c.Ports) == 0 {
		return errors.New("canary: enabled but no ports specified")
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return errors.New("canary: port out of range")
		}
	}
	return nil
}
