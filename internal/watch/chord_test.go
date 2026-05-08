package watch

import (
	"testing"
)

func TestNewChord_EmptyRequired(t *testing.T) {
	_, err := NewChord([]int{}, func([]int) {})
	if err == nil {
		t.Fatal("expected error for empty required ports")
	}
}

func TestNewChord_NilCallback(t *testing.T) {
	_, err := NewChord([]int{80, 443}, nil)
	if err == nil {
		t.Fatal("expected error for nil callback")
	}
}

func TestNewChord_InvalidPort(t *testing.T) {
	_, err := NewChord([]int{80, 99999}, func([]int) {})
	if err == nil {
		t.Fatal("expected error for out-of-range port")
	}
}

func TestNewChord_Valid(t *testing.T) {
	c, err := NewChord([]int{80, 443}, func([]int) {})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 2 {
		t.Fatalf("expected len 2, got %d", c.Len())
	}
}

func TestChord_Observe_AllPresent_FiresCallback(t *testing.T) {
	fired := false
	c, _ := NewChord([]int{80, 443}, func(ports []int) {
		fired = true
		if len(ports) != 2 {
			t.Errorf("expected 2 matched ports, got %d", len(ports))
		}
	})
	c.Observe([]int{22, 80, 443, 8080})
	if !fired {
		t.Fatal("expected callback to fire")
	}
}

func TestChord_Observe_PartialPresent_NoCallback(t *testing.T) {
	fired := false
	c, _ := NewChord([]int{80, 443}, func([]int) { fired = true })
	c.Observe([]int{80})
	if fired {
		t.Fatal("callback should not fire when not all ports are present")
	}
}

func TestChord_Observe_NoPorts_NoCallback(t *testing.T) {
	fired := false
	c, _ := NewChord([]int{80, 443}, func([]int) { fired = true })
	c.Observe([]int{})
	if fired {
		t.Fatal("callback should not fire on empty observation")
	}
}

func TestChord_Observe_ExactMatch_FiresCallback(t *testing.T) {
	fired := false
	c, _ := NewChord([]int{8080}, func([]int) { fired = true })
	c.Observe([]int{8080})
	if !fired {
		t.Fatal("expected callback to fire on exact match")
	}
}
