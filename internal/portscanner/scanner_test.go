package portscanner

import (
	"net"
	"testing"
	"time"
)

// startTestListener opens a TCP listener on an OS-assigned port and returns it.
func startTestListener(t *testing.T) (net.Listener, int) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to start test listener: %v", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	return ln, port
}

func TestNewScanner_Defaults(t *testing.T) {
	s := NewScanner(1, 1024)
	if s.StartPort != 1 || s.EndPort != 1024 {
		t.Errorf("unexpected port range: %d-%d", s.StartPort, s.EndPort)
	}
	if s.Protocol != "tcp" {
		t.Errorf("expected protocol tcp, got %s", s.Protocol)
	}
	if s.Workers <= 0 {
		t.Error("workers should be positive")
	}
}

func TestScan_InvalidRange(t *testing.T) {
	s := NewScanner(1000, 500)
	_, err := s.Scan()
	if err == nil {
		t.Error("expected error for invalid port range")
	}
}

func TestScan_DetectsOpenPort(t *testing.T) {
	ln, port := startTestListener(t)
	defer ln.Close()

	s := NewScanner(port, port)
	s.Timeout = 300 * time.Millisecond
	s.Workers = 1

	result, err := s.Scan()
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	open := OpenPorts(result)
	if len(open) != 1 || open[0].Port != port {
		t.Errorf("expected port %d to be open, got %+v", port, open)
	}
}

func TestScan_DetectsClosedPort(t *testing.T) {
	// Find a free port then immediately release it so it's closed during scan.
	ln, port := startTestListener(t)
	ln.Close()

	s := NewScanner(port, port)
	s.Timeout = 200 * time.Millisecond
	s.Workers = 1

	result, err := s.Scan()
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	open := OpenPorts(result)
	if len(open) != 0 {
		t.Errorf("expected no open ports, got %+v", open)
	}
}

func TestScanResult_Timestamp(t *testing.T) {
	before := time.Now()
	s := NewScanner(65534, 65535)
	s.Timeout = 50 * time.Millisecond
	result, err := s.Scan()
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}
	if result.Timestamp.Before(before) {
		t.Error("result timestamp should be after scan start")
	}
}
