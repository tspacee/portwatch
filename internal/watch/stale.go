package watch

import (
	"errors"
	"sync"
	"time"
)

// Stale tracks ports that have not been observed within a configurable TTL.
// Once a port exceeds the TTL without being refreshed it is considered stale.
type Stale struct {
	mu    sync.Mutex
	ttl   time.Duration
	seen  map[int]time.Time
	now   func() time.Time
}

// ErrInvalidTTL is returned when the provided TTL is not positive.
var ErrInvalidTTL = errors.New("stale: ttl must be greater than zero")

// NewStale creates a Stale tracker with the given TTL.
func NewStale(ttl time.Duration) (*Stale, error) {
	if ttl <= 0 {
		return nil, ErrInvalidTTL
	}
	return &Stale{
		ttl:  ttl,
		seen: make(map[int]time.Time),
		now:  time.Now,
	}, nil
}

// Refresh records the current time for the given port, resetting its staleness.
func (s *Stale) Refresh(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen[port] = s.now()
}

// IsStale reports whether the given port has exceeded the TTL since its last
// refresh, or has never been seen.
func (s *Stale) IsStale(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.seen[port]
	if !ok {
		return true
	}
	return s.now().Sub(t) > s.ttl
}

// Evict removes the tracking entry for the given port.
func (s *Stale) Evict(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.seen, port)
}

// Snapshot returns a copy of all tracked ports and their last-seen times.
func (s *Stale) Snapshot() map[int]time.Time {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[int]time.Time, len(s.seen))
	for k, v := range s.seen {
		out[k] = v
	}
	return out
}
