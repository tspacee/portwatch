package config

import (
	"testing"
	"time"
)

func TestDefaultHorizonConfig_Valid(t *testing.T) {
	cfg := defaultHorizonConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestDefaultHorizonConfig_DisabledByDefault(t *testing.T) {
	cfg := defaultHorizonConfig()
	if cfg.Enabled {
		t.Error("expected horizon disabled by default")
	}
}

func TestDefaultHorizonConfig_DefaultCutoff(t *testing.T) {
	cfg := defaultHorizonConfig()
	if cfg.Cutoff != 24*time.Hour {
		t.Errorf("expected 24h cutoff, got %v", cfg.Cutoff)
	}
}

func TestHorizonConfig_Validate_Disabled(t *testing.T) {
	cfg := HorizonConfig{Enabled: false, Cutoff: 0}
	if err := cfg.Validate(); err != nil {
		t.Errorf("expected no error when disabled, got %v", err)
	}
}

func TestHorizonConfig_Validate_EnabledWithValidCutoff(t *testing.T) {
	cfg := HorizonConfig{Enabled: true, Cutoff: time.Hour}
	if err := cfg.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHorizonConfig_Validate_EnabledWithZeroCutoff(t *testing.T) {
	cfg := HorizonConfig{Enabled: true, Cutoff: 0}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for zero cutoff when enabled")
	}
}

func TestHorizonConfig_Validate_EnabledWithNegativeCutoff(t *testing.T) {
	cfg := HorizonConfig{Enabled: true, Cutoff: -time.Minute}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for negative cutoff when enabled")
	}
}
