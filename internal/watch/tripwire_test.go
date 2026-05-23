package watch

import (
	"sync/atomic"
	"testing"
)

func TestNewTripwire_EmptyPorts(t *testing.T) {
	_, err := NewTripwire(nil, func(int) {})
	if err == nil {
		t.Fatal("expected error for empty ports")
	}
}

func TestNewTripwire_NilCallback(t *testing.T) {
	_, err := NewTripwire([]int{80}, nil)
	if err == nil {
		t.Fatal("expected error for nil callback")
	}
}

func TestNewTripwire_InvalidPort(t *testing.T) {
	_, err := NewTripwire([]int{0}, func(int) {})
	if err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestNewTripwire_Valid(t *testing.T) {
	tw, err := NewTripwire([]int{80, 443}, func(int) {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tw.Tripped() {
		t.Fatal("expected not tripped initially")
	}
}

func TestTripwire_Observe_NoTrip(t *testing.T) {
	var fired int32
	tw, _ := NewTripwire([]int{80}, func(int) { atomic.AddInt32(&fired, 1) })
	tw.Observe([]int{80, 443})
	if atomic.LoadInt32(&fired) != 0 {
		t.Fatal("expected no trip when watched port is present")
	}
	if tw.Tripped() {
		t.Fatal("expected not tripped")
	}
}

func TestTripwire_Observe_TripsOnAbsence(t *testing.T) {
	var trippedPort int
	tw, _ := NewTripwire([]int{8080}, func(p int) { trippedPort = p })
	tw.Observe([]int{80})
	if !tw.Tripped() {
		t.Fatal("expected tripped")
	}
	if trippedPort != 8080 {
		t.Fatalf("expected tripped port 8080, got %d", trippedPort)
	}
}

func TestTripwire_Observe_FiresOnlyOnce(t *testing.T) {
	var count int32
	tw, _ := NewTripwire([]int{8080}, func(int) { atomic.AddInt32(&count, 1) })
	tw.Observe([]int{})
	tw.Observe([]int{})
	if atomic.LoadInt32(&count) != 1 {
		t.Fatalf("expected callback once, got %d", atomic.LoadInt32(&count))
	}
}

func TestTripwire_Reset_AllowsRetrigger(t *testing.T) {
	var count int32
	tw, _ := NewTripwire([]int{22}, func(int) { atomic.AddInt32(&count, 1) })
	tw.Observe([]int{})
	tw.Reset()
	if tw.Tripped() {
		t.Fatal("expected not tripped after reset")
	}
	tw.Observe([]int{})
	if atomic.LoadInt32(&count) != 2 {
		t.Fatalf("expected 2 callbacks after reset, got %d", atomic.LoadInt32(&count))
	}
}
