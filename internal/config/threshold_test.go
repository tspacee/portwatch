package config

import "testing"

func TestDefaultThresholdConfig_Valid(t *testing.T) {
	cfg := defaultThresholdConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid default, got: %v", err)
	}
}

func TestDefaultThresholdConfig_DisabledByDefault(t *testing.T) {
	cfg := defaultThresholdConfig()
	if cfg.Enabled {
		t.Fatal("expected threshold disabled by default")
	}
}

func TestDefaultThresholdConfig_DefaultLimit(t *testing.T) {
	cfg := defaultThresholdConfig()
	if cfg.Limit != 5 {
		t.Fatalf("expected default limit 5, got %d", cfg.Limit)
	}
}

func TestThresholdConfig_Validate_EnabledWithValidLimit(t *testing.T) {
	cfg := ThresholdConfig{Enabled: true, Limit: 3}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}

func TestThresholdConfig_Validate_EnabledWithZeroLimit(t *testing.T) {
	cfg := ThresholdConfig{Enabled: true, Limit: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for limit=0 when enabled")
	}
}

func TestThresholdConfig_Validate_EnabledWithNegativeLimit(t *testing.T) {
	cfg := ThresholdConfig{Enabled: true, Limit: -1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative limit when enabled")
	}
}

func TestThresholdConfig_Validate_DisabledWithZeroLimit_Valid(t *testing.T) {
	cfg := ThresholdConfig{Enabled: false, Limit: 0}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid when disabled, got: %v", err)
	}
}
