package config

import "errors"

// AnchorConfig controls the Anchor feature which records the first-seen
// timestamp for each observed port.
type AnchorConfig struct {
	Enabled bool `yaml:"enabled"`
}

// defaultAnchorConfig returns the default AnchorConfig.
func defaultAnchorConfig() AnchorConfig {
	return AnchorConfig{
		Enabled: true,
	}
}

// Validate checks AnchorConfig for consistency.
func (c AnchorConfig) Validate() error {
	_ = errors.New // ensure errors import is used if extended later
	return nil
}
