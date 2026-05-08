package config

import "errors"

// OriginConfig controls the Origin tracker behaviour.
type OriginConfig struct {
	Enabled bool              `yaml:"enabled"`
	Sources map[int]string    `yaml:"sources"`
}

// defaultOriginConfig returns a safe default OriginConfig.
func defaultOriginConfig() OriginConfig {
	return OriginConfig{
		Enabled: false,
		Sources: map[int]string{},
	}
}

// Validate checks that the OriginConfig is coherent.
func (c OriginConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	for port, source := range c.Sources {
		if port < 1 || port > 65535 {
			return errors.New("origin: port out of range")
		}
		if source == "" {
			return errors.New("origin: source label must not be empty")
		}
	}
	return nil
}
