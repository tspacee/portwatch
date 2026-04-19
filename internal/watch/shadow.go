package watch

import (
	"errors"
	"sync"
)

// Shadow maintains a secondary copy of the last known port state
// for comparison without mutating the primary snapshot.
type Shadow struct {
	mu    sync.RWMutex
	ports map[int]struct{}
}

// NewShadow returns an empty Shadow store.
func NewShadow() *Shadow {
	return &Shadow{
		ports: make(map[int]struct{}),
	}
}

// Update replaces the shadow state with the provided ports.
func (s *Shadow) Update(ports []int) error {
	if ports == nil {
		return errors.New("shadow: ports must not be nil")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ports = make(map[int]struct{}, len(ports))
	for _, p := range ports {
		s.ports[p] = struct{}{}
	}
	return nil
}

// Contains reports whether the given port exists in the shadow state.
func (s *Shadow) Contains(port int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.ports[port]
	return ok
}

// Snapshot returns a sorted copy of all shadowed ports.
func (s *Shadow) Snapshot() []int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]int, 0, len(s.ports))
	for p := range s.ports {
		out = append(out, p)
	}
	return out
}

// Len returns the number of ports in the shadow.
func (s *Shadow) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.ports)
}
