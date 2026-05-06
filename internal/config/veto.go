package config

import "errors"

// VetoEntry represents a single port veto rule in configuration.
type VetoEntry struct {
	Port   int    `yaml:"port"`
	Reason string `yaml:"reason"`
}

// VetoConfig holds the veto configuration block.
type VetoConfig struct {
	Enabled bool        `yaml:"enabled"`
	Ports   []VetoEntry `yaml:"ports"`
}

// defaultVetoConfig returns a disabled veto configuration.
func defaultVetoConfig() VetoConfig {
	return VetoConfig{
		Enabled: false,
		Ports:   []VetoEntry{},
	}
}

// Validate checks that all veto entries are well-formed.
func (c VetoConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	for _, e := range c.Ports {
		if e.Port < 1 || e.Port > 65535 {
			return errors.New("veto: port out of range")
		}
		if e.Reason == "" {
			return errors.New("veto: reason must not be empty")
		}
	}
	return nil
}
