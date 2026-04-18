package watch

import (
	"net"
	"testing"
	"time"
)

func startProbeListener(t *testing.T) (int, func()) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return port, func() { _ = ln.Close() }
}

func TestNewProbe_InvalidTimeout(t *testing.T) {
	_, err := NewProbe(0)
	if err == nil {
		t.Fatal("expected error for zero timeout")
	}
}

func TestNewProbe_Valid(t *testing.T) {
	p, err := NewProbe(time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p == nil {
		t.Fatal("expected non-nil probe")
	}
}

func TestProbe_Check_OpenPort(t *testing.T) {
	port, stop := startProbeListener(t)
	defer stop()

	p, _ := NewProbe(500 * time.Millisecond)
	if !p.Check(port) {
		t.Errorf("expected port %d to be reachable", port)
	}
}

func TestProbe_Check_ClosedPort(t *testing.T) {
	p, _ := NewProbe(200 * time.Millisecond)
	// Port 1 is almost certainly closed in test environments.
	if p.Check(1) {
		t.Error("expected port 1 to be unreachable")
	}
}

func TestProbe_CheckAll_ReturnsOnlyOpen(t *testing.T) {
	port, stop := startProbeListener(t)
	defer stop()

	p, _ := NewProbe(500 * time.Millisecond)
	result := p.CheckAll([]int{1, port})

	if len(result) != 1 || result[0] != port {
		t.Errorf("expected only port %d, got %v", port, result)
	}
}

func TestProbe_CheckAll_EmptyInput(t *testing.T) {
	p, _ := NewProbe(time.Second)
	result := p.CheckAll([]int{})
	if len(result) != 0 {
		t.Errorf("expected empty result, got %v", result)
	}
}
