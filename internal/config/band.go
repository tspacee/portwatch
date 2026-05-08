package config

import "fmt"

// BandEntry describes a single named port band.
type BandEntry struct {
	Name string `yaml:"name"`
	Min  int    `yaml:"min"`
	Max  int    `yaml:"max"`
}

// BandConfig holds the configuration for the Band classifier.
type BandConfig struct {
	Enabled bool        `yaml:"enabled"`
	Bands   []BandEntry `yaml:"bands"`
}

// defaultBandConfig returns a BandConfig pre-populated with the three
// canonical IANA port ranges.
func defaultBandConfig() BandConfig {
	return BandConfig{
		Enabled: true,
		Bands: []BandEntry{
			{Name: "system", Min: 1, Max: 1023},
			{Name: "registered", Min: 1024, Max: 49151},
			{Name: "dynamic", Min: 49152, Max: 65535},
		},
	}
}

// Validate checks that all band entries are well-formed.
func (c BandConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	for _, b := range c.Bands {
		if b.Name == "" {
			return fmt.Errorf("band entry has empty name")
		}
		if b.Min < 1 || b.Max > 65535 || b.Min > b.Max {
			return fmt.Errorf("band %q: invalid port range [%d, %d]", b.Name, b.Min, b.Max)
		}
	}
	return nil
}
