package config

import "fmt"

// ShadowConfig controls the shadow port-state store.
type ShadowConfig struct {
	Enabled bool `yaml:"enabled"`
}

// defaultShadowConfig returns safe defaults.
func defaultShadowConfig() ShadowConfig {
	return ShadowConfig{
		Enabled: true,
	}
}

// Validate checks ShadowConfig for correctness.
func (c ShadowConfig) Validate() error {
	// No fields to validate beyond presence; reserved for future options.
	_ = fmt.Sprintf("%v", c.Enabled)
	return nil
}
