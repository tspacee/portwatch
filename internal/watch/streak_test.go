package watch

import (
	"testing"
)

func TestNewStreak_Empty(t *testing.T) {
	s := NewStreak()
	if s.Current(80) != 0 {
		t.Fatal("expected zero streak for unseen port")
	}
	if s.Peak(80) != 0 {
		t.Fatal("expected zero peak for unseen port")
	}
}

func TestStreak_Observe_NilPorts(t *testing.T) {
	s := NewStreak()
	if err := s.Observe(nil); err == nil {
		t.Fatal("expected error for nil ports")
	}
}

func TestStreak_Observe_IncrementsOnPresence(t *testing.T) {
	s := NewStreak()
	_ = s.Observe([]int{80, 443})
	_ = s.Observe([]int{80, 443})
	if got := s.Current(80); got != 2 {
		t.Fatalf("expected streak 2, got %d", got)
	}
	if got := s.Current(443); got != 2 {
		t.Fatalf("expected streak 2, got %d", got)
	}
}

func TestStreak_Observe_ResetsOnAbsence(t *testing.T) {
	s := NewStreak()
	_ = s.Observe([]int{80})
	_ = s.Observe([]int{80})
	_ = s.Observe([]int{443}) // 80 absent
	if got := s.Current(80); got != 0 {
		t.Fatalf("expected streak 0 after absence, got %d", got)
	}
	if got := s.Current(443); got != 1 {
		t.Fatalf("expected streak 1 for 443, got %d", got)
	}
}

func TestStreak_Peak_TracksHighest(t *testing.T) {
	s := NewStreak()
	_ = s.Observe([]int{80})
	_ = s.Observe([]int{80})
	_ = s.Observe([]int{80})
	_ = s.Observe([]int{}) // resets
	_ = s.Observe([]int{80})
	if got := s.Current(80); got != 1 {
		t.Fatalf("expected current streak 1, got %d", got)
	}
	if got := s.Peak(80); got != 3 {
		t.Fatalf("expected peak 3, got %d", got)
	}
}

func TestStreak_Reset_ClearsAll(t *testing.T) {
	s := NewStreak()
	_ = s.Observe([]int{80})
	_ = s.Observe([]int{80})
	s.Reset()
	if got := s.Current(80); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if got := s.Peak(80); got != 0 {
		t.Fatalf("expected peak 0 after reset, got %d", got)
	}
}

func TestStreak_Observe_EmptyPorts_ResetsAll(t *testing.T) {
	s := NewStreak()
	_ = s.Observe([]int{80, 443})
	_ = s.Observe([]int{})
	if got := s.Current(80); got != 0 {
		t.Fatalf("expected 0 after empty scan, got %d", got)
	}
	if got := s.Current(443); got != 0 {
		t.Fatalf("expected 0 after empty scan, got %d", got)
	}
}
