package watch

import (
	"testing"
	"time"
)

func TestNewPoll_InvalidInterval(t *testing.T) {
	_, err := NewPoll(0)
	if err == nil {
		t.Fatal("expected error for zero interval")
	}
}

func TestNewPoll_Valid(t *testing.T) {
	p, err := NewPoll(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil Poll")
	}
}

func TestPoll_Due_NeverPolled(t *testing.T) {
	p, _ := NewPoll(time.Second)
	if !p.Due(8080) {
		t.Error("expected port never polled to be due")
	}
}

func TestPoll_Mark_And_NotDue(t *testing.T) {
	p, _ := NewPoll(time.Hour)
	if err := p.Mark(8080); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Due(8080) {
		t.Error("expected port to not be due immediately after mark")
	}
}

func TestPoll_Mark_InvalidPort(t *testing.T) {
	p, _ := NewPoll(time.Second)
	if err := p.Mark(0); err == nil {
		t.Error("expected error for invalid port")
	}
}

func TestPoll_Set_InvalidPort(t *testing.T) {
	p, _ := NewPoll(time.Second)
	if err := p.Set(99999, time.Second); err == nil {
		t.Error("expected error for port out of range")
	}
}

func TestPoll_Set_InvalidInterval(t *testing.T) {
	p, _ := NewPoll(time.Second)
	if err := p.Set(8080, 0); err == nil {
		t.Error("expected error for zero interval")
	}
}

func TestPoll_Interval_Default(t *testing.T) {
	p, _ := NewPoll(5 * time.Second)
	if got := p.Interval(1234); got != 5*time.Second {
		t.Errorf("expected default interval, got %v", got)
	}
}

func TestPoll_Interval_Custom(t *testing.T) {
	p, _ := NewPoll(5 * time.Second)
	_ = p.Set(443, 30*time.Second)
	if got := p.Interval(443); got != 30*time.Second {
		t.Errorf("expected custom interval 30s, got %v", got)
	}
}

func TestPoll_Due_AfterCustomInterval(t *testing.T) {
	p, _ := NewPoll(time.Hour)
	_ = p.Set(22, time.Millisecond)
	_ = p.Mark(22)
	time.Sleep(5 * time.Millisecond)
	if !p.Due(22) {
		t.Error("expected port to be due after custom interval elapsed")
	}
}
