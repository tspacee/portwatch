package config

import (
	"errors"
	"time"
)

// ProbeRegistryConfig configures the probe registry feature.
type ProbeRegistryConfig struct {
	Enabled bool          `yaml:"enabled"`
	Timeout time.Duration `yaml:"timeout"`
}

// defaultProbeRegistryConfig returns a safe default configuration.
func defaultProbeRegistryConfig() ProbeRegistryConfig {
	return ProbeRegistryConfig{
		Enabled: false,
		Timeout: 2 * time.Second,
	}
}

// Validate checks that the ProbeRegistryConfig is sensible.
func (c ProbeRegistryConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Timeout <= 0 {
		return errors.New("probe_registry: timeout must be positive when enabled")
	}
	if c.Timeout > 30*time.Second {
		return errors.New("probe_registry: timeout must not exceed 30s")
	}
	return nil
}
