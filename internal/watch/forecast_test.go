package watch

import (
	"testing"
)

func TestNewForecast_InvalidWindow(t *testing.T) {
	_, err := NewForecast(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewForecast_Valid(t *testing.T) {
	f, err := NewForecast(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f == nil {
		t.Fatal("expected non-nil Forecast")
	}
}

func TestForecast_Likelihood_NeverObserved(t *testing.T) {
	f, _ := NewForecast(5)
	if got := f.Likelihood(80); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestForecast_Observe_InvalidPort(t *testing.T) {
	f, _ := NewForecast(5)
	if err := f.Observe(0, true); err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestForecast_Likelihood_AllPresent(t *testing.T) {
	f, _ := NewForecast(4)
	for i := 0; i < 4; i++ {
		_ = f.Observe(443, true)
	}
	if got := f.Likelihood(443); got != 1.0 {
		t.Fatalf("expected 1.0, got %f", got)
	}
}

func TestForecast_Likelihood_HalfPresent(t *testing.T) {
	f, _ := NewForecast(4)
	_ = f.Observe(443, true)
	_ = f.Observe(443, false)
	_ = f.Observe(443, true)
	_ = f.Observe(443, false)
	if got := f.Likelihood(443); got != 0.5 {
		t.Fatalf("expected 0.5, got %f", got)
	}
}

func TestForecast_Window_Evicts(t *testing.T) {
	f, _ := NewForecast(3)
	_ = f.Observe(22, true)
	_ = f.Observe(22, true)
	_ = f.Observe(22, true)
	_ = f.Observe(22, false) // evicts oldest true
	// window: [true, true, false] => 2/3
	got := f.Likelihood(22)
	if got < 0.66 || got > 0.68 {
		t.Fatalf("expected ~0.667, got %f", got)
	}
}

func TestForecast_Predict_AboveThreshold(t *testing.T) {
	f, _ := NewForecast(2)
	_ = f.Observe(8080, true)
	_ = f.Observe(8080, true)
	if !f.Predict(8080, 0.9) {
		t.Fatal("expected predict=true")
	}
}

func TestForecast_Reset_ClearsHistory(t *testing.T) {
	f, _ := NewForecast(5)
	_ = f.Observe(3306, true)
	f.Reset(3306)
	if got := f.Likelihood(3306); got != 0 {
		t.Fatalf("expected 0 after reset, got %f", got)
	}
}
