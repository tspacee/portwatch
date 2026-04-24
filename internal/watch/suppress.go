package watch

import (
	"errors"
	"sync"
	"time"
)

// Suppress tracks ports that should be temporarily silenced from alerting.
// Once a port is suppressed, it remains quiet until the suppression expires.
type Suppress struct {
	mu      sync.Mutex
	entries map[int]time.Time
	window  time.Duration
}

// NewSuppress creates a Suppress with the given quiet window duration.
// Returns an error if window is zero or negative.
func NewSuppress(window time.Duration) (*Suppress, error) {
	if window <= 0 {
		return nil, errors.New("suppress: window must be positive")
	}
	return &Suppress{
		entries: make(map[int]time.Time),
		window:  window,
	}, nil
}

// Mute marks a port as suppressed, starting the quiet window now.
// Returns an error if port is out of range.
func (s *Suppress) Mute(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("suppress: port out of range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[port] = time.Now().Add(s.window)
	return nil
}

// IsMuted reports whether the given port is currently suppressed.
func (s *Suppress) IsMuted(port int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	expiry, ok := s.entries[port]
	if !ok {
		return false
	}
	if time.Now().After(expiry) {
		delete(s.entries, port)
		return false
	}
	return true
}

// Unmute removes a suppression for the given port immediately.
func (s *Suppress) Unmute(port int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, port)
}

// Len returns the number of currently active suppressions.
func (s *Suppress) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	count := 0
	for port, expiry := range s.entries {
		if now.After(expiry) {
			delete(s.entries, port)
		} else {
			count++
		}
	}
	return count
}
