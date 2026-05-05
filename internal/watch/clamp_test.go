package watch

import (
	"testing"
)

func TestNewClamp_InvalidMin(t *testing.T) {
	_, err := NewClamp(-1, 10)
	if err == nil {
		t.Fatal("expected error for negative min")
	}
}

func TestNewClamp_MaxLessThanMin(t *testing.T) {
	_, err := NewClamp(5, 3)
	if err == nil {
		t.Fatal("expected error when max < min")
	}
}

func TestNewClamp_Valid(t *testing.T) {
	c, err := NewClamp(1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil Clamp")
	}
}

func TestClamp_Set_InvalidPort(t *testing.T) {
	c, _ := NewClamp(0, 100)
	if err := c.Set(0, 5); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := c.Set(65536, 5); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestClamp_Set_ClampsToMin(t *testing.T) {
	c, _ := NewClamp(3, 10)
	_ = c.Set(80, 1)
	v, ok := c.Get(80)
	if !ok {
		t.Fatal("expected port to be present")
	}
	if v != 3 {
		t.Fatalf("expected clamped value 3, got %d", v)
	}
}

func TestClamp_Set_ClampsToMax(t *testing.T) {
	c, _ := NewClamp(0, 5)
	_ = c.Set(443, 99)
	v, _ := c.Get(443)
	if v != 5 {
		t.Fatalf("expected clamped value 5, got %d", v)
	}
}

func TestClamp_Set_WithinRange(t *testing.T) {
	c, _ := NewClamp(1, 10)
	_ = c.Set(8080, 7)
	v, _ := c.Get(8080)
	if v != 7 {
		t.Fatalf("expected 7, got %d", v)
	}
}

func TestClamp_Get_Missing(t *testing.T) {
	c, _ := NewClamp(0, 10)
	v, ok := c.Get(9999)
	if ok {
		t.Fatal("expected missing port to return false")
	}
	if v != 0 {
		t.Fatalf("expected 0 for missing port, got %d", v)
	}
}

func TestClamp_Reset_ClearsAll(t *testing.T) {
	c, _ := NewClamp(0, 100)
	_ = c.Set(22, 10)
	_ = c.Set(80, 20)
	c.Reset()
	if c.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", c.Len())
	}
}

func TestClamp_Len_Tracks(t *testing.T) {
	c, _ := NewClamp(0, 50)
	_ = c.Set(22, 1)
	_ = c.Set(80, 2)
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
}
