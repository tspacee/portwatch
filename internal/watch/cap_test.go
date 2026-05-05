package watch

import (
	"testing"
)

func TestNewCap_InvalidMax(t *testing.T) {
	_, err := NewCap(0)
	if err == nil {
		t.Fatal("expected error for max=0")
	}
}

func TestNewCap_Valid(t *testing.T) {
	c, err := NewCap(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Len() != 0 {
		t.Fatalf("expected 0 active, got %d", c.Len())
	}
}

func TestCap_Acquire_Valid(t *testing.T) {
	c, _ := NewCap(2)
	ok, err := c.Acquire(80)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !ok {
		t.Fatal("expected acquire to succeed")
	}
	if c.Len() != 1 {
		t.Fatalf("expected len=1, got %d", c.Len())
	}
}

func TestCap_Acquire_InvalidPort(t *testing.T) {
	c, _ := NewCap(2)
	_, err := c.Acquire(0)
	if err == nil {
		t.Fatal("expected error for port=0")
	}
}

func TestCap_Acquire_ExceedsCap(t *testing.T) {
	c, _ := NewCap(1)
	c.Acquire(80)
	ok, err := c.Acquire(443)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ok {
		t.Fatal("expected acquire to fail when cap reached")
	}
}

func TestCap_Acquire_Idempotent(t *testing.T) {
	c, _ := NewCap(1)
	c.Acquire(80)
	ok, _ := c.Acquire(80)
	if !ok {
		t.Fatal("expected idempotent acquire to succeed")
	}
	if c.Len() != 1 {
		t.Fatalf("expected len=1 after idempotent acquire, got %d", c.Len())
	}
}

func TestCap_Release_RemovesPort(t *testing.T) {
	c, _ := NewCap(2)
	c.Acquire(80)
	c.Release(80)
	if c.Len() != 0 {
		t.Fatalf("expected len=0 after release, got %d", c.Len())
	}
}

func TestCap_Full_ReflectsState(t *testing.T) {
	c, _ := NewCap(1)
	if c.Full() {
		t.Fatal("expected not full initially")
	}
	c.Acquire(80)
	if !c.Full() {
		t.Fatal("expected full after acquiring up to max")
	}
}
