package config

import (
	"testing"
	"time"
)

func TestDefaultTombstoneConfig_Valid(t *testing.T) {
	cfg := defaultTombstoneConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid default config, got: %v", err)
	}
}

func TestDefaultTombstoneConfig_DisabledByDefault(t *testing.T) {
	cfg := defaultTombstoneConfig()
	if cfg.Enabled {
		t.Fatal("expected tombstone to be disabled by default")
	}
}

func TestDefaultTombstoneConfig_DefaultTTL(t *testing.T) {
	cfg := defaultTombstoneConfig()
	if cfg.TTL != 30*time.Minute {
		t.Fatalf("expected 30m TTL, got %s", cfg.TTL)
	}
}

func TestTombstoneConfig_Validate_Disabled(t *testing.T) {
	cfg := TombstoneConfig{Enabled: false, TTL: 0}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid, got: %v", err)
	}
}

func TestTombstoneConfig_Validate_EnabledWithValidTTL(t *testing.T) {
	cfg := TombstoneConfig{Enabled: true, TTL: time.Hour, Ports: []int{80, 443}}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}

func TestTombstoneConfig_Validate_ZeroTTL(t *testing.T) {
	cfg := TombstoneConfig{Enabled: true, TTL: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero TTL when enabled")
	}
}

func TestTombstoneConfig_Validate_InvalidPort(t *testing.T) {
	cfg := TombstoneConfig{Enabled: true, TTL: time.Minute, Ports: []int{0}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for invalid port 0")
	}
}

func TestTombstoneConfig_Validate_PortTooHigh(t *testing.T) {
	cfg := TombstoneConfig{Enabled: true, TTL: time.Minute, Ports: []int{65536}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 65536")
	}
}
