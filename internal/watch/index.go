package watch

import (
	"fmt"
	"sync"
)

// Index maintains a mapping from port numbers to arbitrary string categories,
// allowing callers to group or classify observed ports at runtime.
type Index struct {
	mu      sync.RWMutex
	entries map[int]string
}

// NewIndex returns an empty Index ready for use.
func NewIndex() *Index {
	return &Index{
		entries: make(map[int]string),
	}
}

// Set assigns a category to a port. The port must be in the range [1, 65535]
// and category must not be empty.
func (idx *Index) Set(port int, category string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("index: port %d out of range [1, 65535]", port)
	}
	if category == "" {
		return fmt.Errorf("index: category must not be empty")
	}
	idx.mu.Lock()
	defer idx.mu.Unlock()
	idx.entries[port] = category
	return nil
}

// Get returns the category for a port and whether it was found.
func (idx *Index) Get(port int) (string, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	v, ok := idx.entries[port]
	return v, ok
}

// Delete removes a port from the index. It is a no-op if the port is absent.
func (idx *Index) Delete(port int) {
	idx.mu.Lock()
	defer idx.mu.Unlock()
	delete(idx.entries, port)
}

// Snapshot returns a copy of all current index entries.
func (idx *Index) Snapshot() map[int]string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	copy := make(map[int]string, len(idx.entries))
	for k, v := range idx.entries {
		copy[k] = v
	}
	return copy
}

// Len returns the number of entries in the index.
func (idx *Index) Len() int {
	idx.mu.RLock()
	defer idx.mu.RUnlock()
	return len(idx.entries)
}
