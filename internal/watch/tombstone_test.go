package watch

import (
	"testing"
	"time"
)

func TestNewTombstone_InvalidTTL(t *testing.T) {
	_, err := NewTombstone(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
	_, err = NewTombstone(-1 * time.Second)
	if err == nil {
		t.Fatal("expected error for negative TTL")
	}
}

func TestNewTombstone_Valid(t *testing.T) {
	ts, err := NewTombstone(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ts.Len() != 0 {
		t.Fatalf("expected empty tombstone, got %d", ts.Len())
	}
}

func TestTombstone_Bury_InvalidPort(t *testing.T) {
	ts, _ := NewTombstone(time.Minute)
	if err := ts.Bury(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := ts.Bury(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestTombstone_IsBuried_True(t *testing.T) {
	ts, _ := NewTombstone(time.Minute)
	if err := ts.Bury(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ts.IsBuried(8080) {
		t.Fatal("expected port 8080 to be buried")
	}
}

func TestTombstone_IsBuried_False_NeverBuried(t *testing.T) {
	ts, _ := NewTombstone(time.Minute)
	if ts.IsBuried(9090) {
		t.Fatal("expected port 9090 to not be buried")
	}
}

func TestTombstone_Unbury_RemovesEntry(t *testing.T) {
	ts, _ := NewTombstone(time.Minute)
	_ = ts.Bury(443)
	ts.Unbury(443)
	if ts.IsBuried(443) {
		t.Fatal("expected port 443 to be unburied")
	}
}

func TestTombstone_Sweep_RemovesExpired(t *testing.T) {
	ts, _ := NewTombstone(10 * time.Millisecond)
	_ = ts.Bury(22)
	_ = ts.Bury(80)
	time.Sleep(20 * time.Millisecond)
	removed := ts.Sweep()
	if removed != 2 {
		t.Fatalf("expected 2 removed, got %d", removed)
	}
	if ts.Len() != 0 {
		t.Fatalf("expected 0 active entries after sweep, got %d", ts.Len())
	}
}

func TestTombstone_Len_CountsOnlyActive(t *testing.T) {
	ts, _ := NewTombstone(time.Minute)
	_ = ts.Bury(22)
	_ = ts.Bury(80)
	if ts.Len() != 2 {
		t.Fatalf("expected 2, got %d", ts.Len())
	}
	ts.Unbury(22)
	if ts.Len() != 1 {
		t.Fatalf("expected 1 after unbury, got %d", ts.Len())
	}
}
