package watch

import (
	"testing"
	"time"
)

func TestNewSuppress_InvalidWindow(t *testing.T) {
	_, err := NewSuppress(0)
	if err == nil {
		t.Fatal("expected error for zero window")
	}
	_, err = NewSuppress(-time.Second)
	if err == nil {
		t.Fatal("expected error for negative window")
	}
}

func TestNewSuppress_Valid(t *testing.T) {
	s, err := NewSuppress(time.Minute)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil Suppress")
	}
}

func TestSuppress_Mute_InvalidPort(t *testing.T) {
	s, _ := NewSuppress(time.Minute)
	if err := s.Mute(0); err == nil {
		t.Error("expected error for port 0")
	}
	if err := s.Mute(65536); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestSuppress_IsMuted_FirstCall_ReturnsFalse(t *testing.T) {
	s, _ := NewSuppress(time.Minute)
	if s.IsMuted(443) {
		t.Error("expected port 443 to be unmuted initially")
	}
}

func TestSuppress_Mute_ThenIsMuted(t *testing.T) {
	s, _ := NewSuppress(time.Minute)
	_ = s.Mute(8080)
	if !s.IsMuted(8080) {
		t.Error("expected port 8080 to be muted")
	}
}

func TestSuppress_Unmute_ClearsMute(t *testing.T) {
	s, _ := NewSuppress(time.Minute)
	_ = s.Mute(9090)
	s.Unmute(9090)
	if s.IsMuted(9090) {
		t.Error("expected port 9090 to be unmuted after Unmute")
	}
}

func TestSuppress_IsMuted_AfterExpiry_ReturnsFalse(t *testing.T) {
	s, _ := NewSuppress(10 * time.Millisecond)
	_ = s.Mute(3000)
	time.Sleep(30 * time.Millisecond)
	if s.IsMuted(3000) {
		t.Error("expected port 3000 to be unmuted after window expiry")
	}
}

func TestSuppress_Len_CountsActive(t *testing.T) {
	s, _ := NewSuppress(time.Minute)
	_ = s.Mute(1000)
	_ = s.Mute(2000)
	if got := s.Len(); got != 2 {
		t.Errorf("expected Len 2, got %d", got)
	}
}

func TestSuppress_Len_EvictsExpired(t *testing.T) {
	s, _ := NewSuppress(10 * time.Millisecond)
	_ = s.Mute(5000)
	time.Sleep(30 * time.Millisecond)
	if got := s.Len(); got != 0 {
		t.Errorf("expected Len 0 after expiry, got %d", got)
	}
}
