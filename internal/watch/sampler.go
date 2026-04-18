package watch

import (
	"errors"
	"sync"
	"time"
)

// Sampler records port scan samples over a sliding time window and
// exposes a rate (samples per second) for adaptive interval tuning.
type Sampler struct {
	mu      sync.Mutex
	window  time.Duration
	samples []time.Time
}

// NewSampler creates a Sampler with the given sliding window duration.
func NewSampler(window time.Duration) (*Sampler, error) {
	if window <= 0 {
		return nil, errors.New("sampler: window must be positive")
	}
	return &Sampler{window: window}, nil
}

// Record adds a sample at the current time.
func (s *Sampler) Record() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	s.samples = append(s.samples, now)
	s.evict(now)
}

// Rate returns the number of samples per second within the window.
func (s *Sampler) Rate() float64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict(time.Now())
	if len(s.samples) == 0 {
		return 0
	}
	return float64(len(s.samples)) / s.window.Seconds()
}

// Count returns the number of samples currently in the window.
func (s *Sampler) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evict(time.Now())
	return len(s.samples)
}

// Reset clears all samples.
func (s *Sampler) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.samples = s.samples[:0]
}

// evict removes samples older than the window. Must be called with mu held.
func (s *Sampler) evict(now time.Time) {
	cutoff := now.Add(-s.window)
	i := 0
	for i < len(s.samples) && s.samples[i].Before(cutoff) {
		i++
	}
	s.samples = s.samples[i:]
}
