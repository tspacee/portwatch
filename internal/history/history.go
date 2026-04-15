package history

import (
	"sync"
	"time"
)

const defaultMaxSize = 1000

// Entry represents a single scan-cycle record stored in history.
type Entry struct {
	Timestamp time.Time
	Added     []int
	Removed   []int
}

// History is a thread-safe, bounded ring-buffer of scan entries.
type History struct {
	mu      sync.RWMutex
	entries []Entry
	maxSize int
}

// New creates a History with the given maximum number of entries.
// If maxSize <= 0 the defaultMaxSize is used.
func New(maxSize int) *History {
	if maxSize <= 0 {
		maxSize = defaultMaxSize
	}
	return &History{maxSize: maxSize}
}

// Add appends an entry, evicting the oldest when the buffer is full.
func (h *History) Add(e Entry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.entries) >= h.maxSize {
		h.entries = h.entries[1:]
	}
	h.entries = append(h.entries, e)
}

// Len returns the current number of stored entries.
func (h *History) Len() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.entries)
}

// Entries returns a shallow copy of all stored entries.
func (h *History) Entries() []Entry {
	h.mu.RLock()
	defer h.mu.RUnlock()
	out := make([]Entry, len(h.entries))
	copy(out, h.entries)
	return out
}

// Last returns the most recent entry and true, or a zero Entry and false if
// the history is empty.
func (h *History) Last() (Entry, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	if len(h.entries) == 0 {
		return Entry{}, false
	}
	return h.entries[len(h.entries)-1], true
}

// Replace atomically replaces all entries with the provided slice.
func (h *History) Replace(entries []Entry) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.entries = make([]Entry, len(entries))
	copy(h.entries, entries)
}
