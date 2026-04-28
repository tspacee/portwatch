package watch

import (
	"errors"
	"sync"
)

// Streak tracks consecutive scan occurrences of a port.
// A streak increments each time a port is observed and resets
// when the port is absent from a scan.
type Streak struct {
	mu      sync.Mutex
	counts  map[int]int
	max     map[int]int
}

// NewStreak creates a new Streak tracker.
func NewStreak() *Streak {
	return &Streak{
		counts: make(map[int]int),
		max:    make(map[int]int),
	}
}

// Observe records the current set of open ports for one scan cycle.
// Ports present are incremented; ports absent are reset to zero.
func (s *Streak) Observe(ports []int) error {
	if ports == nil {
		return errors.New("streak: ports must not be nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	seen := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		seen[p] = struct{}{}
	}

	// Increment streak for observed ports.
	for p := range seen {
		s.counts[p]++
		if s.counts[p] > s.max[p] {
			s.max[p] = s.counts[p]
		}
	}

	// Reset streak for ports not observed this cycle.
	for p := range s.counts {
		if _, ok := seen[p]; !ok {
			s.counts[p] = 0
		}
	}
	return nil
}

// Current returns the current consecutive streak for the given port.
func (s *Streak) Current(port int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.counts[port]
}

// Peak returns the highest streak ever recorded for the given port.
func (s *Streak) Peak(port int) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.max[port]
}

// Reset clears all streak data.
func (s *Streak) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts = make(map[int]int)
	s.max = make(map[int]int)
}
