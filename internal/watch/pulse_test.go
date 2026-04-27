package watch

import (
	"testing"
	"time"
)

func TestNewPulse_InitialState(t *testing.T) {
	p := NewPulse()
	_, err := p.Summary()
	if err == nil {
		t.Fatal("expected error for empty pulse, got nil")
	}
}

func TestPulse_Beat_FirstCallNoInterval(t *testing.T) {
	p := NewPulse()
	p.Beat()
	_, err := p.Summary()
	if err == nil {
		t.Fatal("expected error after single beat, got nil")
	}
}

func TestPulse_Beat_RecordsInterval(t *testing.T) {
	p := NewPulse()
	p.Beat()
	time.Sleep(20 * time.Millisecond)
	p.Beat()

	s, err := p.Summary()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count != 1 {
		t.Errorf("expected count 1, got %d", s.Count)
	}
	if s.Avg <= 0 {
		t.Errorf("expected positive avg, got %v", s.Avg)
	}
}

func TestPulse_Summary_MinMaxAvg(t *testing.T) {
	p := NewPulse()
	p.Beat()
	time.Sleep(10 * time.Millisecond)
	p.Beat()
	time.Sleep(30 * time.Millisecond)
	p.Beat()

	s, err := p.Summary()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.Count != 2 {
		t.Errorf("expected count 2, got %d", s.Count)
	}
	if s.Min > s.Max {
		t.Errorf("min %v > max %v", s.Min, s.Max)
	}
	if s.Avg < s.Min || s.Avg > s.Max {
		t.Errorf("avg %v out of [%v, %v]", s.Avg, s.Min, s.Max)
	}
}

func TestPulse_Reset_ClearsState(t *testing.T) {
	p := NewPulse()
	p.Beat()
	time.Sleep(10 * time.Millisecond)
	p.Beat()

	p.Reset()
	_, err := p.Summary()
	if err == nil {
		t.Fatal("expected error after reset, got nil")
	}
}

func TestPulse_Reset_AllowsReuse(t *testing.T) {
	p := NewPulse()
	p.Beat()
	time.Sleep(10 * time.Millisecond)
	p.Beat()
	p.Reset()

	p.Beat()
	time.Sleep(10 * time.Millisecond)
	p.Beat()

	s, err := p.Summary()
	if err != nil {
		t.Fatalf("unexpected error after reuse: %v", err)
	}
	if s.Count != 1 {
		t.Errorf("expected count 1 after reuse, got %d", s.Count)
	}
}
