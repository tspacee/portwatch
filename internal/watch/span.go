package watch

import (
	"errors"
	"sync"
	"time"
)

// Span tracks the active duration of a port — how long it has been
// continuously observed as open. Once a port disappears, its span is reset.
type Span struct {
	mu      sync.Mutex
	started map[int]time.Time
}

// NewSpan returns an initialised Span tracker.
func NewSpan() *Span {
	return &Span{
		started: make(map[int]time.Time),
	}
}

// Observe marks a port as currently open. If the port is seen for the first
// time (or after a reset), the start time is recorded.
func (s *Span) Observe(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("span: port out of range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.started[port]; !ok {
		s.started[port] = time.Now()
	}
	return nil
}

// Reset removes the port from tracking, clearing its start time.
func (s *Span) Reset(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("span: port out of range")
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.started, port)
	return nil
}

// Duration returns how long the port has been continuously open.
// Returns 0 and false if the port is not currently tracked.
func (s *Span) Duration(port int) (time.Duration, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.started[port]
	if !ok {
		return 0, false
	}
	return time.Since(t), true
}

// Len returns the number of ports currently being tracked.
func (s *Span) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.started)
}
