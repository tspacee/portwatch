package watch

import (
	"errors"
	"net"
	"strconv"
	"time"
)

// Probe attempts a TCP dial against a specific port to verify it is reachable.
type Probe struct {
	timeout time.Duration
}

// NewProbe returns a Probe with the given timeout.
// Returns an error if timeout is non-positive.
func NewProbe(timeout time.Duration) (*Probe, error) {
	if timeout <= 0 {
		return nil, errors.New("probe: timeout must be positive")
	}
	return &Probe{timeout: timeout}, nil
}

// Check dials the given port on localhost and returns true if the connection
// succeeds within the configured timeout.
func (p *Probe) Check(port int) bool {
	addr := net.JoinHostPort("127.0.0.1", strconv.Itoa(port))
	conn, err := net.DialTimeout("tcp", addr, p.timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

// CheckAll returns the subset of ports from the provided list that are
// reachable according to Check.
func (p *Probe) CheckAll(ports []int) []int {
	open := make([]int, 0, len(ports))
	for _, port := range ports {
		if p.Check(port) {
			open = append(open, port)
		}
	}
	return open
}
