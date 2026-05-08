package watch

import (
	"errors"
	"sync"
)

// Chord tracks whether a required set of ports are all simultaneously open.
// It fires a callback when all watched ports are present in an observed scan.
type Chord struct {
	mu       sync.Mutex
	required map[int]struct{}
	callback func([]int)
}

// NewChord creates a Chord that fires cb when all ports in required are open.
// Returns an error if required is empty or cb is nil.
func NewChord(required []int, cb func([]int)) (*Chord, error) {
	if len(required) == 0 {
		return nil, errors.New("chord: required ports must not be empty")
	}
	if cb == nil {
		return nil, errors.New("chord: callback must not be nil")
	}
	for _, p := range required {
		if p < 1 || p > 65535 {
			return nil, errors.New("chord: port out of range")
		}
	}
	req := make(map[int]struct{}, len(required))
	for _, p := range required {
		req[p] = struct{}{}
	}
	return &Chord{required: req, callback: cb}, nil
}

// Observe checks the given open ports against the required set.
// If all required ports are present, the callback is invoked with matching ports.
func (c *Chord) Observe(open []int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	present := make(map[int]struct{}, len(open))
	for _, p := range open {
		present[p] = struct{}{}
	}

	matched := make([]int, 0, len(c.required))
	for p := range c.required {
		if _, ok := present[p]; !ok {
			return
		}
		matched = append(matched, p)
	}
	c.callback(matched)
}

// Len returns the number of required ports in the chord.
func (c *Chord) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.required)
}
