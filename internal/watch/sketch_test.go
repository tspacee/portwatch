package watch

import (
	"testing"
)

func TestNewSketch_Empty(t *testing.T) {
	s := NewSketch()
	if s.Total() != 0 {
		t.Fatalf("expected total 0, got %d", s.Total())
	}
}

func TestSketch_Record_Valid(t *testing.T) {
	s := NewSketch()
	if err := s.Record(80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Estimate(80) != 1 {
		t.Fatalf("expected estimate 1, got %d", s.Estimate(80))
	}
	if s.Total() != 1 {
		t.Fatalf("expected total 1, got %d", s.Total())
	}
}

func TestSketch_Record_InvalidPort(t *testing.T) {
	s := NewSketch()
	if err := s.Record(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := s.Record(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestSketch_Record_Accumulates(t *testing.T) {
	s := NewSketch()
	for i := 0; i < 5; i++ {
		_ = s.Record(443)
	}
	if s.Estimate(443) != 5 {
		t.Fatalf("expected estimate 5, got %d", s.Estimate(443))
	}
	if s.Total() != 5 {
		t.Fatalf("expected total 5, got %d", s.Total())
	}
}

func TestSketch_Estimate_UnseenPort_ReturnsZero(t *testing.T) {
	s := NewSketch()
	if got := s.Estimate(8080); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestSketch_Reset_ClearsState(t *testing.T) {
	s := NewSketch()
	_ = s.Record(22)
	_ = s.Record(80)
	s.Reset()
	if s.Total() != 0 {
		t.Fatalf("expected total 0 after reset, got %d", s.Total())
	}
	if s.Estimate(22) != 0 {
		t.Fatalf("expected estimate 0 after reset")
	}
}

func TestSketch_TopN_ReturnsHighestPorts(t *testing.T) {
	s := NewSketch()
	_ = s.Record(80)
	_ = s.Record(80)
	_ = s.Record(80)
	_ = s.Record(443)
	_ = s.Record(443)
	_ = s.Record(22)

	top := s.TopN(2)
	if len(top) != 2 {
		t.Fatalf("expected 2 results, got %d", len(top))
	}
	if top[0] != 80 {
		t.Errorf("expected port 80 at index 0, got %d", top[0])
	}
	if top[1] != 443 {
		t.Errorf("expected port 443 at index 1, got %d", top[1])
	}
}

func TestSketch_TopN_FewerThanN(t *testing.T) {
	s := NewSketch()
	_ = s.Record(8080)
	top := s.TopN(5)
	if len(top) != 1 {
		t.Fatalf("expected 1 result, got %d", len(top))
	}
}
