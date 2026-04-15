package config

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config-*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	_ = f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTemp(t, `
scan_interval: 10
port_range:
  start: 1
  end: 1024
protocols: [tcp]
snapshot_path: /tmp/snap.json
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.ScanInterval != 10 {
		t.Errorf("expected scan_interval 10, got %d", cfg.ScanInterval)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidPortRange(t *testing.T) {
	path := writeTemp(t, `
scan_interval: 5
port_range:
  start: 1000
  end: 100
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected error for invalid port range")
	}
}

func TestLoad_DefaultsApplied(t *testing.T) {
	path := writeTemp(t, `
port_range:
  start: 1
  end: 1024
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	def := Default()
	if cfg.ScanInterval != def.ScanInterval {
		t.Errorf("expected default scan_interval %d, got %d", def.ScanInterval, cfg.ScanInterval)
	}
}

func TestLoad_EmailAlertConfig(t *testing.T) {
	path := writeTemp(t, `
scan_interval: 5
port_range:
  start: 1
  end: 1024
alerts:
  email:
    enabled: true
    smtp_host: smtp.example.com
    smtp_port: 465
    from: noreply@example.com
    to:
      - admin@example.com
`)
	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.Alerts.Email.Enabled {
		t.Error("expected email alerts to be enabled")
	}
	if cfg.Alerts.Email.SMTPHost != "smtp.example.com" {
		t.Errorf("unexpected smtp host: %s", cfg.Alerts.Email.SMTPHost)
	}
	if len(cfg.Alerts.Email.To) != 1 || cfg.Alerts.Email.To[0] != "admin@example.com" {
		t.Errorf("unexpected recipients: %v", cfg.Alerts.Email.To)
	}
}
