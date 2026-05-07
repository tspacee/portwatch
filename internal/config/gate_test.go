package config

import "testing"

func TestDefaultGateConfig_Valid(t *testing.T) {
	cfg := defaultGateConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid default, got: %v", err)
	}
}

func TestDefaultGateConfig_DisabledByDefault(t *testing.T) {
	cfg := defaultGateConfig()
	if cfg.Enabled {
		t.Error("expected gate to be disabled by default")
	}
}

func TestDefaultGateConfig_ClosedByDefault(t *testing.T) {
	cfg := defaultGateConfig()
	if cfg.OpenByDefault {
		t.Error("expected gate to be closed by default")
	}
}

func TestGateConfig_Validate_Disabled(t *testing.T) {
	cfg := GateConfig{Enabled: false, AllowedPorts: []int{0}}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled gate should skip validation, got: %v", err)
	}
}

func TestGateConfig_Validate_InvalidPort(t *testing.T) {
	cfg := GateConfig{Enabled: true, AllowedPorts: []int{0}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for port 0")
	}
}

func TestGateConfig_Validate_PortTooHigh(t *testing.T) {
	cfg := GateConfig{Enabled: true, AllowedPorts: []int{65536}}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for port 65536")
	}
}

func TestGateConfig_Validate_ValidPorts(t *testing.T) {
	cfg := GateConfig{Enabled: true, AllowedPorts: []int{80, 443, 8080}}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected valid config, got: %v", err)
	}
}
