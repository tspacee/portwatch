package config

import "testing"

func TestDefaultSequenceConfig_Valid(t *testing.T) {
	c := defaultSequenceConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDefaultSequenceConfig_DisabledByDefault(t *testing.T) {
	c := defaultSequenceConfig()
	if c.Enabled {
		t.Fatal("expected sequence to be disabled by default")
	}
}

func TestSequenceConfig_Validate_Disabled(t *testing.T) {
	c := SequenceConfig{Enabled: false}
	if err := c.Validate(); err != nil {
		t.Fatalf("disabled config should be valid, got %v", err)
	}
}

func TestSequenceConfig_Validate_Enabled(t *testing.T) {
	c := SequenceConfig{Enabled: true}
	if err := c.Validate(); err != nil {
		t.Fatalf("enabled config should be valid, got %v", err)
	}
}
