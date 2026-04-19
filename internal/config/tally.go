package config

import "errors"

// TallyConfig controls the port tally feature.
type TallyConfig struct {
	Enabled  bool `yaml:"enabled"`
	TopN     int  `yaml:"top_n"`
}

// defaultTallyConfig returns a sensible default TallyConfig.
func defaultTallyConfig() TallyConfig {
	return TallyConfig{
		Enabled: true,
		TopN:    10,
	}
}

// Validate checks TallyConfig fields for correctness.
func (c TallyConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.TopN < 1 {
		return errors.New("tally: top_n must be at least 1")
	}
	if c.TopN > 65535 {
		return errors.New("tally: top_n exceeds maximum port count")
	}
	return nil
}
