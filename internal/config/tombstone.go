package config

import (
	"fmt"
	"time"
)

// TombstoneConfig holds configuration for the Tombstone feature.
type TombstoneConfig struct {
	Enabled bool          `yaml:"enabled"`
	TTL     time.Duration `yaml:"ttl"`
	Ports   []int         `yaml:"ports"`
}

// defaultTombstoneConfig returns a safe default TombstoneConfig.
func defaultTombstoneConfig() TombstoneConfig {
	return TombstoneConfig{
		Enabled: false,
		TTL:     30 * time.Minute,
		Ports:   []int{},
	}
}

// Validate checks that the TombstoneConfig values are valid.
func (c TombstoneConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.TTL <= 0 {
		return fmt.Errorf("tombstone: ttl must be positive, got %s", c.TTL)
	}
	for _, p := range c.Ports {
		if p < 1 || p > 65535 {
			return fmt.Errorf("tombstone: invalid port %d", p)
		}
	}
	return nil
}
