package watch

import (
	"errors"
	"sync"
)

// Quorum tracks port observations across multiple scan sources and
// considers a port "confirmed" only when it has been seen by at
// least a configurable number of distinct sources.
type Quorum struct {
	mu       sync.Mutex
	minVotes int
	votes    map[int]map[string]struct{}
}

// NewQuorum creates a Quorum that requires at least minVotes distinct
// source names before a port is considered confirmed.
func NewQuorum(minVotes int) (*Quorum, error) {
	if minVotes < 1 {
		return nil, errors.New("quorum: minVotes must be at least 1")
	}
	return &Quorum{
		minVotes: minVotes,
		votes:    make(map[int]map[string]struct{}),
	}, nil
}

// Vote records that source has observed port. It returns true if the
// port has now reached quorum.
func (q *Quorum) Vote(port int, source string) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("quorum: port out of range")
	}
	if source == "" {
		return false, errors.New("quorum: source must not be empty")
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if _, ok := q.votes[port]; !ok {
		q.votes[port] = make(map[string]struct{})
	}
	q.votes[port][source] = struct{}{}
	return len(q.votes[port]) >= q.minVotes, nil
}

// Confirmed returns true if port has reached quorum.
func (q *Quorum) Confirmed(port int) bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.votes[port]) >= q.minVotes
}

// Reset clears all votes for port.
func (q *Quorum) Reset(port int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	delete(q.votes, port)
}

// VoteCount returns the number of distinct sources that have voted for port.
func (q *Quorum) VoteCount(port int) int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.votes[port])
}
