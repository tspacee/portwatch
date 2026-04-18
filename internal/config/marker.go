package config

import "errors"

// MarkerEntry represents a single pre-configured port mark.
type MarkerEntry struct {
	Port   int    `yaml:"port"`
	Reason string `yaml:"reason"`
}

// MarkerConfig holds startup port marks loaded from configuration.
type MarkerConfig struct {
	Entries []MarkerEntry `yaml:"entries"`
}

func defaultMarkerConfig() MarkerConfig {
	return MarkerConfig{}
}

func (c MarkerConfig) Validate() error {
	for _, e := range c.Entries {
		if e.Port < 1 || e.Port > 65535 {
			return errors.New("marker: port out of range")
		}
		if e.Reason == "" {
			return errors.New("marker: reason must not be empty")
		}
	}
	return nil
}
