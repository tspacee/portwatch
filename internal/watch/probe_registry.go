package watch

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// ProbeRegistry tracks active port probes and their last results.
type ProbeRegistry struct {
	mu      sync.RWMutex
	probes  map[int]*ProbeRecord
	timeout time.Duration
}

// ProbeRecord holds the result of a single probe check.
type ProbeRecord struct {
	Port      int
	Open      bool
	LastCheck time.Time
	Err       error
}

// NewProbeRegistry creates a ProbeRegistry with the given probe timeout.
func NewProbeRegistry(timeout time.Duration) (*ProbeRegistry, error) {
	if timeout <= 0 {
		return nil, errors.New("probe registry: timeout must be positive")
	}
	return &ProbeRegistry{
		probes:  make(map[int]*ProbeRecord),
		timeout: timeout,
	}, nil
}

// Record stores a probe result for the given port.
func (r *ProbeRegistry) Record(port int, open bool, err error) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("probe registry: invalid port %d", port)
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.probes[port] = &ProbeRecord{
		Port:      port,
		Open:      open,
		LastCheck: time.Now(),
		Err:       err,
	}
	return nil
}

// Get returns the probe record for a port, or nil if not found.
func (r *ProbeRegistry) Get(port int) *ProbeRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if rec, ok := r.probes[port]; ok {
		copy := *rec
		return &copy
	}
	return nil
}

// Len returns the number of tracked ports.
func (r *ProbeRegistry) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.probes)
}

// Clear removes all probe records.
func (r *ProbeRegistry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.probes = make(map[int]*ProbeRecord)
}
