package watch

import (
	"testing"
)

func TestNewWatermark_Empty(t *testing.T) {
	wm := NewWatermark()
	if wm == nil {
		t.Fatal("expected non-nil Watermark")
	}
	if len(wm.Snapshot()) != 0 {
		t.Error("expected empty snapshot on creation")
	}
}

func TestWatermark_Record_Valid(t *testing.T) {
	wm := NewWatermark()
	if err := wm.Record(8080, 5); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	v, ok := wm.Peak(8080)
	if !ok {
		t.Fatal("expected entry for port 8080")
	}
	if v != 5 {
		t.Errorf("expected 5, got %d", v)
	}
}

func TestWatermark_Record_InvalidPort(t *testing.T) {
	wm := NewWatermark()
	if err := wm.Record(0, 1); err == nil {
		t.Error("expected error for port 0")
	}
	if err := wm.Record(65536, 1); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestWatermark_Record_KeepsHighest(t *testing.T) {
	wm := NewWatermark()
	_ = wm.Record(443, 3)
	_ = wm.Record(443, 7)
	_ = wm.Record(443, 2)
	v, _ := wm.Peak(443)
	if v != 7 {
		t.Errorf("expected high-water mark 7, got %d", v)
	}
}

func TestWatermark_Peak_Missing(t *testing.T) {
	wm := NewWatermark()
	_, ok := wm.Peak(9999)
	if ok {
		t.Error("expected false for unseen port")
	}
}

func TestWatermark_Reset_ClearsEntry(t *testing.T) {
	wm := NewWatermark()
	_ = wm.Record(22, 10)
	wm.Reset(22)
	_, ok := wm.Peak(22)
	if ok {
		t.Error("expected entry to be cleared after Reset")
	}
}

func TestWatermark_Snapshot_ReturnsCopy(t *testing.T) {
	wm := NewWatermark()
	_ = wm.Record(80, 4)
	_ = wm.Record(443, 9)
	snap := wm.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	// Mutating the copy must not affect the watermark.
	snap[80] = 999
	v, _ := wm.Peak(80)
	if v == 999 {
		t.Error("Snapshot returned a reference, not a copy")
	}
}
