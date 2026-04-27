package watch

import (
	"errors"
	"sync"
	"time"
)

// Drift tracks how much a port's observed open/close timing deviates from its
// historical baseline. A high drift score may indicate instability or churn.
type Drift struct {
	mu       sync.Mutex
	baseline map[int]time.Duration
	last     map[int]time.Time
	scores   map[int]float64
	decay    float64
}

// NewDrift creates a Drift tracker with the given exponential decay factor
// (0 < decay <= 1). A decay close to 1 weights recent observations heavily.
func NewDrift(decay float64) (*Drift, error) {
	if decay <= 0 || decay > 1 {
		return nil, errors.New("drift: decay must be in range (0, 1]")
	}
	return &Drift{
		baseline: make(map[int]time.Duration),
		last:     make(map[int]time.Time),
		scores:   make(map[int]float64),
		decay:    decay,
	}, nil
}

// Observe records an event for the given port at the given time and updates
// the drift score using an exponentially weighted moving average.
func (d *Drift) Observe(port int, at time.Time) error {
	if port < 1 || port > 65535 {
		return errors.New("drift: port out of range")
	}
	d.mu.Lock()
	defer d.mu.Unlock()

	if prev, ok := d.last[port]; ok {
		interval := at.Sub(prev)
		if base, exists := d.baseline[port]; exists {
			delta := float64(interval - base)
			if delta < 0 {
				delta = -delta
			}
			d.scores[port] = d.decay*delta + (1-d.decay)*d.scores[port]
			d.baseline[port] = time.Duration(d.decay*float64(interval) + (1-d.decay)*float64(base))
		} else {
			d.baseline[port] = interval
			d.scores[port] = 0
		}
	}
	d.last[port] = at
	return nil
}

// Score returns the current drift score for the given port. A score of 0
// means the port is behaving consistently with its baseline.
func (d *Drift) Score(port int) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.scores[port]
}

// Reset clears all state for the given port.
func (d *Drift) Reset(port int) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.baseline, port)
	delete(d.last, port)
	delete(d.scores, port)
}
