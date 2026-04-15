package metrics_test

import (
	"testing"
	"time"

	"github.com/user/portwatch/internal/metrics"
)

func TestNewCollector_EmptySummary(t *testing.T) {
	c := metrics.NewCollector()
	s := c.Summary()

	if s.TotalScans != 0 {
		t.Fatalf("expected 0 total scans, got %d", s.TotalScans)
	}
	if s.TotalViolations != 0 {
		t.Fatalf("expected 0 violations, got %d", s.TotalViolations)
	}
}

func TestCollector_Record_IncrementsTotals(t *testing.T) {
	c := metrics.NewCollector()

	c.Record(metrics.ScanResult{
		Timestamp:    time.Now(),
		OpenPorts:    5,
		Violations:   2,
		ScanDuration: 10 * time.Millisecond,
	})
	c.Record(metrics.ScanResult{
		Timestamp:    time.Now(),
		OpenPorts:    3,
		Violations:   1,
		ScanDuration: 20 * time.Millisecond,
	})

	s := c.Summary()
	if s.TotalScans != 2 {
		t.Fatalf("expected 2 scans, got %d", s.TotalScans)
	}
	if s.TotalViolations != 3 {
		t.Fatalf("expected 3 violations, got %d", s.TotalViolations)
	}
}

func TestCollector_RecordError_IncrementsErrorCount(t *testing.T) {
	c := metrics.NewCollector()
	c.RecordError()
	c.RecordError()

	if c.Summary().ErroredScans != 2 {
		t.Fatalf("expected 2 errored scans")
	}
}

func TestCollector_Summary_AvgDuration(t *testing.T) {
	c := metrics.NewCollector()
	c.Record(metrics.ScanResult{ScanDuration: 10 * time.Millisecond})
	c.Record(metrics.ScanResult{ScanDuration: 30 * time.Millisecond})

	s := c.Summary()
	if s.AvgScanDuration != 20*time.Millisecond {
		t.Fatalf("expected 20ms avg, got %v", s.AvgScanDuration)
	}
}

func TestCollector_Summary_LastScanTracked(t *testing.T) {
	c := metrics.NewCollector()

	old := time.Now().Add(-1 * time.Minute)
	recent := time.Now()

	c.Record(metrics.ScanResult{Timestamp: old, OpenPorts: 1})
	c.Record(metrics.ScanResult{Timestamp: recent, OpenPorts: 7})

	s := c.Summary()
	if !s.LastScan.Equal(recent) {
		t.Fatalf("expected last scan to be the most recent timestamp")
	}
	if s.LastOpenPorts != 7 {
		t.Fatalf("expected 7 open ports from last scan, got %d", s.LastOpenPorts)
	}
}
