package config

import "testing"

func TestDefaultNormConfig_Valid(t *testing.T) {
	c := defaultNormConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultNormConfig_DisabledByDefault(t *testing.T) {
	c := defaultNormConfig()
	if c.Enabled {
		t.Fatal("expected norm to be disabled by default")
	}
}

func TestDefaultNormConfig_EmptyPorts(t *testing.T) {
	c := defaultNormConfig()
	if len(c.Ports) != 0 {
		t.Fatalf("expected empty ports slice, got %d entries", len(c.Ports))
	}
}

func TestNormConfig_Validate_Disabled(t *testing.T) {
	c := NormConfig{Enabled: false, Ports: []int{0}}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected no error when disabled, got: %v", err)
	}
}

func TestNormConfig_Validate_EnabledWithValidPorts(t *testing.T) {
	c := NormConfig{Enabled: true, Ports: []int{80, 443, 8080}}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNormConfig_Validate_InvalidPort(t *testing.T) {
	c := NormConfig{Enabled: true, Ports: []int{80, 0}}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestNormConfig_Validate_PortTooHigh(t *testing.T) {
	c := NormConfig{Enabled: true, Ports: []int{65536}}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for port 65536")
	}
}
