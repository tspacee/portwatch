package config

import "testing"

func TestDefaultCapConfig_Valid(t *testing.T) {
	c := defaultCapConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid default, got: %v", err)
	}
}

func TestDefaultCapConfig_DisabledByDefault(t *testing.T) {
	c := defaultCapConfig()
	if c.Enabled {
		t.Fatal("expected cap to be disabled by default")
	}
}

func TestDefaultCapConfig_DefaultMax(t *testing.T) {
	c := defaultCapConfig()
	if c.Max != 256 {
		t.Fatalf("expected default max=256, got %d", c.Max)
	}
}

func TestCapConfig_Validate_Disabled(t *testing.T) {
	c := CapConfig{Enabled: false, Max: 0}
	if err := c.Validate(); err != nil {
		t.Fatalf("disabled cap with zero max should be valid, got: %v", err)
	}
}

func TestCapConfig_Validate_EnabledWithValidMax(t *testing.T) {
	c := CapConfig{Enabled: true, Max: 50}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}

func TestCapConfig_Validate_EnabledWithZeroMax(t *testing.T) {
	c := CapConfig{Enabled: true, Max: 0}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for enabled cap with max=0")
	}
}
