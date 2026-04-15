package config

import (
	"errors"
	"fmt"
)

// FilterConfig holds port-filtering settings loaded from the config file.
type FilterConfig struct {
	// ExcludePorts lists individual ports that should never be reported.
	ExcludePorts []int `yaml:"exclude_ports"`

	// RangeMin and RangeMax constrain the ports that are scanned.
	// Both default to the global scan range when zero.
	RangeMin int `yaml:"range_min"`
	RangeMax int `yaml:"range_max"`
}

// defaultFilterConfig returns a FilterConfig with sensible defaults.
func defaultFilterConfig() FilterConfig {
	return FilterConfig{
		ExcludePorts: []int{},
		RangeMin:     1,
		RangeMax:     65535,
	}
}

// validate checks that the FilterConfig values are self-consistent.
func (f FilterConfig) validate() error {
	for _, p := range f.ExcludePorts {
		if p < 1 || p > 65535 {
			return fmt.Errorf("filter: excluded port %d is out of range [1, 65535]", p)
		}
	}
	if f.RangeMin < 1 || f.RangeMax > 65535 {
		return errors.New("filter: port range must be between 1 and 65535")
	}
	if f.RangeMin > f.RangeMax {
		return fmt.Errorf("filter: range_min (%d) must not exceed range_max (%d)", f.RangeMin, f.RangeMax)
	}
	return nil
}
