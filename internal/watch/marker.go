package watch

import (
	"errors"
	"sync"
)

// Marker tracks which ports have been manually marked for attention.
// Marks persist in memory for the lifetime of the process.
type Marker struct {
	mu    sync.RWMutex
	marks map[int]string
}

// NewMarker returns an empty Marker.
func NewMarker() *Marker {
	return &Marker{
		marks: make(map[int]string),
	}
}

// Mark associates a reason with a port. Port must be in [1, 65535].
func (m *Marker) Mark(port int, reason string) error {
	if port < 1 || port > 65535 {
		return errors.New("marker: port out of range")
	}
	if reason == "" {
		return errors.New("marker: reason must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.marks[port] = reason
	return nil
}

// Unmark removes the mark for a port. No-op if not marked.
func (m *Marker) Unmark(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.marks, port)
}

// IsMarked returns true if the port has been marked.
func (m *Marker) IsMarked(port int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.marks[port]
	return ok
}

// Reason returns the reason for a marked port, or empty string if not marked.
func (m *Marker) Reason(port int) string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.marks[port]
}

// Snapshot returns a copy of all current marks.
func (m *Marker) Snapshot() map[int]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[int]string, len(m.marks))
	for k, v := range m.marks {
		out[k] = v
	}
	return out
}
