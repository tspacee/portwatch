package watch

import (
	"fmt"
	"sync"
	"time"
)

// AuditEntry records a single port state change event for audit purposes.
type AuditEntry struct {
	Timestamp time.Time
	Port      int
	Protocol  string
	Event     string // "opened" or "closed"
	Source    string
}

// AuditLog maintains an in-memory ordered audit trail of port events.
type AuditLog struct {
	mu      sync.RWMutex
	entries []AuditEntry
	maxSize int
}

// NewAuditLog creates an AuditLog with the given maximum size.
// Returns an error if maxSize is less than 1.
func NewAuditLog(maxSize int) (*AuditLog, error) {
	if maxSize < 1 {
		return nil, fmt.Errorf("audit: maxSize must be at least 1, got %d", maxSize)
	}
	return &AuditLog{
		entries: make([]AuditEntry, 0, maxSize),
		maxSize: maxSize,
	}, nil
}

// Add appends an entry to the audit log, evicting the oldest if at capacity.
func (a *AuditLog) Add(entry AuditEntry) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.entries) >= a.maxSize {
		a.entries = a.entries[1:]
	}
	a.entries = append(a.entries, entry)
}

// Entries returns a copy of all current audit entries.
func (a *AuditLog) Entries() []AuditEntry {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := make([]AuditEntry, len(a.entries))
	copy(out, a.entries)
	return out
}

// Len returns the number of entries currently stored.
func (a *AuditLog) Len() int {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return len(a.entries)
}

// Clear removes all entries from the audit log.
func (a *AuditLog) Clear() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.entries = a.entries[:0]
}
