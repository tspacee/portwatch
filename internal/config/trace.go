package config

import "errors"

// TraceConfig controls the per-port observation trace buffer.
type TraceConfig struct {
	Enabled    bool `yaml:"enabled"`
	MaxEntries int  `yaml:"max_entries"`
}

// defaultTraceConfig returns a TraceConfig with sensible defaults.
func defaultTraceConfig() TraceConfig {
	return TraceConfig{
		Enabled:    false,
		MaxEntries: 100,
	}
}

// Validate checks that the TraceConfig fields are consistent.
func (c TraceConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MaxEntries < 1 {
		return errors.New("trace: max_entries must be at least 1 when enabled")
	}
	return nil
}
