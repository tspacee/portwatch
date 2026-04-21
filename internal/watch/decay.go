package watch

import (
	"errors"
	"sync"
	"time"
)

// Decay tracks per-port scores that diminish over time toward zero.
// Each call to Tick reduces all scores by the configured rate.
type Decay struct {
	mu     sync.Mutex
	scores map[int]float64
	rate   float64
	last   time.Time
	now    func() time.Time
}

// NewDecay creates a Decay tracker where rate is the fraction lost per second (0 < rate <= 1).
func NewDecay(rate float64) (*Decay, error) {
	if rate <= 0 || rate > 1 {
		return nil, errors.New("decay: rate must be in range (0, 1]")
	}
	return &Decay{
		scores: make(map[int]float64),
		rate:   rate,
		now:    time.Now,
	}, nil
}

// Add increases the score for the given port by delta.
func (d *Decay) Add(port int, delta float64) error {
	if port < 1 || port > 65535 {
		return errors.New("decay: port out of range")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	d.scores[port] += delta
	return nil
}

// Score returns the current score for the given port.
func (d *Decay) Score(port int) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.scores[port]
}

// Tick applies exponential decay based on elapsed time since the last tick.
// Entries that decay below 0.001 are pruned.
func (d *Decay) Tick() {
	d.mu.Lock()
	defer d.mu.Unlock()
	now := d.now()
	if !d.last.IsZero() {
		elapsed := now.Sub(d.last).Seconds()
		factor := 1.0 - d.rate*elapsed
		if factor < 0 {
			factor = 0
		}
		for port, score := range d.scores {
			newScore := score * factor
			if newScore < 0.001 {
				delete(d.scores, port)
			} else {
				d.scores[port] = newScore
			}
		}
	}
	d.last = now
}

// Snapshot returns a copy of all current scores.
func (d *Decay) Snapshot() map[int]float64 {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make(map[int]float64, len(d.scores))
	for k, v := range d.scores {
		out[k] = v
	}
	return out
}
