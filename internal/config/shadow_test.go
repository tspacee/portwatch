package config

import "testing"

func TestDefaultShadowConfig_Valid(t *testing.T) {
	cfg := defaultShadowConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultShadowConfig_EnabledByDefault(t *testing.T) {
	cfg := defaultShadowConfig()
	if !cfg.Enabled {
		t.Fatal("expected shadow to be enabled by default")
	}
}

func TestShadowConfig_Validate_Disabled(t *testing.T) {
	cfg := ShadowConfig{Enabled: false}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
