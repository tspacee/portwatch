package config

import (
	"testing"
	"time"
)

func TestDefaultExpiryConfig_Valid(t *testing.T) {
	c := defaultExpiryConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultExpiryConfig_DisabledByDefault(t *testing.T) {
	c := defaultExpiryConfig()
	if c.Enabled {
		t.Error("expected expiry to be disabled by default")
	}
}

func TestDefaultExpiryConfig_DefaultTTL(t *testing.T) {
	c := defaultExpiryConfig()
	if c.TTL != 30*time.Second {
		t.Errorf("expected 30s TTL, got %v", c.TTL)
	}
}

func TestExpiryConfig_Validate_Disabled(t *testing.T) {
	c := ExpiryConfig{Enabled: false, TTL: 0}
	if err := c.Validate(); err != nil {
		t.Fatalf("expected no error for disabled config, got %v", err)
	}
}

func TestExpiryConfig_Validate_EnabledWithValidTTL(t *testing.T) {
	c := ExpiryConfig{Enabled: true, TTL: time.Minute}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestExpiryConfig_Validate_EnabledWithZeroTTL(t *testing.T) {
	c := ExpiryConfig{Enabled: true, TTL: 0}
	if err := c.Validate(); err == nil {
		t.Error("expected error for zero TTL when enabled")
	}
}

func TestExpiryConfig_Validate_EnabledWithNegativeTTL(t *testing.T) {
	c := ExpiryConfig{Enabled: true, TTL: -time.Second}
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative TTL when enabled")
	}
}
