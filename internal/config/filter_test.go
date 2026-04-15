package config

import "testing"

func TestDefaultFilterConfig_Valid(t *testing.T) {
	f := defaultFilterConfig()
	if err := f.validate(); err != nil {
		t.Fatalf("default filter config should be valid: %v", err)
	}
}

func TestFilterConfig_Validate_ExcludedPortOutOfRange(t *testing.T) {
	f := defaultFilterConfig()
	f.ExcludePorts = []int{0}
	if err := f.validate(); err == nil {
		t.Error("expected error for excluded port 0")
	}
	f.ExcludePorts = []int{70000}
	if err := f.validate(); err == nil {
		t.Error("expected error for excluded port 70000")
	}
}

func TestFilterConfig_Validate_ValidExcludedPorts(t *testing.T) {
	f := defaultFilterConfig()
	f.ExcludePorts = []int{22, 80, 443}
	if err := f.validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestFilterConfig_Validate_RangeMinExceedsMax(t *testing.T) {
	f := defaultFilterConfig()
	f.RangeMin = 9000
	f.RangeMax = 1000
	if err := f.validate(); err == nil {
		t.Error("expected error when range_min > range_max")
	}
}

func TestFilterConfig_Validate_RangeBelowOne(t *testing.T) {
	f := defaultFilterConfig()
	f.RangeMin = 0
	if err := f.validate(); err == nil {
		t.Error("expected error for range_min < 1")
	}
}

func TestFilterConfig_Validate_RangeAbove65535(t *testing.T) {
	f := defaultFilterConfig()
	f.RangeMax = 70000
	if err := f.validate(); err == nil {
		t.Error("expected error for range_max > 65535")
	}
}

func TestFilterConfig_Validate_NarrowRange(t *testing.T) {
	f := FilterConfig{
		ExcludePorts: []int{8080},
		RangeMin:     1024,
		RangeMax:     9000,
	}
	if err := f.validate(); err != nil {
		t.Errorf("unexpected error for valid narrow range: %v", err)
	}
}
