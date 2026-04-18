package config_test

import (
	"testing"

	"github.com/user/portwatch/internal/config"
)

func TestDefaultLabelerConfig_Valid(t *testing.T) {
	cfg := config.LabelerConfig{
		Labels: []config.PortLabel{
			{Port: 22, Label: "ssh"},
			{Port: 80, Label: "http"},
		},
	}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid config, got: %v", err)
	}
}

func TestLabelerConfig_Validate_InvalidPort(t *testing.T) {
	cfg := config.LabelerConfig{
		Labels: []config.PortLabel{
			{Port: 0, Label: "zero"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestLabelerConfig_Validate_PortTooHigh(t *testing.T) {
	cfg := config.LabelerConfig{
		Labels: []config.PortLabel{
			{Port: 70000, Label: "toobig"},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for port 70000")
	}
}

func TestLabelerConfig_Validate_EmptyLabel(t *testing.T) {
	cfg := config.LabelerConfig{
		Labels: []config.PortLabel{
			{Port: 8080, Label: ""},
		},
	}
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty label")
	}
}

func TestLabelerConfig_Validate_Empty_IsValid(t *testing.T) {
	cfg := config.LabelerConfig{}
	if err := cfg.Validate(); err != nil {
		t.Fatalf("empty config should be valid, got: %v", err)
	}
}
