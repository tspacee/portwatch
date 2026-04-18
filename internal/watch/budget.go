package watch

import (
	"errors"
	"sync"
	"time"
)

// Budget tracks how much scan time has been consumed within a rolling window
// and rejects new scans when the budget is exhausted.
type Budget struct {
	mu       sync.Mutex
	window   time.Duration
	max      time.Duration
	entries  []budgetEntry
}

type budgetEntry struct {
	at  time.Time
	dur time.Duration
}

// NewBudget creates a Budget that allows at most maxUsage of scan time per window.
func NewBudget(window, maxUsage time.Duration) (*Budget, error) {
	if window <= 0 {
		return nil, errors.New("budget: window must be positive")
	}
	if maxUsage <= 0 {
		return nil, errors.New("budget: maxUsage must be positive")
	}
	return &Budget{window: window, max: maxUsage}, nil
}

// Record registers a completed scan duration.
func (b *Budget) Record(d time.Duration) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.evict()
	b.entries = append(b.entries, budgetEntry{at: time.Now(), dur: d})
}

// Allow returns true if a new scan is permitted under the current budget.
func (b *Budget) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.evict()
	var total time.Duration
	for _, e := range b.entries {
		total += e.dur
	}
	return total < b.max
}

// Used returns the total scan time consumed in the current window.
func (b *Budget) Used() time.Duration {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.evict()
	var total time.Duration
	for _, e := range b.entries {
		total += e.dur
	}
	return total
}

func (b *Budget) evict() {
	cutoff := time.Now().Add(-b.window)
	i := 0
	for i < len(b.entries) && b.entries[i].at.Before(cutoff) {
		i++
	}
	b.entries = b.entries[i:]
}
