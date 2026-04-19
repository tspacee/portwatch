package watch

import (
	"testing"
)

func TestNewShadow_Empty(t *testing.T) {
	s := NewShadow()
	if s.Len() != 0 {
		t.Fatalf("expected 0, got %d", s.Len())
	}
}

func TestShadow_Update_NilPorts(t *testing.T) {
	s := NewShadow()
	if err := s.Update(nil); err == nil {
		t.Fatal("expected error for nil ports")
	}
}

func TestShadow_Update_StoresPorts(t *testing.T) {
	s := NewShadow()
	if err := s.Update([]int{80, 443}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Len() != 2 {
		t.Fatalf("expected 2, got %d", s.Len())
	}
}

func TestShadow_Contains_True(t *testing.T) {
	s := NewShadow()
	_ = s.Update([]int{8080})
	if !s.Contains(8080) {
		t.Fatal("expected port 8080 to be present")
	}
}

func TestShadow_Contains_False(t *testing.T) {
	s := NewShadow()
	_ = s.Update([]int{8080})
	if s.Contains(9090) {
		t.Fatal("expected port 9090 to be absent")
	}
}

func TestShadow_Snapshot_ReturnsCopy(t *testing.T) {
	s := NewShadow()
	_ = s.Update([]int{22, 80})
	snap := s.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2, got %d", len(snap))
	}
	snap[0] = 9999
	if s.Contains(9999) {
		t.Fatal("mutation of snapshot affected shadow")
	}
}

func TestShadow_Update_Replaces(t *testing.T) {
	s := NewShadow()
	_ = s.Update([]int{22, 80})
	_ = s.Update([]int{443})
	if s.Len() != 1 {
		t.Fatalf("expected 1 after replace, got %d", s.Len())
	}
	if !s.Contains(443) {
		t.Fatal("expected port 443")
	}
}
