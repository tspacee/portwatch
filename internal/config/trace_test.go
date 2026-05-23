package config

import "testing"

func TestDefaultTraceConfig_Valid(t *testing.T) {
	cfg := defaultTraceConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected validation error: %v", err)
	}
}

func TestDefaultTraceConfig_DisabledByDefault(t *testing.T) {
	cfg := defaultTraceConfig()
	if cfg.Enabled {
		t.Fatal("expected trace to be disabled by default")
	}
}

func TestDefaultTraceConfig_DefaultMaxEntries(t *testing.T) {
	cfg := defaultTraceConfig()
	if cfg.MaxEntries != 100 {
		t.Fatalf("expected default max_entries 100, got %d", cfg.MaxEntries)
	}
}

func TestTraceConfig_Validate_Disabled(t *testing.T) {
	cfg := TraceConfig{Enabled: false, MaxEntries: 0}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid, got: %v", err)
	}
}

func TestTraceConfig_Validate_EnabledWithValidMaxEntries(t *testing.T) {
	cfg := TraceConfig{Enabled: true, MaxEntries: 50}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestTraceConfig_Validate_EnabledWithZeroMaxEntries(t *testing.T) {
	cfg := TraceConfig{Enabled: true, MaxEntries: 0}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for zero max_entries when enabled")
	}
}
