package config

import "errors"

// IntentEntry maps a single port to its declared purpose.
type IntentEntry struct {
	Port   int    `yaml:"port"`
	Reason string `yaml:"reason"`
}

// IntentConfig holds operator-declared port intents loaded from config.
type IntentConfig struct {
	Enabled bool          `yaml:"enabled"`
	Ports   []IntentEntry `yaml:"ports"`
}

// defaultIntentConfig returns an IntentConfig with safe defaults.
func defaultIntentConfig() IntentConfig {
	return IntentConfig{
		Enabled: false,
		Ports:   []IntentEntry{},
	}
}

// Validate checks that all declared intent entries are well-formed.
func (c IntentConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	for _, e := range c.Ports {
		if e.Port < 1 || e.Port > 65535 {
			return errors.New("intent: port out of range")
		}
		if e.Reason == "" {
			return errors.New("intent: reason must not be empty")
		}
	}
	return nil
}
