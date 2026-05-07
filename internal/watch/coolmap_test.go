package watch

import (
	"testing"
	"time"
)

func TestNewCoolMap_InvalidWindow(t *testing.T) {
	_, err := NewCoolMap(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewCoolMap_Valid(t *testing.T) {
	cm, err := NewCoolMap(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cm == nil {
		t.Fatal("expected non-nil CoolMap")
	}
}

func TestCoolMap_Ready_FirstCall_ReturnsTrue(t *testing.T) {
	cm, _ := NewCoolMap(100 * time.Millisecond)
	if !cm.Ready(8080) {
		t.Fatal("expected Ready to return true on first call")
	}
}

func TestCoolMap_Ready_SecondCall_ReturnsFalse(t *testing.T) {
	cm, _ := NewCoolMap(200 * time.Millisecond)
	cm.Ready(9090)
	if cm.Ready(9090) {
		t.Fatal("expected Ready to return false within cooldown window")
	}
}

func TestCoolMap_Ready_AfterExpiry_ReturnsTrue(t *testing.T) {
	cm, _ := NewCoolMap(20 * time.Millisecond)
	cm.Ready(7070)
	time.Sleep(40 * time.Millisecond)
	if !cm.Ready(7070) {
		t.Fatal("expected Ready to return true after window expires")
	}
}

func TestCoolMap_Reset_ClearsCooldown(t *testing.T) {
	cm, _ := NewCoolMap(500 * time.Millisecond)
	cm.Ready(1234)
	cm.Reset(1234)
	if !cm.Ready(1234) {
		t.Fatal("expected Ready to return true after Reset")
	}
}

func TestCoolMap_Len_TracksEntries(t *testing.T) {
	cm, _ := NewCoolMap(100 * time.Millisecond)
	cm.Ready(1)
	cm.Ready(2)
	cm.Ready(3)
	if cm.Len() != 3 {
		t.Fatalf("expected Len 3, got %d", cm.Len())
	}
}

func TestCoolMap_Purge_RemovesExpired(t *testing.T) {
	cm, _ := NewCoolMap(20 * time.Millisecond)
	cm.Ready(100)
	cm.Ready(200)
	time.Sleep(40 * time.Millisecond)
	cm.Purge()
	if cm.Len() != 0 {
		t.Fatalf("expected Len 0 after Purge, got %d", cm.Len())
	}
}
