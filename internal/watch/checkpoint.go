package watch

import (
	"errors"
	"sync"
	"time"
)

// Checkpoint records the last successfully processed scan position,
// allowing the watcher to resume from a known-good state after a restart
// or failure without reprocessing stale data.
type Checkpoint struct {
	mu        sync.RWMutex
	seq       uint64
	timestamp time.Time
	ports     []int
}

// NewCheckpoint returns an empty Checkpoint ready for use.
func NewCheckpoint() *Checkpoint {
	return &Checkpoint{}
}

// Commit saves the current scan position identified by seq, the scan
// time, and the set of open ports observed at that position.
func (c *Checkpoint) Commit(seq uint64, ts time.Time, ports []int) error {
	if seq == 0 {
		return errors.New("checkpoint: seq must be greater than zero")
	}
	copy := make([]int, len(ports))
	for i, p := range ports {
		copy[i] = p
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.seq = seq
	c.timestamp = ts
	c.ports = copy
	return nil
}

// Seq returns the sequence number of the last committed checkpoint.
func (c *Checkpoint) Seq() uint64 {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.seq
}

// Timestamp returns the time recorded at the last commit.
func (c *Checkpoint) Timestamp() time.Time {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.timestamp
}

// Ports returns a copy of the port list recorded at the last commit.
func (c *Checkpoint) Ports() []int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	out := make([]int, len(c.ports))
	copy(out, c.ports)
	return out
}

// Reset clears all checkpoint state, returning the Checkpoint to its
// initial empty condition.
func (c *Checkpoint) Reset() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.seq = 0
	c.timestamp = time.Time{}
	c.ports = nil
}
