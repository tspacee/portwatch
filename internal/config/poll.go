package config

import (
	"errors"
	"time"
)

// PollEntry defines a custom polling interval for a specific port.
type PollEntry struct {
	Port     int           `yaml:"port"`
	Interval time.Duration `yaml:"interval"`
}

// PollConfig holds configuration for the per-port poll scheduler.
type PollConfig struct {
	Enabled         bool          `yaml:"enabled"`
	DefaultInterval time.Duration `yaml:"default_interval"`
	Ports           []PollEntry   `yaml:"ports"`
}

// defaultPollConfig returns a PollConfig with sensible defaults.
func defaultPollConfig() PollConfig {
	return PollConfig{
		Enabled:         false,
		DefaultInterval: 30 * time.Second,
		Ports:           nil,
	}
}

// Validate checks that the PollConfig fields are within acceptable bounds.
func (c PollConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.DefaultInterval <= 0 {
		return errors.New("poll: default_interval must be positive")
	}
	for _, e := range c.Ports {
		if e.Port < 1 || e.Port > 65535 {
			return errors.New("poll: port out of range")
		}
		if e.Interval <= 0 {
			return errors.New("poll: port interval must be positive")
		}
	}
	return nil
}
