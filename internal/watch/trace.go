package watch

import (
	"errors"
	"sync"
	"time"
)

// Trace records a timestamped sequence of port observations for a single port,
// allowing callers to replay or inspect the recent history of state changes.
type Trace struct {
	mu      sync.Mutex
	max     int
	entries []TraceEntry
}

// TraceEntry holds a single observation recorded by Trace.
type TraceEntry struct {
	Port      int
	ObservedAt time.Time
	Open      bool
}

// NewTrace creates a Trace that retains at most maxEntries observations.
func NewTrace(maxEntries int) (*Trace, error) {
	if maxEntries < 1 {
		return nil, errors.New("trace: maxEntries must be at least 1")
	}
	return &Trace{max: maxEntries}, nil
}

// Record appends an observation for the given port.
func (t *Trace) Record(port int, open bool) error {
	if port < 1 || port > 65535 {
		return errors.New("trace: port out of range")
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.entries) >= t.max {
		t.entries = t.entries[1:]
	}
	t.entries = append(t.entries, TraceEntry{
		Port:       port,
		ObservedAt: time.Now(),
		Open:       open,
	})
	return nil
}

// Entries returns a copy of all recorded observations.
func (t *Trace) Entries() []TraceEntry {
	t.mu.Lock()
	defer t.mu.Unlock()
	out := make([]TraceEntry, len(t.entries))
	copy(out, t.entries)
	return out
}

// Len returns the current number of recorded observations.
func (t *Trace) Len() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.entries)
}

// Clear removes all recorded observations.
func (t *Trace) Clear() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.entries = t.entries[:0]
}
