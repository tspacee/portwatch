package metrics

import (
	"sync"
	"time"
)

// ScanResult holds statistics for a single scan cycle.
type ScanResult struct {
	Timestamp    time.Time
	OpenPorts    int
	Violations   int
	ScanDuration time.Duration
}

// Collector accumulates scan metrics over the lifetime of the daemon.
type Collector struct {
	mu      sync.RWMutex
	scans   []ScanResult
	total   int
	errored int
}

// NewCollector returns an initialised Collector.
func NewCollector() *Collector {
	return &Collector{}
}

// Record appends a scan result to the collector.
func (c *Collector) Record(r ScanResult) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.scans = append(c.scans, r)
	c.total++
}

// RecordError increments the errored-scan counter.
func (c *Collector) RecordError() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.errored++
}

// Summary returns an aggregated view of all recorded scans.
func (c *Collector) Summary() Summary {
	c.mu.RLock()
	defer c.mu.RUnlock()

	s := Summary{
		TotalScans:   c.total,
		ErroredScans: c.errored,
	}

	var totalDur time.Duration
	for _, r := range c.scans {
		s.TotalViolations += r.Violations
		totalDur += r.ScanDuration
		if r.Timestamp.After(s.LastScan) {
			s.LastScan = r.Timestamp
			s.LastOpenPorts = r.OpenPorts
		}
	}

	if c.total > 0 {
		s.AvgScanDuration = totalDur / time.Duration(c.total)
	}

	return s
}

// Summary is a snapshot of aggregated metrics.
type Summary struct {
	TotalScans      int
	ErroredScans    int
	TotalViolations int
	LastScan        time.Time
	LastOpenPorts   int
	AvgScanDuration time.Duration
}
