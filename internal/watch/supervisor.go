package watch

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrSupervisorStopped is returned when the supervisor has been stopped.
var ErrSupervisorStopped = errors.New("supervisor: stopped")

// WorkerFunc is a function that can be supervised.
type WorkerFunc func(ctx context.Context) error

// Supervisor restarts a worker function on failure with backoff.
type Supervisor struct {
	worker  WorkerFunc
	backoff *Backoff
	maxRestarts int
	mu       sync.Mutex
	restarts int
}

// NewSupervisor creates a Supervisor for the given worker.
// maxRestarts <= 0 means unlimited restarts.
func NewSupervisor(worker WorkerFunc, backoff *Backoff, maxRestarts int) (*Supervisor, error) {
	if worker == nil {
		return nil, errors.New("supervisor: worker must not be nil")
	}
	if backoff == nil {
		return nil, errors.New("supervisor: backoff must not be nil")
	}
	return &Supervisor{
		worker:      worker,
		backoff:     backoff,
		maxRestarts: maxRestarts,
	}, nil
}

// Restarts returns the number of times the worker has been restarted.
func (s *Supervisor) Restarts() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.restarts
}

// Run starts the worker and supervises it until ctx is cancelled or
// maxRestarts is exceeded.
func (s *Supervisor) Run(ctx context.Context) error {
	s.backoff.Reset()
	for {
		err := s.worker(ctx)
		if ctx.Err() != nil {
			return ctx.Err()
		}
		if err == nil {
			return nil
		}
		s.mu.Lock()
		s.restarts++
		current := s.restarts
		s.mu.Unlock()
		if s.maxRestarts > 0 && current >= s.maxRestarts {
			return ErrSupervisorStopped
		}
		delay := s.backoff.Next()
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(delay):
		}
	}
}
