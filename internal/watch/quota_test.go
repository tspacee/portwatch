package watch

import (
	"testing"
	"time"
)

func TestNewQuota_InvalidWindow(t *testing.T) {
	_, err := NewQuota(0, 5)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewQuota_InvalidMax(t *testing.T) {
	_, err := NewQuota(time.Minute, 0)
	if err == nil {
		t.Fatal("expected error for zero max")
	}
}

func TestNewQuota_Valid(t *testing.T) {
	q, err := NewQuota(time.Minute, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if q == nil {
		t.Fatal("expected non-nil quota")
	}
}

func TestQuota_Allow_WithinLimit(t *testing.T) {
	q, _ := NewQuota(time.Minute, 3)
	for i := 0; i < 3; i++ {
		if !q.Allow() {
			t.Fatalf("expected Allow()=true on call %d", i+1)
		}
	}
}

func TestQuota_Allow_ExceedsLimit(t *testing.T) {
	q, _ := NewQuota(time.Minute, 2)
	q.Allow()
	q.Allow()
	if q.Allow() {
		t.Fatal("expected Allow()=false after limit reached")
	}
}

func TestQuota_Remaining_DecreasesOnAllow(t *testing.T) {
	q, _ := NewQuota(time.Minute, 5)
	if q.Remaining() != 5 {
		t.Fatalf("expected 5 remaining, got %d", q.Remaining())
	}
	q.Allow()
	if q.Remaining() != 4 {
		t.Fatalf("expected 4 remaining, got %d", q.Remaining())
	}
}

func TestQuota_Reset_ClearsState(t *testing.T) {
	q, _ := NewQuota(time.Minute, 2)
	q.Allow()
	q.Allow()
	if q.Allow() {
		t.Fatal("expected blocked before reset")
	}
	q.Reset()
	if !q.Allow() {
		t.Fatal("expected Allow()=true after reset")
	}
}

func TestQuota_Evicts_ExpiredTimestamps(t *testing.T) {
	q, _ := NewQuota(50*time.Millisecond, 2)
	q.Allow()
	q.Allow()
	time.Sleep(60 * time.Millisecond)
	if !q.Allow() {
		t.Fatal("expected Allow()=true after window expired")
	}
}
