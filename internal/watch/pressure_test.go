package watch

import (
	"testing"
	"time"
)

func TestNewPressure_InvalidWindow(t *testing.T) {
	_, err := NewPressure(0, 5)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewPressure_InvalidThreshold(t *testing.T) {
	_, err := NewPressure(time.Second, 0)
	if err == nil {
		t.Fatal("expected error for zero threshold")
	}
}

func TestNewPressure_Valid(t *testing.T) {
	p, err := NewPressure(time.Second, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Pressure")
	}
}

func TestPressure_Count_Empty(t *testing.T) {
	p, _ := NewPressure(time.Second, 3)
	if p.Count() != 0 {
		t.Errorf("expected 0, got %d", p.Count())
	}
}

func TestPressure_Record_IncrementsCount(t *testing.T) {
	p, _ := NewPressure(time.Second, 3)
	p.Record()
	p.Record()
	if{
		t.Errorf("expected 2, got %d", p.Count())
	}
}

func TestPressure_High_BelowThreshold(t *testing.T) {
	p, _ := NewPressure(time.Second, 5)
	p.Record()
	p.Record()
	if p.High() {
		t.Error("expected High() to be false below threshold")
	}
}

func TestPressure_High_AtThreshold(t *testing.T) {
	p, _ := NewPressure(time.Second, 3)
	p.Record()
	p.Record()
	p.Record()
	if !p.High() {
		t.Error("expected High() to be true at threshold")
	}
}

func TestPressure_Evicts_ExpiredEvents(t *testing.T) {
	p, _ := NewPressure(50*time.Millisecond, 3)
	p.Record()
	p.Record()
	time.Sleep(60 * time.Millisecond)
	p.Record()
	if p.Count() != 1 {
		t.Errorf("expected 1 after eviction, got %d", p.Count())
	}
}
