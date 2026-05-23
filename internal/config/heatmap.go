package config

import "errors"

// HeatMapConfig controls the behaviour of the HeatMap tracker.
type HeatMapConfig struct {
	// Enabled determines whether heat tracking is active.
	Enabled bool `yaml:"enabled"`
	// TopN is the number of hottest ports to report in summaries.
	TopN int `yaml:"top_n"`
}

// defaultHeatMapConfig returns a sensible default configuration.
func defaultHeatMapConfig() HeatMapConfig {
	return HeatMapConfig{
		Enabled: true,
		TopN:    10,
	}
}

// Validate checks that the HeatMapConfig fields are within acceptable bounds.
func (c HeatMapConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.TopN < 1 {
		return errors.New("heatmap: top_n must be at least 1")
	}
	if c.TopN > 65535 {
		return errors.New("heatmap: top_n exceeds maximum port count")
	}
	return nil
}
