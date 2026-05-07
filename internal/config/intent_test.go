package config

import "testing"

func TestDefaultIntentConfig_Valid(t *testing.T) {
	c := defaultIntentConfig()
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDefaultIntentConfig_DisabledByDefault(t *testing.T) {
	c := defaultIntentConfig()
	if c.Enabled {
		t.Fatal("expected intent to be disabled by default")
	}
}

func TestDefaultIntentConfig_EmptyPorts(t *testing.T) {
	c := defaultIntentConfig()
	if len(c.Ports) != 0 {
		t.Fatalf("expected no ports, got %d", len(c.Ports))
	}
}

func TestIntentConfig_Validate_Disabled(t *testing.T) {
	c := IntentConfig{Enabled: false, Ports: []IntentEntry{{Port: 0, Reason: ""}}}
	if err := c.Validate(); err != nil {
		t.Fatalf("disabled config should not validate entries: %v", err)
	}
}

func TestIntentConfig_Validate_InvalidPort(t *testing.T) {
	c := IntentConfig{
		Enabled: true,
		Ports:   []IntentEntry{{Port: 0, Reason: "bad"}},
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for port 0")
	}
}

func TestIntentConfig_Validate_EmptyReason(t *testing.T) {
	c := IntentConfig{
		Enabled: true,
		Ports:   []IntentEntry{{Port: 80, Reason: ""}},
	}
	if err := c.Validate(); err == nil {
		t.Fatal("expected error for empty reason")
	}
}

func TestIntentConfig_Validate_Valid(t *testing.T) {
	c := IntentConfig{
		Enabled: true,
		Ports: []IntentEntry{
			{Port: 443, Reason: "https"},
			{Port: 22, Reason: "ssh"},
		},
	}
	if err := c.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
