package watch

import (
	"testing"
)

func TestNewAnomaly_InvalidTolerance(t *testing.T) {
	_, err := NewAnomaly(0)
	if err == nil {
		t.Fatal("expected error for zero tolerance")
	}
	_, err = NewAnomaly(-0.1)
	if err == nil {
		t.Fatal("expected error for negative tolerance")
	}
}

func TestNewAnomaly_Valid(t *testing.T) {
	a, err := NewAnomaly(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a == nil {
		t.Fatal("expected non-nil Anomaly")
	}
}

func TestAnomaly_SetBaseline_InvalidPort(t *testing.T) {
	a, _ := NewAnomaly(0.5)
	if err := a.SetBaseline(0, 1.0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := a.SetBaseline(65536, 1.0); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestAnomaly_SetBaseline_NegativeFreq(t *testing.T) {
	a, _ := NewAnomaly(0.5)
	if err := a.SetBaseline(80, -1.0); err == nil {
		t.Fatal("expected error for negative frequency")
	}
}

func TestAnomaly_Observe_WithinTolerance(t *testing.T) {
	a, _ := NewAnomaly(0.5)
	_ = a.SetBaseline(80, 10.0)

	// 14.0 is 40% above baseline — within 50% tolerance
	anom, err := a.Observe(80, 14.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if anom {
		t.Error("expected port 80 to be within tolerance")
	}
}

func TestAnomaly_Observe_ExceedsTolerance(t *testing.T) {
	a, _ := NewAnomaly(0.5)
	_ = a.SetBaseline(80, 10.0)

	// 16.0 is 60% above baseline — exceeds 50% tolerance
	anom, err := a.Observe(80, 16.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !anom {
		t.Error("expected port 80 to be anomalous")
	}
}

func TestAnomaly_Observe_NoBaseline_IsAnomalous(t *testing.T) {
	a, _ := NewAnomaly(0.5)
	anom, err := a.Observe(443, 5.0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !anom {
		t.Error("expected port with no baseline to be anomalous")
	}
}

func TestAnomaly_IsAnomalous_ReflectsLastObservation(t *testing.T) {
	a, _ := NewAnomaly(0.2)
	_ = a.SetBaseline(22, 5.0)
	_, _ = a.Observe(22, 10.0) // 100% deviation — anomalous
	if !a.IsAnomalous(22) {
		t.Error("expected port 22 to be marked anomalous")
	}
}

func TestAnomaly_Reset_ClearsState(t *testing.T) {
	a, _ := NewAnomaly(0.1)
	_, _ = a.Observe(8080, 99.0)
	if !a.IsAnomalous(8080) {
		t.Fatal("expected anomaly before reset")
	}
	a.Reset()
	if a.IsAnomalous(8080) {
		t.Error("expected anomaly cleared after reset")
	}
}
