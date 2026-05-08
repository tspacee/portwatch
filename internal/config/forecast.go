package config

import "errors"

// ForecastConfig controls the Forecast rolling-window predictor.
type ForecastConfig struct {
	Enabled   bool    `yaml:"enabled"`
	Window    int     `yaml:"window"`
	Threshold float64 `yaml:"threshold"`
}

// defaultForecastConfig returns a safe default configuration.
func defaultForecastConfig() ForecastConfig {
	return ForecastConfig{
		Enabled:   false,
		Window:    10,
		Threshold: 0.75,
	}
}

// Validate checks that the ForecastConfig fields are within acceptable ranges.
func (c ForecastConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.Window < 1 {
		return errors.New("forecast: window must be at least 1")
	}
	if c.Threshold < 0 || c.Threshold > 1 {
		return errors.New("forecast: threshold must be between 0.0 and 1.0")
	}
	return nil
}
