package watch

import "errors"

// ErrInvalidConcurrency is returned when the concurrency limit is less than 1.
var ErrInvalidConcurrency = errors.New("watch: concurrency limit must be at least 1")

// Semaphore limits the number of concurrent operations.
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore creates a Semaphore with the given concurrency limit.
func NewSemaphore(limit int) (*Semaphore, error) {
	if limit < 1 {
		return nil, ErrInvalidConcurrency
	}
	return &Semaphore{ch: make(chan struct{}, limit)}, nil
}

// Acquire blocks until a slot is available.
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release frees a slot.
func (s *Semaphore) Release() {
	<-s.ch
}

// TryAcquire attempts to acquire a slot without blocking.
// Returns true if successful.
func (s *Semaphore) TryAcquire() bool {
	select {
	case s.ch <- struct{}{}:
		return true
	default:
		return false
	}
}

// Available returns the number of free slots.
func (s *Semaphore) Available() int {
	return cap(s.ch) - len(s.ch)
}
