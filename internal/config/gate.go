package config

import "errors"

// GateConfig controls the gate allow-list feature.
type GateConfig struct {
	Enabled      bool  `yaml:"enabled"`
	OpenByDefault bool  `yaml:"open_by_default"`
	AllowedPorts []int `yaml:"allowed_ports"`
}

// defaultGateConfig returns a safe default gate configuration.
func defaultGateConfig() GateConfig {
	return GateConfig{
		Enabled:      false,
		OpenByDefault: false,
		AllowedPorts: []int{},
	}
}

// Validate checks that the GateConfig is consistent.
func (g GateConfig) Validate() error {
	if !g.Enabled {
		return nil
	}
	for _, p := range g.AllowedPorts {
		if p < 1 || p > 65535 {
			return errors.New("gate: allowed_ports contains out-of-range port")
		}
	}
	return nil
}
