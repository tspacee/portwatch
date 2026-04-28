package watch

import (
	"testing"
)

func TestNewCensus_Empty(t *testing.T) {
	c := NewCensus()
	if c.Scans() != 0 {
		t.Fatalf("expected 0 scans, got %d", c.Scans())
	}
	if len(c.Snapshot()) != 0 {
		t.Fatal("expected empty snapshot")
	}
}

func TestCensus_Record_IncrementsScanCount(t *testing.T) {
	c := NewCensus()
	_ = c.Record([]int{80, 443})
	_ = c.Record([]int{80})
	if c.Scans() != 2 {
		t.Fatalf("expected 2 scans, got %d", c.Scans())
	}
}

func TestCensus_Record_TracksFrequency(t *testing.T) {
	c := NewCensus()
	_ = c.Record([]int{80, 443})
	_ = c.Record([]int{80})
	if c.Frequency(80) != 2 {
		t.Fatalf("expected frequency 2 for port 80, got %d", c.Frequency(80))
	}
	if c.Frequency(443) != 1 {
		t.Fatalf("expected frequency 1 for port 443, got %d", c.Frequency(443))
	}
}

func TestCensus_Frequency_UnseenPort_ReturnsZero(t *testing.T) {
	c := NewCensus()
	if c.Frequency(8080) != 0 {
		t.Fatal("expected 0 for unseen port")
	}
}

func TestCensus_Record_InvalidPort_ReturnsError(t *testing.T) {
	c := NewCensus()
	if err := c.Record([]int{0}); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := c.Record([]int{65536}); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestCensus_Snapshot_ReturnsCopy(t *testing.T) {
	c := NewCensus()
	_ = c.Record([]int{22, 80})
	snap := c.Snapshot()
	snap[22] = 999
	if c.Frequency(22) == 999 {
		t.Fatal("snapshot mutation affected census")
	}
}

func TestCensus_Reset_ClearsState(t *testing.T) {
	c := NewCensus()
	_ = c.Record([]int{80})
	c.Reset()
	if c.Scans() != 0 {
		t.Fatalf("expected 0 scans after reset, got %d", c.Scans())
	}
	if c.Frequency(80) != 0 {
		t.Fatalf("expected 0 frequency after reset, got %d", c.Frequency(80))
	}
}

func TestCensus_Record_EmptyPorts_IncrementsScanOnly(t *testing.T) {
	c := NewCensus()
	_ = c.Record([]int{})
	if c.Scans() != 1 {
		t.Fatalf("expected 1 scan, got %d", c.Scans())
	}
	if len(c.Snapshot()) != 0 {
		t.Fatal("expected empty counts for empty port list")
	}
}
