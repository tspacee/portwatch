package watch

import (
	"errors"
	"testing"
	"time"
)

func TestNewRetryPolicy_InvalidMaxAttempts(t *testing.T) {
	_, err := NewRetryPolicy(0, time.Millisecond, 1)
	if err == nil {
		t.Fatal("expected error for zero maxAttempts")
	}
}

func TestNewRetryPolicy_InvalidMultiplier(t *testing.T) {
	_, err := NewRetryPolicy(3, time.Millisecond, 0.5)
	if err == nil {
		t.Fatal("expected error for multiplier < 1")
	}
}

func TestNewRetryPolicy_Valid(t *testing.T) {
	p, err := NewRetryPolicy(3, 10*time.Millisecond, 2.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.MaxAttempts != 3 {
		t.Errorf("expected MaxAttempts=3, got %d", p.MaxAttempts)
	}
}

func TestRetry_SucceedsOnFirstAttempt(t *testing.T) {
	p, _ := NewRetryPolicy(3, 0, 1)
	calls := 0
	err := p.Retry(func() error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Errorf("expected 1 call, got %d", calls)
	}
}

func TestRetry_SucceedsOnSecondAttempt(t *testing.T) {
	p, _ := NewRetryPolicy(3, 0, 1)
	calls := 0
	err := p.Retry(func() error {
		calls++
		if calls < 2 {
			return errors.New("fail")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 2 {
		t.Errorf("expected 2 calls, got %d", calls)
	}
}

func TestRetry_ExhaustsAttempts(t *testing.T) {
	p, _ := NewRetryPolicy(3, 0, 1)
	calls := 0
	err := p.Retry(func() error {
		calls++
		return errors.New("always fail")
	})
	if !errors.Is(err, ErrMaxRetries) {
		t.Fatalf("expected ErrMaxRetries, got %v", err)
	}
	if calls != 3 {
		t.Errorf("expected 3 calls, got %d", calls)
	}
}

func TestRetry_ZeroMaxAttempts(t *testing.T) {
	p := RetryPolicy{MaxAttempts: 0}
	err := p.Retry(func() error { return nil })
	if !errors.Is(err, ErrMaxRetries) {
		t.Fatalf("expected ErrMaxRetries, got %v", err)
	}
}
