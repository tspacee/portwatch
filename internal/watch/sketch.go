package watch

import (
	"errors"
	"sync"
)

// Sketch maintains a probabilistic frequency summary of observed ports
// using a simple count-min sketch approach with a fixed-size frequency table.
type Sketch struct {
	mu    sync.Mutex
	table map[int]uint64
	total uint64
}

// NewSketch creates an empty Sketch ready to record port observations.
func NewSketch() *Sketch {
	return &Sketch{
		table: make(map[int]uint64),
	}
}

// Record increments the observation count for the given port.
// Returns an error if the port is out of the valid range [1, 65535].
func (s *Sketch) Record(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("sketch: port out of range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.table[port]++
	s.total++
	return nil
}

// Estimate returns the recorded frequency count for the given port.
// Returns 0 for ports that have never been observed.
func (s *Sketch) Estimate(port int) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.table[port]
}

// Total returns the total number of observations recorded across all ports.
func (s *Sketch) Total() uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.total
}

// Reset clears all recorded observations.
func (s *Sketch) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.table = make(map[int]uint64)
	s.total = 0
}

// TopN returns the n ports with the highest observation counts.
// The result is sorted descending by count. If fewer than n ports
// have been observed, all observed ports are returned.
func (s *Sketch) TopN(n int) []int {
	s.mu.Lock()
	defer s.mu.Unlock()
	type entry struct {
		port  int
		count uint64
	}
	entries := make([]entry, 0, len(s.table))
	for p, c := range s.table {
		entries = append(entries, entry{p, c})
	}
	// simple insertion sort — table is typically small
	for i := 1; i < len(entries); i++ {
		for j := i; j > 0 && entries[j].count > entries[j-1].count; j-- {
			entries[j], entries[j-1] = entries[j-1], entries[j]
		}
	}
	if n > len(entries) {
		n = len(entries)
	}
	out := make([]int, n)
	for i := 0; i < n; i++ {
		out[i] = entries[i].port
	}
	return out
}
