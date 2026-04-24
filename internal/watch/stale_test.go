package watch

import (
	"testing"
	"time"
)

func TestNewStale_InvalidTTL(t *testing.T) {
	_, err := NewStale(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
	_, err = NewStale(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative TTL")
	}
}

func TestNewStale_Valid(t *testing.T) {
	s, err := NewStale(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Stale")
	}
}

func TestStale_IsStale_NeverSeen(t *testing.T) {
	s, _ := NewStale(time.Minute)
	if !s.IsStale(8080) {
		t.Error("expected unseen port to be stale")
	}
}

func TestStale_Refresh_MarksNotStale(t *testing.T) {
	s, _ := NewStale(time.Minute)
	s.Refresh(8080)
	if s.IsStale(8080) {
		t.Error("expected refreshed port to not be stale")
	}
}

func TestStale_IsStale_AfterTTL(t *testing.T) {
	s, _ := NewStale(time.Second)
	now := time.Now()
	s.seen[9000] = now.Add(-2 * time.Second)
	if !s.IsStale(9000) {
		t.Error("expected port to be stale after TTL exceeded")
	}
}

func TestStale_Evict_RemovesEntry(t *testing.T) {
	s, _ := NewStale(time.Minute)
	s.Refresh(443)
	s.Evict(443)
	if !s.IsStale(443) {
		t.Error("expected evicted port to be stale")
	}
}

func TestStale_Snapshot_ReturnsCopy(t *testing.T) {
	s, _ := NewStale(time.Minute)
	s.Refresh(80)
	s.Refresh(443)
	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	// Mutating snapshot must not affect internal state.
	delete(snap, 80)
	if s.IsStale(80) {
		t.Error("snapshot mutation affected internal state")
	}
}
