package watch

import (
	"fmt"
	"sync"
	"time"
)

// LedgerEntry records a single port observation with a timestamp and hit count.
type LedgerEntry struct {
	Port      int
	FirstSeen time.Time
	LastSeen  time.Time
	HitCount  int
}

// Ledger maintains a persistent record of every port that has been observed
// during the lifetime of a watch session. It tracks first-seen and last-seen
// timestamps along with a cumulative hit count per port.
type Ledger struct {
	mu      sync.RWMutex
	entries map[int]*LedgerEntry
}

// NewLedger creates an empty Ledger ready for recording port observations.
func NewLedger() *Ledger {
	return &Ledger{
		entries: make(map[int]*LedgerEntry),
	}
}

// Record registers a port observation. If the port has been seen before its
// LastSeen timestamp and HitCount are updated; otherwise a new entry is
// created with FirstSeen set to now.
func (l *Ledger) Record(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("ledger: port %d out of valid range [1, 65535]", port)
	}

	now := time.Now()

	l.mu.Lock()
	defer l.mu.Unlock()

	if e, ok := l.entries[port]; ok {
		e.LastSeen = now
		e.HitCount++
		return nil
	}

	l.entries[port] = &LedgerEntry{
		Port:      port,
		FirstSeen: now,
		LastSeen:  now,
		HitCount:  1,
	}
	return nil
}

// Get returns the LedgerEntry for the given port and a boolean indicating
// whether the port has ever been recorded.
func (l *Ledger) Get(port int) (LedgerEntry, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	e, ok := l.entries[port]
	if !ok {
		return LedgerEntry{}, false
	}
	return *e, true
}

// Entries returns a snapshot copy of all recorded ledger entries.
func (l *Ledger) Entries() []LedgerEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	out := make([]LedgerEntry, 0, len(l.entries))
	for _, e := range l.entries {
		out = append(out, *e)
	}
	return out
}

// Len returns the number of distinct ports recorded in the ledger.
func (l *Ledger) Len() int {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return len(l.entries)
}

// Reset clears all entries from the ledger.
func (l *Ledger) Reset() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = make(map[int]*LedgerEntry)
}
