package watch

import (
	"errors"
	"sync"
)

// Registry tracks the known state of ports across scans,
// associating each port with metadata such as first-seen time.

var ErrInvalidRegistryPort = errors.New("registry: port must be between 1 and 65535")

// PortMeta holds metadata about a tracked port.
type PortMeta struct {
	Port     int
	Protocol string
	SeenCount int
}

// Registry stores port metadata in a thread-safe map.
type Registry struct {
	mu      sync.RWMutex
	entries map[int]*PortMeta
}

// NewRegistry creates an empty Registry.
func NewRegistry() *Registry {
	return &Registry{
		entries: make(map[int]*PortMeta),
	}
}

// Track records a port observation, creating or updating its metadata.
func (r *Registry) Track(port int, protocol string) error {
	if port < 1 || port > 65535 {
		return ErrInvalidRegistryPort
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if m, ok := r.entries[port]; ok {
		m.SeenCount++
	} else {
		r.entries[port] = &PortMeta{Port: port, Protocol: protocol, SeenCount: 1}
	}
	return nil
}

// Get returns a copy of the metadata for a port, and whether it exists.
func (r *Registry) Get(port int) (PortMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	m, ok := r.entries[port]
	if !ok {
		return PortMeta{}, false
	}
	return *m, true
}

// Len returns the number of tracked ports.
func (r *Registry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.entries)
}

// Remove deletes a port from the registry.
func (r *Registry) Remove(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.entries, port)
}
