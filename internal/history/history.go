// Package history provides a rolling history of port scan events
// for trend analysis and repeated-violation detection.
package history

import (
	"sync"
	"time"
)

// Entry records a single scan event and any violations detected.
type Entry struct {
	Timestamp  time.Time
	OpenPorts  []int
	Violations []string
}

// History holds a bounded, in-memory ring of scan entries.
type History struct {
	mu      sync.RWMutex
	entries []Entry
	maxSize int
}

// New creates a History that retains at most maxSize entries.
// If maxSize is <= 0 it defaults to 100.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 100
	}
	return &History{
		entries: make([]Entry, 0, maxSize),
		maxSize: maxSize,
	}
}

// Add appends a new entry, evicting the oldest when the buffer is full.
func (h *History) Add(e Entry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, e)
}

// Entries returns a shallow copy of all stored entries, oldest first.
func (h *History) Entries() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Len returns the current number of stored entries.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}

// ViolationCount returns the total number of violations recorded across
// all entries in the history window.
func (h *History) ViolationCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	total := 0
	for _, e := range h.entries {
		total += len(e.Violations)
	}
	return total
}
