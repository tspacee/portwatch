package watch

import (
	"errors"
	"sync"
)

// Matcher maps ports to pattern labels for flexible port classification.
type Matcher struct {
	mu       sync.RWMutex
	patterns map[int][]string
}

// NewMatcher returns an empty Matcher.
func NewMatcher() *Matcher {
	return &Matcher{
		patterns: make(map[int][]string),
	}
}

// Add registers a label pattern for the given port.
// Returns an error if the port is out of range or the label is empty.
func (m *Matcher) Add(port int, label string) error {
	if port < 1 || port > 65535 {
		return errors.New("matcher: port out of range")
	}
	if label == "" {
		return errors.New("matcher: label must not be empty")
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.patterns[port] = append(m.patterns[port], label)
	return nil
}

// Match returns all labels registered for the given port.
// Returns nil if no labels are registered.
func (m *Matcher) Match(port int) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	labels, ok := m.patterns[port]
	if !ok {
		return nil
	}
	out := make([]string, len(labels))
	copy(out, labels)
	return out
}

// Has reports whether any labels are registered for the given port.
func (m *Matcher) Has(port int) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, ok := m.patterns[port]
	return ok
}

// Remove deletes all labels for the given port.
func (m *Matcher) Remove(port int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.patterns, port)
}

// Snapshot returns a copy of all registered port-to-labels mappings.
func (m *Matcher) Snapshot() map[int][]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make(map[int][]string, len(m.patterns))
	for port, labels := range m.patterns {
		cp := make([]string, len(labels))
		copy(cp, labels)
		out[port] = cp
	}
	return out
}
