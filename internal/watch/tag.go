package watch

import (
	"errors"
	"sync"
)

// Tag associates a string label with a port number for display and filtering.
type Tag struct {
	Port  int
	Label string
}

// TagRegistry maps ports to human-readable labels.
type TagRegistry struct {
	mu   sync.RWMutex
	tags map[int]string
}

// NewTagRegistry creates an empty TagRegistry.
func NewTagRegistry() *TagRegistry {
	return &TagRegistry{tags: make(map[int]string)}
}

// Set adds or updates a label for the given port.
func (r *TagRegistry) Set(port int, label string) error {
	if port < 1 || port > 65535 {
		return errors.New("port out of range")
	}
	if label == "" {
		return errors.New("label must not be empty")
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tags[port] = label
	return nil
}

// Get returns the label for a port and whether it exists.
func (r *TagRegistry) Get(port int) (string, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	label, ok := r.tags[port]
	return label, ok
}

// Delete removes the tag for a port.
func (r *TagRegistry) Delete(port int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.tags, port)
}

// All returns a copy of all tags.
func (r *TagRegistry) All() []Tag {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Tag, 0, len(r.tags))
	for port, label := range r.tags {
		out = append(out, Tag{Port: port, Label: label})
	}
	return out
}
