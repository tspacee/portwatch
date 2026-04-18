package config

import "testing"

func TestDefaultResolverConfig_Valid(t *testing.T) {
	c := defaultResolverConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolverConfig_Validate_InvalidPort(t *testing.T) {
	c := ResolverConfig{
		Custom: []ResolverEntry{{Port: 0, Service: "zero"}},
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for port 0")
	}
}

func TestResolverConfig_Validate_PortTooHigh(t *testing.T) {
	c := ResolverConfig{
		Custom: []ResolverEntry{{Port: 70000, Service: "overflow"}},
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for port 70000")
	}
}

func TestResolverConfig_Validate_EmptyService(t *testing.T) {
	c := ResolverConfig{
		Custom: []ResolverEntry{{Port: 8080, Service: ""}},
	}
	if err := c.Validate(); err == nil {
		t.Error("expected error for empty service")
	}
}

func TestResolverConfig_Validate_MultipleEntries_Valid(t *testing.T) {
	c := ResolverConfig{
		Custom: []ResolverEntry{
			{Port: 9090, Service: "prometheus"},
			{Port: 9091, Service: "grafana"},
		},
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
