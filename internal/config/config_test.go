package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/portwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "portwatch-*.yaml")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
scan_interval_seconds: 60
port_range:
  start: 1024
  end: 9000
protocols:
  - tcp
  - udp
rules:
  - name: allow-ssh
    port: 22
    protocol: tcp
    action: allow
alert:
  output: stdout
`
	path := writeTemp(t, yaml)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 60 {
		t.Errorf("expected scan_interval 60, got %d", cfg.ScanInterval)
	}
	if cfg.PortRange.Start != 1024 || cfg.PortRange.End != 9000 {
		t.Errorf("unexpected port range: %+v", cfg.PortRange)
	}
	if len(cfg.Rules) != 1 || cfg.Rules[0].Name != "allow-ssh" {
		t.Errorf("unexpected rules: %+v", cfg.Rules)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidPortRange(t *testing.T) {
	yaml := `
port_range:
  start: 9000
  end: 1024
`
	path := writeTemp(t, yaml)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected validation error for invalid port range")
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTemp(t, "{}")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 30 {
		t.Errorf("expected default scan_interval 30, got %d", cfg.ScanInterval)
	}
	if len(cfg.Protocols) == 0 || cfg.Protocols[0] != "tcp" {
		t.Errorf("expected default protocol tcp, got %v", cfg.Protocols)
	}
}

func TestDefault(t *testing.T) {
	cfg := config.Default()
	if cfg.ScanInterval != 30 {
		t.Errorf("expected 30, got %d", cfg.ScanInterval)
	}
	if cfg.PortRange.Start != 1 || cfg.PortRange.End != 65535 {
		t.Errorf("unexpected default port range: %+v", cfg.PortRange)
	}
}
