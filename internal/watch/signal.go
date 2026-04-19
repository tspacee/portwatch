package watch

import (
	"errors"
	"sync"
)

// Signal is a broadcast notification mechanism that allows multiple listeners
// to wait for a named event to be fired.
type Signal struct {
	mu       sync.Mutex
	channels map[string][]chan struct{}
}

// NewSignal returns an initialised Signal.
func NewSignal() *Signal {
	return &Signal{
		channels: make(map[string][]chan struct{}),
	}
}

// Subscribe returns a channel that will be closed when the named event fires.
// The caller must consume or discard the channel to avoid goroutine leaks.
func (s *Signal) Subscribe(event string) (<-chan struct{}, error) {
	if event == "" {
		return nil, errors.New("signal: event name must not be empty")
	}
	ch := make(chan struct{}, 1)
	s.mu.Lock()
	s.channels[event] = append(s.channels[event], ch)
	s.mu.Unlock()
	return ch, nil
}

// Fire closes all subscriber channels registered under event, notifying
// every listener exactly once.
func (s *Signal) Fire(event string) error {
	if event == "" {
		return errors.New("signal: event name must not be empty")
	}
	s.mu.Lock()
	listeners := s.channels[event]
	delete(s.channels, event)
	s.mu.Unlock()
	for _, ch := range listeners {
		close(ch)
	}
	return nil
}

// Len returns the number of active subscribers for the given event.
func (s *Signal) Len(event string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.channels[event])
}
