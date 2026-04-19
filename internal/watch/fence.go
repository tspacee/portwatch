package watch

import (
	"errors"
	"sync"
)

// Fence enforces a port boundary: ports outside the allowed set are flagged.
type Fence struct {
	mu      sync.RWMutex
	allowed map[int]struct{}
}

// NewFence creates a Fence from the given allowed ports.
func NewFence(ports []int) (*Fence, error) {
	if len(ports) == 0 {
		return nil, errors.New("fence: allowed port list must not be empty")
	}
	allowed := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		if p < 1 || p > 65535 {
			return nil, errors.New("fence: port out of range")
		}
		allowed[p] = struct{}{}
	}
	return &Fence{allowed: allowed}, nil
}

// Allow reports whether port is within the allowed set.
func (f *Fence) Allow(port int) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()
	_, ok := f.allowed[port]
	return ok
}

// Add adds a port to the allowed set.
func (f *Fence) Add(port int) error {
	if port < 1 || port > 65535 {
		return errors.New("fence: port out of range")
	}
	f.mu.Lock()
	defer f.mu.Unlock()
	f.allowed[port] = struct{}{}
	return nil
}

// Remove removes a port from the allowed set.
func (f *Fence) Remove(port int) {
	f.mu.Lock()
	defer f.mu.Unlock()
	delete(f.allowed, port)
}

// Snapshot returns a copy of the current allowed ports.
func (f *Fence) Snapshot() []int {
	f.mu.RLock()
	defer f.mu.RUnlock()
	out := make([]int, 0, len(f.allowed))
	for p := range f.allowed {
		out = append(out, p)
	}
	return out
}
