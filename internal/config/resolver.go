package config

import "fmt"

// ResolverEntry maps a port to a service name in configuration.
type ResolverEntry struct {
	Port    int    `yaml:"port"`
	Service string `yaml:"service"`
}

// ResolverConfig holds custom port-to-service mappings.
type ResolverConfig struct {
	Custom []ResolverEntry `yaml:"custom"`
}

// defaultResolverConfig returns an empty resolver config.
func defaultResolverConfig() ResolverConfig {
	return ResolverConfig{}
}

// Validate checks that all entries have valid ports and non-empty service names.
func (c ResolverConfig) Validate() error {
	for _, e := range c.Custom {
		if e.Port < 1 || e.Port > 65535 {
			return fmt.Errorf("resolver: invalid port %d", e.Port)
		}
		if e.Service == "" {
			return fmt.Errorf("resolver: service name must not be empty for port %d", e.Port)
		}
	}
	return nil
}
