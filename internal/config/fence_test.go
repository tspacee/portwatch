package config

import "testing"

func TestDefaultFenceConfig_Valid(t *testing.T) {
	cfg := defaultFenceConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFenceConfig_Validate_DisabledEmptyPorts_Valid(t *testing.T) {
	cfg := FenceConfig{Enabled: false, AllowedPorts: []int{}}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestFenceConfig_Validate_EnabledEmptyPorts_Invalid(t *testing.T) {
	cfg := FenceConfig{Enabled: true, AllowedPorts: []int{}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for enabled fence with empty allowed_ports")
	}
}

func TestFenceConfig_Validate_InvalidPort(t *testing.T) {
	cfg := FenceConfig{Enabled: true, AllowedPorts: []int{0}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestFenceConfig_Validate_PortTooHigh(t *testing.T) {
	cfg := FenceConfig{Enabled: true, AllowedPorts: []int{99999}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port exceeding 65535")
	}
}

func TestFenceConfig_Validate_ValidPorts(t *testing.T) {
	cfg := FenceConfig{Enabled: true, AllowedPorts: []int{80, 443, 8080}}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
