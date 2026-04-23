package watch

import (
	"fmt"
	"sync"
)

// Mask suppresses alerts for a set of ports for a configurable duration.
// Ports added to the mask are ignored during evaluation until the mask
// is explicitly cleared or the port is removed.
type Mask struct {
	mu    sync.RWMutex
	ports map[int]string
}

// NewMask returns an initialised, empty Mask.
func NewMask() *Mask {
	return &Mask{
		ports: make(map[int]string),
	}
}

// Add suppresses alerts for port with the given reason.
// Returns an error if the port is out of range or the reason is empty.
func (m *Mask) Add(port int, reason string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("mask: port %d out of range [1, 65535]", port)
	}
	if reason == "" {
		return fmt.Errorf("mask: reason must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ports[port] = reason
	return nil
}

// Remove lifts the mask for the given port.
func (m *Mask) Remove(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.ports, port)
}

// Masked reports whether the given port is currently masked.
func (m *Mask) Masked(port int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.ports[port]
	return ok
}

// Reason returns the suppression reason for port, or an empty string
// if the port is not masked.
func (m *Mask) Reason(port int) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.ports[port]
}

// Len returns the number of currently masked ports.
func (m *Mask) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.ports)
}

// Clear removes all masked ports.
func (m *Mask) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.ports = make(map[int]string)
}
