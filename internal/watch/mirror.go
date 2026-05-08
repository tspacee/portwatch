package watch

import (
	"errors"
	"sync"
)

// Mirror maintains a read-only copy of the last observed port set,
// allowing concurrent readers to access the most recent snapshot
// without blocking the scan loop.
type Mirror struct {
	mu    sync.RWMutex
	ports []int
	seq   uint64
}

// NewMirror returns an empty Mirror.
func NewMirror() *Mirror {
	return &Mirror{}
}

// Update replaces the mirrored port set with a copy of the provided slice.
// seq must be greater than the current sequence number.
func (m *Mirror) Update(seq uint64, ports []int) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if seq <= m.seq && m.seq != 0 {
		return errors.New("mirror: sequence number must be monotonically increasing")
	}

	copy_ := make([]int, len(ports))
	copy(copy_, ports)
	m.ports = copy_
	m.seq = seq
	return nil
}

// Ports returns a copy of the currently mirrored port list.
func (m *Mirror) Ports() []int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.ports) == 0 {
		return nil
	}
	out := make([]int, len(m.ports))
	copy(out, m.ports)
	return out
}

// Seq returns the sequence number of the most recent update.
func (m *Mirror) Seq() uint64 {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.seq
}

// Len returns the number of ports currently mirrored.
func (m *Mirror) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.ports)
}
