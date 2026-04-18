package watch

import (
	"testing"
	"time"
)

func TestNewBudget_InvalidWindow(t *testing.T) {
	_, err := NewBudget(0, time.Second)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewBudget_InvalidMaxUsage(t *testing.T) {
	_, err := NewBudget(time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero maxUsage")
	}
}

func TestNewBudget_Valid(t *testing.T) {
	b, err := NewBudget(time.Minute, time.Second*10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b == nil {
		t.Fatal("expected non-nil Budget")
	}
}

func TestBudget_Allow_WhenEmpty(t *testing.T) {
	b, _ := NewBudget(time.Minute, time.Second*5)
	if !b.Allow() {
		t.Fatal("expected Allow() true when no usage recorded")
	}
}

func TestBudget_Allow_WhenExhausted(t *testing.T) {
	b, _ := NewBudget(time.Minute, time.Second*2)
	b.Record(time.Second * 3)
	if b.Allow() {
		t.Fatal("expected Allow() false when budget exhausted")
	}
}

func TestBudget_Allow_PermitsAfterWindowExpiry(t *testing.T) {
	b, _ := NewBudget(50*time.Millisecond, time.Second)
	b.Record(time.Second * 2)
	time.Sleep(60 * time.Millisecond)
	if !b.Allow() {
		t.Fatal("expected Allow() true after window expired")
	}
}

func TestBudget_Used_ReflectsRecorded(t *testing.T) {
	b, _ := NewBudget(time.Minute, time.Second*10)
	b.Record(200 * time.Millisecond)
	b.Record(300 * time.Millisecond)
	if b.Used() != 500*time.Millisecond {
		t.Fatalf("expected 500ms used, got %v", b.Used())
	}
}

func TestBudget_Used_EvictsExpired(t *testing.T) {
	b, _ := NewBudget(50*time.Millisecond, time.Second*10)
	b.Record(500 * time.Millisecond)
	time.Sleep(60 * time.Millisecond)
	if b.Used() != 0 {
		t.Fatalf("expected 0 after eviction, got %v", b.Used())
	}
}
