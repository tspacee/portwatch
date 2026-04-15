package portscanner

import (
	"fmt"
	"net"
	"sync"
	"time"
)

// PortState represents the state of a single port.
type PortState struct {
	Port     int
	Protocol string
	Open     bool
}

// ScanResult holds the results of a full port scan.
type ScanResult struct {
	Timestamp time.Time
	Ports     []PortState
}

// Scanner scans local open ports within a given range.
type Scanner struct {
	StartPort int
	EndPort   int
	Protocol  string
	Timeout   time.Duration
	Workers   int
}

// NewScanner creates a Scanner with sensible defaults.
func NewScanner(start, end int) *Scanner {
	return &Scanner{
		StartPort: start,
		EndPort:   end,
		Protocol:  "tcp",
		Timeout:   500 * time.Millisecond,
		Workers:   100,
	}
}

// Scan performs a concurrent port scan and returns the result.
func (s *Scanner) Scan() (*ScanResult, error) {
	if s.StartPort < 1 || s.EndPort > 65535 || s.StartPort > s.EndPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", s.StartPort, s.EndPort)
	}

	ports := make(chan int, s.Workers)
	var mu sync.Mutex
	var wg sync.WaitGroup
	var results []PortState

	for i := 0; i < s.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range ports {
				address := fmt.Sprintf("127.0.0.1:%d", port)
				conn, err := net.DialTimeout(s.Protocol, address, s.Timeout)
				open := err == nil
				if open {
					conn.Close()
				}
				mu.Lock()
				results = append(results, PortState{Port: port, Protocol: s.Protocol, Open: open})
				mu.Unlock()
			}
		}()
	}

	for port := s.StartPort; port <= s.EndPort; port++ {
		ports <- port
	}
	close(ports)
	wg.Wait()

	return &ScanResult{
		Timestamp: time.Now(),
		Ports:     results,
	}, nil
}

// OpenPorts filters and returns only open ports from a ScanResult.
func OpenPorts(result *ScanResult) []PortState {
	var open []PortState
	for _, p := range result.Ports {
		if p.Open {
			open = append(open, p)
		}
	}
	return open
}
