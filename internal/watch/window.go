package watch

import (
	"errors"
	"sync"
	"time"
)

// Window tracks how many events occurred within a sliding time window.
// It is safe for concurrent use.
type Window struct {
	mu       sync.Mutex
	size     time.Duration
	timestamps []time.Time
}

// NewWindow creates a Window with the given duration.
// Returns an error if size is zero or negative.
func NewWindow(size time.Duration) (*Window, error) {
	if size <= 0 {
		return nil, errors.New("window: size must be positive")
	}
	return &Window{size: size}, nil
}

// Record adds an event at the current time.
func (w *Window) Record() {
	w.mu.Lock()
	defer w.mu.Unlock()
	now := time.Now()
	w.timestamps = append(w.timestamps, now)
	w.evict(now)
}

// Count returns the number of events within the window.
func (w *Window) Count() int {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.evict(time.Now())
	return len(w.timestamps)
}

// Reset clears all recorded events.
func (w *Window) Reset() {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.timestamps = w.timestamps[:0]
}

// evict removes timestamps older than the window size. Must be called with mu held.
func (w *Window) evict(now time.Time) {
	cutoff := now.Add(-w.size)
	i := 0
	for i < len(w.timestamps) && w.timestamps[i].Before(cutoff) {
		i++
	}
	w.timestamps = w.timestamps[i:]
}
