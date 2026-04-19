package config

import "errors"

// FenceConfig holds configuration for the port fence allowlist.
type FenceConfig struct {
	Enabled      bool  `yaml:"enabled"`
	AllowedPorts []int `yaml:"allowed_ports"`
}

func defaultFenceConfig() FenceConfig {
	return FenceConfig{
		Enabled:      false,
		AllowedPorts: []int{},
	}
}

func (f FenceConfig) Validate() error {
	if !f.Enabled {
		return nil
	}
	if len(f.AllowedPorts) == 0 {
		return errors.New("fence: allowed_ports must not be empty when enabled")
	}
	for _, p := range f.AllowedPorts {
		if p < 1 || p > 65535 {
			return errors.New("fence: allowed port out of valid range (1-65535)")
		}
	}
	return nil
}
