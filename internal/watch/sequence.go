package watch

import (
	"errors"
	"sync"
)

// Sequence tracks a monotonically increasing counter per port,
// allowing callers to detect ordering violations or gaps.
type Sequence struct {
	mu      sync.Mutex
	counters map[int]uint64
}

// NewSequence returns an empty Sequence.
func NewSequence() *Sequence {
	return &Sequence{
		counters: make(map[int]uint64),
	}
}

// Next increments and returns the next sequence number for port.
// Returns an error if port is out of range.
func (s *Sequence) Next(port int) (uint64, error) {
	if port < 1 || port > 65535 {
		return 0, errors.New("sequence: port out of range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counters[port]++
	return s.counters[port], nil
}

// Current returns the current sequence number for port without incrementing.
// Returns 0 if port has never been seen.
func (s *Sequence) Current(port int) uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counters[port]
}

// Reset clears the counter for port.
func (s *Sequence) Reset(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.counters, port)
}

// Snapshot returns a copy of all current counters.
func (s *Sequence) Snapshot() map[int]uint64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[int]uint64, len(s.counters))
	for k, v := range s.counters {
		out[k] = v
	}
	return out
}
