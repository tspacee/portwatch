package watch

import (
	"testing"
)

func TestNewPivot_Empty(t *testing.T) {
	pv := NewPivot()
	if pv.Len() != 0 {
		t.Fatalf("expected len 0, got %d", pv.Len())
	}
}

func TestPivot_Observe_InvalidPort(t *testing.T) {
	pv := NewPivot()
	if err := pv.Observe(0, true, 1); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := pv.Observe(65536, true, 1); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestPivot_Observe_RecordsPivot(t *testing.T) {
	pv := NewPivot()
	if err := pv.Observe(80, true, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	idx, ok := pv.Since(80)
	if !ok {
		t.Fatal("expected port 80 to be tracked")
	}
	if idx != 5 {
		t.Fatalf("expected pivot index 5, got %d", idx)
	}
}

func TestPivot_Observe_SameState_DoesNotUpdateIndex(t *testing.T) {
	pv := NewPivot()
	_ = pv.Observe(443, true, 3)
	_ = pv.Observe(443, true, 9) // same state — should not update
	idx, _ := pv.Since(443)
	if idx != 3 {
		t.Fatalf("expected pivot index 3, got %d", idx)
	}
}

func TestPivot_Observe_StateChange_UpdatesIndex(t *testing.T) {
	pv := NewPivot()
	_ = pv.Observe(22, true, 1)
	_ = pv.Observe(22, false, 7) // state changed
	idx, _ := pv.Since(22)
	if idx != 7 {
		t.Fatalf("expected pivot index 7, got %d", idx)
	}
}

func TestPivot_Since_Missing(t *testing.T) {
	pv := NewPivot()
	_, ok := pv.Since(8080)
	if ok {
		t.Fatal("expected false for untracked port")
	}
}

func TestPivot_Reset_RemovesEntry(t *testing.T) {
	pv := NewPivot()
	_ = pv.Observe(3000, true, 2)
	pv.Reset(3000)
	if pv.Len() != 0 {
		t.Fatalf("expected len 0 after reset, got %d", pv.Len())
	}
	_, ok := pv.Since(3000)
	if ok {
		t.Fatal("expected port to be removed after reset")
	}
}

func TestPivot_Len_TracksMultiplePorts(t *testing.T) {
	pv := NewPivot()
	_ = pv.Observe(80, true, 1)
	_ = pv.Observe(443, true, 2)
	_ = pv.Observe(8080, false, 3)
	if pv.Len() != 3 {
		t.Fatalf("expected len 3, got %d", pv.Len())
	}
}
