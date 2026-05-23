package watch

import (
	"testing"
)

func TestNewHeatMap_Empty(t *testing.T) {
	hm := NewHeatMap()
	if hm == nil {
		t.Fatal("expected non-nil HeatMap")
	}
	if got := hm.Heat(80); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestHeatMap_Record_Valid(t *testing.T) {
	hm := NewHeatMap()
	if err := hm.Record(443); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := hm.Heat(443); got != 1 {
		t.Fatalf("expected 1, got %d", got)
	}
}

func TestHeatMap_Record_InvalidPort(t *testing.T) {
	hm := NewHeatMap()
	if err := hm.Record(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := hm.Record(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestHeatMap_Record_Accumulates(t *testing.T) {
	hm := NewHeatMap()
	for i := 0; i < 5; i++ {
		_ = hm.Record(22)
	}
	if got := hm.Heat(22); got != 5 {
		t.Fatalf("expected 5, got %d", got)
	}
}

func TestHeatMap_Hottest_ReturnsHighest(t *testing.T) {
	hm := NewHeatMap()
	_ = hm.Record(80)
	_ = hm.Record(80)
	_ = hm.Record(443)
	if got := hm.Hottest(); got != 80 {
		t.Fatalf("expected 80, got %d", got)
	}
}

func TestHeatMap_Hottest_EmptyReturnsZero(t *testing.T) {
	hm := NewHeatMap()
	if got := hm.Hottest(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestHeatMap_Snapshot_ReturnsCopy(t *testing.T) {
	hm := NewHeatMap()
	_ = hm.Record(8080)
	snap := hm.Snapshot()
	snap[8080] = 999
	if got := hm.Heat(8080); got != 1 {
		t.Fatalf("snapshot mutation affected original: got %d", got)
	}
}

func TestHeatMap_Reset_ClearsAll(t *testing.T) {
	hm := NewHeatMap()
	_ = hm.Record(22)
	_ = hm.Record(80)
	hm.Reset()
	if got := hm.Heat(22); got != 0 {
		t.Fatalf("expected 0 after reset, got %d", got)
	}
	if got := hm.Hottest(); got != 0 {
		t.Fatalf("expected hottest 0 after reset, got %d", got)
	}
}
