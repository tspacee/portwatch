package config

import (
	"testing"
	"time"
)

func TestDefaultProbeRegistryConfig_Valid(t *testing.T) {
	c := defaultProbeRegistryConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultProbeRegistryConfig_DisabledByDefault(t *testing.T) {
	c := defaultProbeRegistryConfig()
	if c.Enabled {
		t.Fatal("expected disabled by default")
	}
}

func TestDefaultProbeRegistryConfig_DefaultTimeout(t *testing.T) {
	c := defaultProbeRegistryConfig()
	if c.Timeout != 2*time.Second {
		t.Fatalf("expected 2s timeout, got %v", c.Timeout)
	}
}

func TestProbeRegistryConfig_Validate_Disabled(t *testing.T) {
	c := ProbeRegistryConfig{Enabled: false, Timeout: 0}
	if err := c.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid: %v", err)
	}
}

func TestProbeRegistryConfig_Validate_EnabledWithValidTimeout(t *testing.T) {
	c := ProbeRegistryConfig{Enabled: true, Timeout: 5 * time.Second}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestProbeRegistryConfig_Validate_EnabledWithZeroTimeout(t *testing.T) {
	c := ProbeRegistryConfig{Enabled: true, Timeout: 0}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for zero timeout when enabled")
	}
}

func TestProbeRegistryConfig_Validate_TimeoutExceedsMax(t *testing.T) {
	c := ProbeRegistryConfig{Enabled: true, Timeout: 60 * time.Second}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for timeout exceeding 30s")
	}
}
