package watch

import (
	"testing"
	"time"
)

func TestNewCanary_Empty(t *testing.T) {
	c := NewCanary()
	if c.Len() != 0 {
		t.Fatalf("expected 0 watched ports, got %d", c.Len())
	}
}

func TestCanary_Watch_Valid(t *testing.T) {
	c := NewCanary()
	if err := c.Watch(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 1 {
		t.Fatalf("expected 1 watched port, got %d", c.Len())
	}
}

func TestCanary_Watch_InvalidPort(t *testing.T) {
	c := NewCanary()
	if err := c.Watch(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := c.Watch(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestCanary_Observe_NoTrip(t *testing.T) {
	c := NewCanary()
	_ = c.Watch(80)
	_ = c.Watch(443)
	c.Observe([]int{80, 443, 8080})
	if len(c.Tripped()) != 0 {
		t.Fatal("expected no tripped ports")
	}
}

func TestCanary_Observe_TripsOnAbsence(t *testing.T) {
	c := NewCanary()
	_ = c.Watch(80)
	_ = c.Watch(443)
	c.Observe([]int{80}) // 443 is absent
	tripped := c.Tripped()
	if _, ok := tripped[443]; !ok {
		t.Fatal("expected port 443 to be tripped")
	}
	if _, ok := tripped[80]; ok {
		t.Fatal("port 80 should not be tripped")
	}
}

func TestCanary_Observe_RecordsTime(t *testing.T) {
	c := NewCanary()
	_ = c.Watch(9000)
	before := time.Now()
	c.Observe([]int{})
	after := time.Now()
	tripped := c.Tripped()
	ts, ok := tripped[9000]
	if !ok {
		t.Fatal("expected port 9000 to be tripped")
	}
	if ts.Before(before) || ts.After(after) {
		t.Fatalf("trip time %v not within expected range", ts)
	}
}

func TestCanary_Observe_ClearsOnReturn(t *testing.T) {
	c := NewCanary()
	_ = c.Watch(22)
	c.Observe([]int{})   // trip port 22
	c.Observe([]int{22}) // port 22 comes back
	if len(c.Tripped()) != 0 {
		t.Fatal("expected tripped map to be cleared after port returned")
	}
}

func TestCanary_Tripped_ReturnsCopy(t *testing.T) {
	c := NewCanary()
	_ = c.Watch(3306)
	c.Observe([]int{})
	tripped := c.Tripped()
	delete(tripped, 3306)
	if len(c.Tripped()) != 1 {
		t.Fatal("Tripped should return a copy, not a reference")
	}
}
