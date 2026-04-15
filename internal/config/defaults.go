package config

// Default returns a Config populated with sensible defaults.
// Useful when no config file is provided.
func Default() *Config {
	return &Config{
		ScanInterval: 30,
		PortRange: PortRange{
			Start: 1,
			End:   65535,
		},
		Protocols: []string{"tcp"},
		Rules:     []RuleConfig{},
		Alert: AlertConfig{
			Output: "stdout",
		},
	}
}
