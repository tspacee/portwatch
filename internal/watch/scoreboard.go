package watch

import (
	"errors"
	"sync"
)

// Scoreboard tracks cumulative violation counts per port across scan cycles.
// It provides a ranked view of the most frequently violating ports.
type Scoreboard struct {
	mu     sync.RWMutex
	counts map[int]int
}

// NewScoreboard returns an initialised Scoreboard.
func NewScoreboard() *Scoreboard {
	return &Scoreboard{
		counts: make(map[int]int),
	}
}

// Record increments the violation count for the given port.
func (s *Scoreboard) Record(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("scoreboard: port out of range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[port]++
	return nil
}

// Score returns the current violation count for the given port.
func (s *Scoreboard) Score(port int) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.counts[port]
}

// Reset clears the violation count for the given port.
func (s *Scoreboard) Reset(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.counts, port)
}

// Snapshot returns a copy of all recorded counts keyed by port.
func (s *Scoreboard) Snapshot() map[int]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[int]int, len(s.counts))
	for k, v := range s.counts {
		out[k] = v
	}
	return out
}

// Len returns the number of ports currently tracked.
func (s *Scoreboard) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.counts)
}
