package config

import "errors"

// CensusConfig controls the port census feature which tracks observation
// frequency across scans to distinguish persistent from transient ports.
type CensusConfig struct {
	Enabled           bool `yaml:"enabled"`
	// PersistenceThreshold is the minimum number of scans a port must appear in
	// before it is considered persistent. Zero disables the threshold check.
	PersistenceThreshold int `yaml:"persistence_threshold"`
}

// defaultCensusConfig returns a CensusConfig with sensible defaults.
func defaultCensusConfig() CensusConfig {
	return CensusConfig{
		Enabled:              true,
		PersistenceThreshold: 3,
	}
}

// Validate checks that the CensusConfig fields are logically consistent.
func (c CensusConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.PersistenceThreshold < 0 {
		return errors.New("census: persistence_threshold must be >= 0")
	}
	return nil
}
