package watch

import (
	"errors"
	"sync"
)

// Origin tracks the first-seen timestamp and source label for each port.
// It is useful for attributing when and where a port was first observed.
type Origin struct {
	mu      sync.RWMutex
	entries map[int]OriginEntry
}

// OriginEntry holds metadata about a port's first observation.
type OriginEntry struct {
	Port   int
	Source string
}

// NewOrigin returns an initialised Origin tracker.
func NewOrigin() *Origin {
	return &Origin{
		entries: make(map[int]OriginEntry),
	}
}

// Record stores the origin of a port if it has not been seen before.
// Returns an error if the port is out of range or source is empty.
func (o *Origin) Record(port int, source string) error {
	if port < 1 || port > 65535 {
		return errors.New("origin: port out of range")
	}
	if source == "" {
		return errors.New("origin: source must not be empty")
	}
	o.mu.Lock()
	defer o.mu.Unlock()
	if _, exists := o.entries[port]; !exists {
		o.entries[port] = OriginEntry{Port: port, Source: source}
	}
	return nil
}

// Get returns the OriginEntry for the given port and whether it exists.
func (o *Origin) Get(port int) (OriginEntry, bool) {
	o.mu.RLock()
	defer o.mu.RUnlock()
	e, ok := o.entries[port]
	return e, ok
}

// Len returns the number of tracked ports.
func (o *Origin) Len() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.entries)
}

// Reset clears all recorded origins.
func (o *Origin) Reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.entries = make(map[int]OriginEntry)
}
