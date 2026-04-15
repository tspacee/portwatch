package history

import (
	"log"
	"time"
)

// Cleaner periodically removes entries from a History that exceed a maximum age.
type Cleaner struct {
	h          *History
	maxAge     time.Duration
	interval   time.Duration
	logger     *log.Logger
}

// NewCleaner creates a Cleaner for the given History. maxAge defines how old
// an entry can be before it is pruned; interval controls how often the sweep
// runs. Returns ErrNilHistory if h is nil.
func NewCleaner(h *History, maxAge, interval time.Duration, logger *log.Logger) (*Cleaner, error) {
	if h == nil {
		return nil, ErrNilHistory
	}
	if logger == nil {
		logger = log.Default()
	}
	return &Cleaner{
		h:        h,
		maxAge:   maxAge,
		interval: interval,
		logger:   logger,
	}, nil
}

// Run blocks, sweeping stale entries on every interval tick, until stop is
// closed.
func (c *Cleaner) Run(stop <-chan struct{}) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-stop:
			return
		case <-ticker.C:
			n := c.sweep()
			if n > 0 {
				c.logger.Printf("history/cleaner: pruned %d stale entries (older than %s)", n, c.maxAge)
			}
		}
	}
}

// sweep removes entries older than maxAge and returns the number removed.
func (c *Cleaner) sweep() int {
	cutoff := time.Now().Add(-c.maxAge)
	before := c.h.Len()
	entries := c.h.Entries()
	var kept []Entry
	for _, e := range entries {
		if !e.Timestamp.Before(cutoff) {
			kept = append(kept, e)
		}
	}
	if len(kept) == before {
		return 0
	}
	c.h.Replace(kept)
	return before - len(kept)
}
