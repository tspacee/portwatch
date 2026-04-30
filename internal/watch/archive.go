package watch

import (
	"errors"
	"sync"
	"time"
)

// ArchiveEntry holds a timestamped snapshot of open ports.
type ArchiveEntry struct {
	Timestamp time.Time
	Ports     []int
}

// Archive stores a bounded, chronological history of port snapshots.
type Archive struct {
	mu      sync.RWMutex
	entries []ArchiveEntry
	maxSize int
}

// NewArchive creates an Archive with the given maximum number of entries.
func NewArchive(maxSize int) (*Archive, error) {
	if maxSize < 1 {
		return nil, errors.New("archive: maxSize must be at least 1")
	}
	return &Archive{maxSize: maxSize}, nil
}

// Store appends a new snapshot entry, evicting the oldest if capacity is exceeded.
func (a *Archive) Store(ports []int) {
	a.mu.Lock()
	defer a.mu.Unlock()

	cp := make([]int, len(ports))
	copy(cp, ports)

	entry := ArchiveEntry{Timestamp: time.Now(), Ports: cp}
	if len(a.entries) >= a.maxSize {
		a.entries = a.entries[1:]
	}
	a.entries = append(a.entries, entry)
}

// Entries returns a copy of all stored entries in chronological order.
func (a *Archive) Entries() []ArchiveEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()

	out := make([]ArchiveEntry, len(a.entries))
	for i, e := range a.entries {
		cp := make([]int, len(e.Ports))
		copy(cp, e.Ports)
		out[i] = ArchiveEntry{Timestamp: e.Timestamp, Ports: cp}
	}
	return out
}

// Len returns the number of stored entries.
func (a *Archive) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}

// Clear removes all stored entries.
func (a *Archive) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = nil
}
