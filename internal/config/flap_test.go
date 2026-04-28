package config

import (
	"testing"
	"time"
)

func TestDefaultFlapConfig_Valid(t *testing.T) {
	cfg := defaultFlapConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid default config, got: %v", err)
	}
}

func TestDefaultFlapConfig_EnabledByDefault(t *testing.T) {
	cfg := defaultFlapConfig()
	if !cfg.Enabled {
		t.Fatal("expected flap detection to be enabled by default")
	}
}

func TestDefaultFlapConfig_DefaultThreshold(t *testing.T) {
	cfg := defaultFlapConfig()
	if cfg.Threshold != 4 {
		t.Fatalf("expected default threshold 4, got %d", cfg.Threshold)
	}
}

func TestFlapConfig_Validate_Disabled(t *testing.T) {
	cfg := FlapConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid, got: %v", err)
	}
}

func TestFlapConfig_Validate_ZeroWindow(t *testing.T) {
	cfg := FlapConfig{Enabled: true, Window: 0, Threshold: 3}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestFlapConfig_Validate_ThresholdTooLow(t *testing.T) {
	cfg := FlapConfig{Enabled: true, Window: time.Minute, Threshold: 1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for threshold < 2")
	}
}

func TestFlapConfig_Validate_ValidCustom(t *testing.T) {
	cfg := FlapConfig{Enabled: true, Window: 10 * time.Second, Threshold: 5}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}
