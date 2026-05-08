package config

import "testing"

func TestDefaultBandConfig_Valid(t *testing.T) {
	cfg := defaultBandConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestDefaultBandConfig_EnabledByDefault(t *testing.T) {
	cfg := defaultBandConfig()
	if !cfg.Enabled {
		t.Error("expected band config to be enabled by default")
	}
}

func TestDefaultBandConfig_ThreeBands(t *testing.T) {
	cfg := defaultBandConfig()
	if len(cfg.Bands) != 3 {
		t.Errorf("expected 3 default bands, got %d", len(cfg.Bands))
	}
}

func TestBandConfig_Validate_Disabled(t *testing.T) {
	cfg := BandConfig{Enabled: false, Bands: []BandEntry{{Name: "", Min: 0, Max: 0}}}
	if err := cfg.Validate(); err != nil {
		t.Errorf("disabled config should skip validation, got: %v", err)
	}
}

func TestBandConfig_Validate_EmptyName(t *testing.T) {
	cfg := BandConfig{
		Enabled: true,
		Bands:   []BandEntry{{Name: "", Min: 1, Max: 1023}},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for empty band name")
	}
}

func TestBandConfig_Validate_InvalidRange(t *testing.T) {
	cfg := BandConfig{
		Enabled: true,
		Bands:   []BandEntry{{Name: "bad", Min: 500, Max: 100}},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for inverted range")
	}
}

func TestBandConfig_Validate_PortBelowOne(t *testing.T) {
	cfg := BandConfig{
		Enabled: true,
		Bands:   []BandEntry{{Name: "zero", Min: 0, Max: 1023}},
	}
	if err := cfg.Validate(); err == nil {
		t.Error("expected error for port 0")
	}
}
