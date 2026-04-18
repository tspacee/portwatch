package watch

import (
	"testing"
	"time"
)

func TestNewSampler_InvalidWindow(t *testing.T) {
	_, err := NewSampler(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewSampler_Valid(t *testing.T) {
	s, err := NewSampler(10 * time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil sampler")
	}
}

func TestSampler_Count_Empty(t *testing.T) {
	s, _ := NewSampler(10 * time.Second)
	if s.Count() != 0 {
		t.Fatalf("expected 0, got %d", s.Count())
	}
}

func TestSampler_Record_IncrementsCount(t *testing.T) {
	s, _ := NewSampler(10 * time.Second)
	s.Record()
	s.Record()
	if s.Count() != 2 {
		t.Fatalf("expected 2, got %d", s.Count())
	}
}

func TestSampler_Rate_NonZero(t *testing.T) {
	s, _ := NewSampler(10 * time.Second)
	s.Record()
	s.Record()
	if s.Rate() <= 0 {
		t.Fatal("expected positive rate")
	}
}

func TestSampler_Rate_Empty(t *testing.T) {
	s, _ := NewSampler(10 * time.Second)
	if s.Rate() != 0 {
		t.Fatal("expected zero rate for empty sampler")
	}
}

func TestSampler_Reset_ClearsCount(t *testing.T) {
	s, _ := NewSampler(10 * time.Second)
	s.Record()
	s.Record()
	s.Reset()
	if s.Count() != 0 {
		t.Fatalf("expected 0 after reset, got %d", s.Count())
	}
}

func TestSampler_Evicts_OldSamples(t *testing.T) {
	s, _ := NewSampler(50 * time.Millisecond)
	s.Record()
	time.Sleep(60 * time.Millisecond)
	if s.Count() != 0 {
		t.Fatalf("expected 0 after expiry, got %d", s.Count())
	}
}
