package config

import (
	"errors"
	"time"
)

// LimiterConfig holds rate-limiting settings for the scan loop.
type LimiterConfig struct {
	// MinInterval is the minimum time between consecutive scans.
	MinInterval time.Duration `yaml:"min_interval"`
	// Window is the rolling window duration for the scan count cap.
	Window time.Duration `yaml:"window"`
	// MaxPerWindow is the maximum number of scans allowed within Window.
	MaxPerWindow int `yaml:"max_per_window"`
}

// defaultLimiterConfig returns safe defaults.
func defaultLimiterConfig() LimiterConfig {
	return LimiterConfig{
		MinInterval:  5 * time.Second,
		Window:       time.Minute,
		MaxPerWindow: 10,
	}
}

// Validate checks that all LimiterConfig fields are sensible.
func (c LimiterConfig) Validate() error {
	if c.MinInterval <= 0 {
		return errors.New("limiter.min_interval must be positive")
	}
	if c.Window <= 0 {
		return errors.New("limiter.window must be positive")
	}
	if c.MaxPerWindow < 1 {
		return errors.New("limiter.max_per_window must be at least 1")
	}
	return nil
}
