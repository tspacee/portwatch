package watch

import (
	"errors"
	"sync"
)

// Threshold tracks a numeric value per port and fires a callback
// when the value exceeds a configured limit.
type Threshold struct {
	mu       sync.Mutex
	limit    int
	counts   map[int]int
	callback func(port, value int)
}

// NewThreshold creates a Threshold with the given limit and callback.
// limit must be >= 1 and callback must not be nil.
func NewThreshold(limit int, callback func(port, value int)) (*Threshold, error) {
	if limit < 1 {
		return nil, errors.New("threshold: limit must be >= 1")
	}
	if callback == nil {
		return nil, errors.New("threshold: callback must not be nil")
	}
	return &Threshold{
		limit:    limit,
		counts:   make(map[int]int),
		callback: callback,
	}, nil
}

// Record increments the count for the given port.
// If the new count exceeds the limit, the callback is invoked.
func (t *Threshold) Record(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.counts[port]++
	if t.counts[port] > t.limit {
		t.callback(port, t.counts[port])
	}
}

// Reset clears the count for a given port.
func (t *Threshold) Reset(port int) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.counts, port)
}

// Count returns the current count for a port.
func (t *Threshold) Count(port int) int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.counts[port]
}
