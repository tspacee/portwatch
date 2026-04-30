package watch

import (
	"testing"
	"time"
)

func TestNewDrain_InvalidWindow(t *testing.T) {
	_, err := NewDrain(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestNewDrain_Valid(t *testing.T) {
	d, err := NewDrain(100 * time.Millisecond)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Len() != 0 {
		t.Fatal("expected empty drain")
	}
}

func TestDrain_Stage_InvalidPort(t *testing.T) {
	d, _ := NewDrain(100 * time.Millisecond)
	if err := d.Stage(0); err == nil {
		t.Fatal("expected error for port 0")
	}
	if err := d.Stage(65536); err == nil {
		t.Fatal("expected error for port 65536")
	}
}

func TestDrain_Stage_IncrementsLen(t *testing.T) {
	d, _ := NewDrain(100 * time.Millisecond)
	if err := d.Stage(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if d.Len() != 1 {
		t.Fatalf("expected len 1, got %d", d.Len())
	}
}

func TestDrain_Stage_Idempotent(t *testing.T) {
	d, _ := NewDrain(100 * time.Millisecond)
	_ = d.Stage(8080)
	_ = d.Stage(8080)
	if d.Len() != 1 {
		t.Fatalf("expected len 1 after duplicate stage, got %d", d.Len())
	}
}

func TestDrain_Unstage_RemovesPort(t *testing.T) {
	d, _ := NewDrain(100 * time.Millisecond)
	_ = d.Stage(9090)
	d.Unstage(9090)
	if d.Len() != 0 {
		t.Fatal("expected empty drain after unstage")
	}
}

func TestDrain_Drained_BeforeWindow_ReturnsEmpty(t *testing.T) {
	d, _ := NewDrain(500 * time.Millisecond)
	_ = d.Stage(443)
	result := d.Drained()
	if len(result) != 0 {
		t.Fatalf("expected no drained ports before window, got %v", result)
	}
	if d.Len() != 1 {
		t.Fatal("port should still be staged")
	}
}

func TestDrain_Drained_AfterWindow_ReturnsPort(t *testing.T) {
	d, _ := NewDrain(20 * time.Millisecond)
	_ = d.Stage(22)
	time.Sleep(40 * time.Millisecond)
	result := d.Drained()
	if len(result) != 1 || result[0] != 22 {
		t.Fatalf("expected port 22 drained, got %v", result)
	}
	if d.Len() != 0 {
		t.Fatal("port should be removed after draining")
	}
}
