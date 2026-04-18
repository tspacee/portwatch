package watch

import (
	"errors"
	"fmt"
	"sync"
)

// Labeler assigns human-readable labels to ports based on registered mappings.
// Unknown ports fall back to a default label format.
type Labeler struct {
	mu     sync.RWMutex
	labels map[int]string
}

// NewLabeler returns a Labeler with an empty label map.
func NewLabeler() *Labeler {
	return &Labeler{
		labels: make(map[int]string),
	}
}

// Register associates a label with a port number.
func (l *Labeler) Register(port int, label string) error {
	if port < 1 || port > 65535 {
		return errors.New("labeler: port out of range")
	}
	if label == "" {
		return errors.New("labeler: label must not be empty")
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.labels[port] = label
	return nil
}

// Label returns the label for a port, or a default formatted string.
func (l *Labeler) Label(port int) string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	if lbl, ok := l.labels[port]; ok {
		return lbl
	}
	return fmt.Sprintf("port/%d", port)
}

// All returns a copy of all registered port-label pairs.
func (l *Labeler) All() map[int]string {
	l.mu.RLock()
	defer l.mu.RUnlock()
	out := make(map[int]string, len(l.labels))
	for k, v := range l.labels {
		out[k] = v
	}
	return out
}
