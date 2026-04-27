package watch

import (
	"errors"
	"sync"
	"time"
)

// Pulse tracks the heartbeat interval between successive port scan cycles,
// exposing the average, minimum, and maximum observed intervals.
type Pulse struct {
	mu      sync.Mutex
	last    time.Time
	count   int
	total   time.Duration
	min     time.Duration
	max     time.Duration
}

// PulseSummary holds a snapshot of pulse statistics.
type PulseSummary struct {
	Count int
	Avg   time.Duration
	Min   time.Duration
	Max   time.Duration
}

// NewPulse returns a new Pulse tracker.
func NewPulse() *Pulse {
	return &Pulse{}
}

// Beat records a single heartbeat. The interval is measured from the
// previous call to Beat. The first call initialises the baseline and
// does not record an interval.
func (p *Pulse) Beat() {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	if p.last.IsZero() {
		p.last = now
		return
	}

	interval := now.Sub(p.last)
	p.last = now
	p.count++
	p.total += interval

	if p.min == 0 || interval < p.min {
		p.min = interval
	}
	if interval > p.max {
		p.max = interval
	}
}

// Summary returns a snapshot of the current pulse statistics.
// Returns an error if no intervals have been recorded yet.
func (p *Pulse) Summary() (PulseSummary, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.count == 0 {
		return PulseSummary{}, errors.New("pulse: no intervals recorded")
	}

	return PulseSummary{
		Count: p.count,
		Avg:   p.total / time.Duration(p.count),
		Min:   p.min,
		Max:   p.max,
	}, nil
}

// Reset clears all recorded intervals and resets the baseline.
func (p *Pulse) Reset() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.last = time.Time{}
	p.count = 0
	p.total = 0
	p.min = 0
	p.max = 0
}
