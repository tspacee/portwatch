package watch

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// ChangeEntry records a single port-change event.
type ChangeEntry struct {
	Timestamp time.Time
	Added     []int
	Removed   []int
}

// ChangeLog keeps an in-memory ring of recent port-change events.
type ChangeLog struct {
	entries []ChangeEntry
	maxSize int
}

// NewChangeLog creates a ChangeLog with the given capacity.
func NewChangeLog(maxSize int) (*ChangeLog, error) {
	if maxSize < 1 {
		return nil, fmt.Errorf("maxSize must be at least 1, got %d", maxSize)
	}
	return &ChangeLog{maxSize: maxSize}, nil
}

// Add appends a new entry, evicting the oldest if at capacity.
func (c *ChangeLog) Add(entry ChangeEntry) {
	if len(c.entries) >= c.maxSize {
		c.entries = c.entries[1:]
	}
	c.entries = append(c.entries, entry)
}

// Entries returns a copy of all stored entries.
func (c *ChangeLog) Entries() []ChangeEntry {
	out := make([]ChangeEntry, len(c.entries))
	copy(out, c.entries)
	return out
}

// Len returns the current number of entries.
func (c *ChangeLog) Len() int { return len(c.entries) }

// Print writes a human-readable summary to w (defaults to os.Stdout).
func (c *ChangeLog) Print(w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	for _, e := range c.entries {
		added := formatInts(e.Added)
		removed := formatInts(e.Removed)
		fmt.Fprintf(w, "[%s] added=%s removed=%s\n",
			e.Timestamp.Format(time.RFC3339), added, removed)
	}
}

func formatInts(vals []int) string {
	if len(vals) == 0 {
		return "none"
	}
	parts := make([]string, len(vals))
	for i, v := range vals {
		parts[i] = fmt.Sprintf("%d", v)
	}
	return strings.Join(parts, ",")
}
