package watch

import (
	"errors"
	"sync"
)

// Nomad tracks ports that have moved between scan cycles — ports that
// disappeared from one range and reappeared in another within the same
// observation window.
type Nomad struct {
	mu      sync.Mutex
	prev    map[int]struct{}
	moved   map[int]int // port -> times moved
	maxMove int
}

// NewNomad creates a Nomad tracker. maxMove is the maximum number of
// moves a port may make before it is considered persistently nomadic.
func NewNomad(maxMove int) (*Nomad, error) {
	if maxMove < 1 {
		return nil, errors.New("nomad: maxMove must be at least 1")
	}
	return &Nomad{
		prev:    make(map[int]struct{}),
		moved:   make(map[int]int),
		maxMove: maxMove,
	}, nil
}

// Observe compares current ports against the previous snapshot and
// increments the move counter for any port that has reappeared after
// being absent.
func (n *Nomad) Observe(ports []int) {
	n.mu.Lock()
	defer n.mu.Unlock()

	current := make(map[int]struct{}, len(ports))
	for _, p := range ports {
		current[p] = struct{}{}
		if _, wasSeen := n.prev[p]; !wasSeen {
			n.moved[p]++
		}
	}
	n.prev = current
}

// IsNomadic returns true if the given port has exceeded the maxMove
// threshold, indicating it is persistently appearing and disappearing.
func (n *Nomad) IsNomadic(port int) bool {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.moved[port] > n.maxMove
}

// MoveCount returns the number of times the port has been observed
// reappearing after an absence.
func (n *Nomad) MoveCount(port int) int {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.moved[port]
}

// Reset clears all move history and the previous snapshot.
func (n *Nomad) Reset() {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.prev = make(map[int]struct{})
	n.moved = make(map[int]int)
}
