package config

import "errors"

// ScoreboardConfig controls whether the violation scoreboard is active.
type ScoreboardConfig struct {
	Enabled bool `yaml:"enabled"`
	// TopN is the maximum number of ports to include in a ranked summary.
	TopN int `yaml:"top_n"`
}

// defaultScoreboardConfig returns safe defaults for the scoreboard.
func defaultScoreboardConfig() ScoreboardConfig {
	return ScoreboardConfig{
		Enabled: true,
		TopN:    10,
	}
}

// Validate checks that the scoreboard configuration is coherent.
func (c ScoreboardConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.TopN < 1 {
		return errors.New("scoreboard: top_n must be at least 1")
	}
	if c.TopN > 65535 {
		return errors.New("scoreboard: top_n exceeds maximum port count")
	}
	return nil
}
