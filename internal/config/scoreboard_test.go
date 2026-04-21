package config

import "testing"

func TestDefaultScoreboardConfig_Valid(t *testing.T) {
	cfg := defaultScoreboardConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid default config, got: %v", err)
	}
}

func TestDefaultScoreboardConfig_EnabledByDefault(t *testing.T) {
	cfg := defaultScoreboardConfig()
	if !cfg.Enabled {
		t.Fatal("expected scoreboard to be enabled by default")
	}
}

func TestDefaultScoreboardConfig_DefaultTopN(t *testing.T) {
	cfg := defaultScoreboardConfig()
	if cfg.TopN != 10 {
		t.Fatalf("expected TopN 10, got %d", cfg.TopN)
	}
}

func TestScoreboardConfig_Validate_Disabled(t *testing.T) {
	cfg := ScoreboardConfig{Enabled: false, TopN: 0}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid, got: %v", err)
	}
}

func TestScoreboardConfig_Validate_ZeroTopN(t *testing.T) {
	cfg := ScoreboardConfig{Enabled: true, TopN: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for TopN=0")
	}
}

func TestScoreboardConfig_Validate_NegativeTopN(t *testing.T) {
	cfg := ScoreboardConfig{Enabled: true, TopN: -1}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for negative TopN")
	}
}

func TestScoreboardConfig_Validate_ValidTopN(t *testing.T) {
	cfg := ScoreboardConfig{Enabled: true, TopN: 25}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestScoreboardConfig_Validate_TopNAtMaxPort(t *testing.T) {
	cfg := ScoreboardConfig{Enabled: true, TopN: 65535}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error at max: %v", err)
	}
}

func TestScoreboardConfig_Validate_TopNExceedsMax(t *testing.T) {
	cfg := ScoreboardConfig{Enabled: true, TopN: 65536}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for TopN exceeding 65535")
	}
}
