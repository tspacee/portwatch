package config

import "errors"

// SpanConfig controls the behaviour of the Span duration tracker.
type SpanConfig struct {
	// Enabled activates span tracking when true.
	Enabled bool `yaml:"enabled"`

	// MinDuration is the minimum continuous open duration (in seconds) that
	// must be reached before a port is considered "established". Zero means
	// no minimum is enforced.
	MinDurationSecs int `yaml:"min_duration_secs"`
}

// defaultSpanConfig returns a safe, conservative SpanConfig.
func defaultSpanConfig() SpanConfig {
	return SpanConfig{
		Enabled:         true,
		MinDurationSecs: 0,
	}
}

// Validate returns an error if the SpanConfig contains invalid values.
func (c SpanConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MinDurationSecs < 0 {
		return errors.New("span: min_duration_secs must be >= 0")
	}
	return nil
}
