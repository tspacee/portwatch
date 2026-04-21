package watch

import (
	"testing"
	"time"
)

func TestNewDecay_InvalidRate_Zero(t *testing.T) {
	_, err := NewDecay(0)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestNewDecay_InvalidRate_Negative(t *testing.T) {
	_, err := NewDecay(-0.5)
	if err == nil {
		t.Fatal("expected error for negative rate")
	}
}

func TestNewDecay_InvalidRate_ExceedsOne(t *testing.T) {
	_, err := NewDecay(1.1)
	if err == nil {
		t.Fatal("expected error for rate > 1")
	}
}

func TestNewDecay_Valid(t *testing.T) {
	d, err := NewDecay(0.5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d == nil {
		t.Fatal("expected non-nil Decay")
	}
}

func TestDecay_Add_InvalidPort(t *testing.T) {
	d, _ := NewDecay(0.5)
	if err := d.Add(0, 1.0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := d.Add(70000, 1.0); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestDecay_Add_And_Score(t *testing.T) {
	d, _ := NewDecay(0.1)
	if err := d.Add(8080, 5.0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := d.Score(8080); got != 5.0 {
		t.Fatalf("expected 5.0, got %f", got)
	}
}

func TestDecay_Score_Missing_ReturnsZero(t *testing.T) {
	d, _ := NewDecay(0.5)
	if got := d.Score(9999); got != 0 {
		t.Fatalf("expected 0, got %f", got)
	}
}

func TestDecay_Tick_ReducesScore(t *testing.T) {
	d, _ := NewDecay(0.5)
	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	d.now = func() time.Time { return fixed }
	_ = d.Add(8080, 10.0)
	d.Tick() // sets last
	d.now = func() time.Time { return fixed.Add(1 * time.Second) }
	d.Tick() // applies decay: factor = 1 - 0.5*1 = 0.5 => 5.0
	got := d.Score(8080)
	if got < 4.9 || got > 5.1 {
		t.Fatalf("expected ~5.0 after decay, got %f", got)
	}
}

func TestDecay_Tick_PrunesSmallScores(t *testing.T) {
	d, _ := NewDecay(1.0)
	fixed := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	d.now = func() time.Time { return fixed }
	_ = d.Add(443, 0.5)
	d.Tick()
	d.now = func() time.Time { return fixed.Add(2 * time.Second) }
	d.Tick() // factor = 1 - 1.0*2 = clamped to 0 => pruned
	snap := d.Snapshot()
	if _, ok := snap[443]; ok {
		t.Fatal("expected port 443 to be pruned")
	}
}

func TestDecay_Snapshot_ReturnsCopy(t *testing.T) {
	d, _ := NewDecay(0.1)
	_ = d.Add(80, 3.0)
	snap := d.Snapshot()
	snap[80] = 999
	if d.Score(80) == 999 {
		t.Fatal("snapshot mutation affected internal state")
	}
}
