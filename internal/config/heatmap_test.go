package config

import "testing"

func TestDefaultHeatMapConfig_Valid(t *testing.T) {
	cfg := defaultHeatMapConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid default config, got: %v", err)
	}
}

func TestDefaultHeatMapConfig_EnabledByDefault(t *testing.T) {
	cfg := defaultHeatMapConfig()
	if !cfg.Enabled {
		t.Fatal("expected heatmap to be enabled by default")
	}
}

func TestDefaultHeatMapConfig_DefaultTopN(t *testing.T) {
	cfg := defaultHeatMapConfig()
	if cfg.TopN != 10 {
		t.Fatalf("expected TopN 10, got %d", cfg.TopN)
	}
}

func TestHeatMapConfig_Validate_Disabled(t *testing.T) {
	cfg := HeatMapConfig{Enabled: false, TopN: 0}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid, got: %v", err)
	}
}

func TestHeatMapConfig_Validate_ZeroTopN(t *testing.T) {
	cfg := HeatMapConfig{Enabled: true, TopN: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for TopN=0")
	}
}

func TestHeatMapConfig_Validate_ValidTopN(t *testing.T) {
	cfg := HeatMapConfig{Enabled: true, TopN: 5}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHeatMapConfig_Validate_TopNExceedsMax(t *testing.T) {
	cfg := HeatMapConfig{Enabled: true, TopN: 70000}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for TopN exceeding max port count")
	}
}
