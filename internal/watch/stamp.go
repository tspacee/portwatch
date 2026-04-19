package watch

import (
	"errors"
	"sync"
	"time"
)

// Stamp records the last-seen timestamp for a given port.
// It is safe for concurrent use.
type Stamp struct {
	mu      sync.RWMutex
	entries map[int]time.Time
}

// NewStamp returns an initialised Stamp.
func NewStamp() *Stamp {
	return &Stamp{entries: make(map[int]time.Time)}
}

// Touch records the current time for the given port.
// Returns an error if the port is out of range.
func (s *Stamp) Touch(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("stamp: port out of range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = time.Now()
	return nil
}

// Last returns the last-seen time for the port and whether it exists.
func (s *Stamp) Last(port int) (time.Time, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.entries[port]
	return t, ok
}

// Delete removes the stamp entry for the given port.
func (s *Stamp) Delete(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// Snapshot returns a copy of all current stamp entries.
func (s *Stamp) Snapshot() map[int]time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[int]time.Time, len(s.entries))
	for k, v := range s.entries {
		out[k] = v
	}
	return out
}
