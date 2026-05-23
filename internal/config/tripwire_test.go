package config

import "testing"

func TestDefaultTripwireConfig_Valid(t *testing.T) {
	cfg := defaultTripwireConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid default config, got: %v", err)
	}
}

func TestDefaultTripwireConfig_DisabledByDefault(t *testing.T) {
	cfg := defaultTripwireConfig()
	if cfg.Enabled {
		t.Fatal("expected disabled by default")
	}
}

func TestDefaultTripwireConfig_EmptyPorts(t *testing.T) {
	cfg := defaultTripwireConfig()
	if len(cfg.Ports) != 0 {
		t.Fatal("expected empty ports by default")
	}
}

func TestTripwireConfig_Validate_Disabled(t *testing.T) {
	cfg := TripwireConfig{Enabled: false, Ports: []int{}}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("disabled config should always be valid, got: %v", err)
	}
}

func TestTripwireConfig_Validate_EnabledWithValidPorts(t *testing.T) {
	cfg := TripwireConfig{Enabled: true, Ports: []int{22, 80, 443}}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid, got: %v", err)
	}
}

func TestTripwireConfig_Validate_EnabledNoPorts(t *testing.T) {
	cfg := TripwireConfig{Enabled: true, Ports: []int{}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for enabled config with no ports")
	}
}

func TestTripwireConfig_Validate_InvalidPort(t *testing.T) {
	cfg := TripwireConfig{Enabled: true, Ports: []int{0}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestTripwireConfig_Validate_PortTooHigh(t *testing.T) {
	cfg := TripwireConfig{Enabled: true, Ports: []int{65536}}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 65536")
	}
}
