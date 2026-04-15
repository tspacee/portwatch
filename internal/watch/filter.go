package watch

import "errors"

// ErrNilFilter is returned when a nil filter is used.
var ErrNilFilter = errors.New("filter: filter cannot be nil")

// PortFilter decides whether a given port should be included in scan results.
type PortFilter interface {
	Allow(port int) bool
}

// RangeFilter allows only ports within an inclusive [Min, Max] range.
type RangeFilter struct {
	Min int
	Max int
}

// NewRangeFilter creates a RangeFilter, returning an error for invalid bounds.
func NewRangeFilter(min, max int) (*RangeFilter, error) {
	if min < 1 || max > 65535 {
		return nil, errors.New("filter: port range must be between 1 and 65535")
	}
	if min > max {
		return nil, errors.New("filter: min must not exceed max")
	}
	return &RangeFilter{Min: min, Max: max}, nil
}

// Allow returns true if port falls within [Min, Max].
func (f *RangeFilter) Allow(port int) bool {
	return port >= f.Min && port <= f.Max
}

// ExcludeFilter blocks a fixed set of ports.
type ExcludeFilter struct {
	excluded map[int]struct{}
}

// NewExcludeFilter creates an ExcludeFilter from a slice of port numbers.
func NewExcludeFilter(ports []int) *ExcludeFilter {
	m := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		m[p] = struct{}{}
	}
	return &ExcludeFilter{excluded: m}
}

// Allow returns true if port is NOT in the exclusion list.
func (f *ExcludeFilter) Allow(port int) bool {
	_, blocked := f.excluded[port]
	return !blocked
}

// ChainFilter applies multiple filters in order; all must allow the port.
type ChainFilter struct {
	filters []PortFilter
}

// NewChainFilter creates a ChainFilter from one or more PortFilter values.
func NewChainFilter(filters ...PortFilter) (*ChainFilter, error) {
	for _, f := range filters {
		if f == nil {
			return nil, ErrNilFilter
		}
	}
	return &ChainFilter{filters: filters}, nil
}

// Allow returns true only if every contained filter allows the port.
func (c *ChainFilter) Allow(port int) bool {
	for _, f := range c.filters {
		if !f.Allow(port) {
			return false
		}
	}
	return true
}
