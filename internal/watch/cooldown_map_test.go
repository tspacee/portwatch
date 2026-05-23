package watch

import (
	"testing"
	"time"
)

func TestNewExpiry_InvalidTTL(t *testing.T) {
	_, err := NewExpiry(0)
	if err == nil {
		t.Fatal("expected error for zero TTL")
	}
}

func TestNewExpiry_Valid(t *testing.T) {
	e, err := NewExpiry(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e == nil {
		t.Fatal("expected non-nil Expiry")
	}
}

func TestExpiry_Set_And_NotExpired(t *testing.T) {
	e, _ := NewExpiry(time.Hour)
	if err := e.Set(80); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if e.Expired(80) {
		t.Error("expected port 80 to not be expired")
	}
}

func TestExpiry_Expired_NeverSet(t *testing.T) {
	e, _ := NewExpiry(time.Hour)
	if !e.Expired(443) {
		t.Error("expected unset port to be expired")
	}
}

func TestExpiry_Expired_AfterTTL(t *testing.T) {
	e, _ := NewExpiry(10 * time.Millisecond)
	_ = e.Set(8080)
	time.Sleep(20 * time.Millisecond)
	if !e.Expired(8080) {
		t.Error("expected port 8080 to be expired after TTL")
	}
}

func TestExpiry_Delete_RemovesEntry(t *testing.T) {
	e, _ := NewExpiry(time.Hour)
	_ = e.Set(22)
	e.Delete(22)
	if !e.Expired(22) {
		t.Error("expected deleted port to be expired")
	}
}

func TestExpiry_Len_TracksEntries(t *testing.T) {
	e, _ := NewExpiry(time.Hour)
	_ = e.Set(22)
	_ = e.Set(80)
	if e.Len() != 2 {
		t.Errorf("expected len 2, got %d", e.Len())
	}
	e.Delete(22)
	if e.Len() != 1 {
		t.Errorf("expected len 1, got %d", e.Len())
	}
}

func TestExpiry_Set_InvalidPort(t *testing.T) {
	e, _ := NewExpiry(time.Hour)
	if err := e.Set(0); err == nil {
		t.Error("expected error for port 0")
	}
	if err := e.Set(65536); err == nil {
		t.Error("expected error for port 65536")
	}
}
