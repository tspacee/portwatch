package watch

import (
	"testing"
	"time"
)

func TestNewStamp_Empty(t *testing.T) {
	s := NewStamp()
	if s == nil {
		t.Fatal("expected non-nil Stamp")
	}
	snap := s.Snapshot()
	if len(snap) != 0 {
		t.Fatalf("expected empty snapshot, got %d entries", len(snap))
	}
}

func TestStamp_Touch_Valid(t *testing.T) {
	s := NewStamp()
	before := time.Now()
	if err := s.Touch(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := time.Now()
	ts, ok := s.Last(8080)
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v out of expected range", ts)
	}
}

func TestStamp_Touch_InvalidPort(t *testing.T) {
	s := NewStamp()
	if err := s.Touch(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := s.Touch(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestStamp_Last_Missing(t *testing.T) {
	s := NewStamp()
	_, ok := s.Last(9999)
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestStamp_Delete_RemovesEntry(t *testing.T) {
	s := NewStamp()
	_ = s.Touch(443)
	s.Delete(443)
	_, ok := s.Last(443)
	if ok {
		t.Fatal("expected entry to be deleted")
	}
}

func TestStamp_Snapshot_ReturnsCopy(t *testing.T) {
	s := NewStamp()
	_ = s.Touch(80)
	_ = s.Touch(443)
	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	// mutating the copy must not affect the original
	delete(snap, 80)
	if _, ok := s.Last(80); !ok {
		t.Fatal("original should still contain port 80")
	}
}
