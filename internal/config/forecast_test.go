package config

import "testing"

func TestDefaultForecastConfig_Valid(t *testing.T) {
	cfg := defaultForecastConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultForecastConfig_DisabledByDefault(t *testing.T) {
	cfg := defaultForecastConfig()
	if cfg.Enabled {
		t.Fatal("expected forecast to be disabled by default")
	}
}

func TestDefaultForecastConfig_DefaultWindow(t *testing.T) {
	cfg := defaultForecastConfig()
	if cfg.Window != 10 {
		t.Fatalf("expected window=10, got %d", cfg.Window)
	}
}

func TestForecastConfig_Validate_Disabled(t *testing.T) {
	cfg := ForecastConfig{Enabled: false, Window: 0, Threshold: -1}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always pass validation, got: %v", err)
	}
}

func TestForecastConfig_Validate_ZeroWindow(t *testing.T) {
	cfg := ForecastConfig{Enabled: true, Window: 0, Threshold: 0.5}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero window")
	}
}

func TestForecastConfig_Validate_ThresholdBelowZero(t *testing.T) {
	cfg := ForecastConfig{Enabled: true, Window: 5, Threshold: -0.1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative threshold")
	}
}

func TestForecastConfig_Validate_ThresholdAboveOne(t *testing.T) {
	cfg := ForecastConfig{Enabled: true, Window: 5, Threshold: 1.1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for threshold > 1")
	}
}

func TestForecastConfig_Validate_EnabledWithValidValues(t *testing.T) {
	cfg := ForecastConfig{Enabled: true, Window: 5, Threshold: 0.75}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
