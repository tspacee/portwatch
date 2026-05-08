package watch

import (
	"testing"
)

func TestNewSequence_Empty(t *testing.T) {
	s := NewSequence()
	if s == nil {
		t.Fatal("expected non-nil Sequence")
	}
	if got := s.Current(80); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSequence_Next_Valid(t *testing.T) {
	s := NewSequence()
	n, err := s.Next(80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Fatalf("expected 1, got %d", n)
	}
}

func TestSequence_Next_InvalidPort(t *testing.T) {
	s := NewSequence()
	_, err := s.Next(0)
	if err == nil {
		t.Fatal("expected error for port 0")
	}
	_, err = s.Next(65536)
	if err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestSequence_Next_Increments(t *testing.T) {
	s := NewSequence()
	for i := uint64(1); i <= 5; i++ {
		n, err := s.Next(443)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if n != i {
			t.Fatalf("expected %d, got %d", i, n)
		}
	}
}

func TestSequence_Reset_ClearsCounter(t *testing.T) {
	s := NewSequence()
	s.Next(8080) //nolint
	s.Next(8080) //nolint
	s.Reset(8080)
	if got := s.Current(8080); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
}

func TestSequence_Snapshot_ReturnsCopy(t *testing.T) {
	s := NewSequence()
	s.Next(22)  //nolint
	s.Next(22)  //nolint
	s.Next(443) //nolint
	snap := s.Snapshot()
	if snap[22] != 2 {
		t.Fatalf("expected 2 for port 22, got %d", snap[22])
	}
	if snap[443] != 1 {
		t.Fatalf("expected 1 for port 443, got %d", snap[443])
	}
	// mutating snapshot must not affect original
	snap[22] = 99
	if s.Current(22) != 2 {
		t.Fatal("snapshot mutation affected internal state")
	}
}

func TestSequence_IndependentPerPort(t *testing.T) {
	s := NewSequence()
	s.Next(80)  //nolint
	s.Next(80)  //nolint
	s.Next(443) //nolint
	if s.Current(80) != 2 {
		t.Fatalf("expected 2 for port 80")
	}
	if s.Current(443) != 1 {
		t.Fatalf("expected 1 for port 443")
	}
}
