package watch

import (
	"testing"
	"time"
)

func TestNewVelocity_InvalidWindow(t *testing.T) {
	_, err := NewVelocity(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewVelocity_Valid(t *testing.T) {
	v, err := NewVelocity(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v == nil {
		t.Fatal("expected non-nil Velocity")
	}
}

func TestVelocity_Record_InvalidPort(t *testing.T) {
	v, _ := NewVelocity(time.Second)
	if err := v.Record(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := v.Record(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestVelocity_Rate_Empty(t *testing.T) {
	v, _ := NewVelocity(time.Second)
	if r := v.Rate(80); r != 0 {
		t.Fatalf("expected 0, got %d", r)
	}
}

func TestVelocity_Record_IncrementsRate(t *testing.T) {
	v, _ := NewVelocity(time.Second)
	_ = v.Record(80)
	_ = v.Record(80)
	if r := v.Rate(80); r != 2 {
		t.Fatalf("expected 2, got %d", r)
	}
}

func TestVelocity_Rate_EvictsExpired(t *testing.T) {
	v, _ := NewVelocity(50 * time.Millisecond)
	_ = v.Record(443)
	time.Sleep(80 * time.Millisecond)
	if r := v.Rate(443); r != 0 {
		t.Fatalf("expected 0 after expiry, got %d", r)
	}
}

func TestVelocity_Reset_ClearsEvents(t *testing.T) {
	v, _ := NewVelocity(time.Second)
	_ = v.Record(8080)
	_ = v.Record(8080)
	v.Reset(8080)
	if r := v.Rate(8080); r != 0 {
		t.Fatalf("expected 0 after reset, got %d", r)
	}
}

func TestVelocity_MultiplePortsIndependent(t *testing.T) {
	v, _ := NewVelocity(time.Second)
	_ = v.Record(22)
	_ = v.Record(22)
	_ = v.Record(80)
	if r := v.Rate(22); r != 2 {
		t.Fatalf("expected 2 for port 22, got %d", r)
	}
	if r := v.Rate(80); r != 1 {
		t.Fatalf("expected 1 for port 80, got %d", r)
	}
}
