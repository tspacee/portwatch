package config

import (
	"errors"
	"time"
)

// ExpiryConfig holds configuration for the per-port expiry tracker.
type ExpiryConfig struct {
	Enabled bool          `yaml:"enabled"`
	TTL     time.Duration `yaml:"ttl"`
}

// defaultExpiryConfig returns a sensible default ExpiryConfig.
func defaultExpiryConfig() ExpiryConfig {
	return ExpiryConfig{
		Enabled: false,
		TTL:     30 * time.Second,
	}
}

// Validate checks the ExpiryConfig for correctness.
func (c ExpiryConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.TTL <= 0 {
		return errors.New("expiry: ttl must be greater than zero")
	}
	return nil
}
