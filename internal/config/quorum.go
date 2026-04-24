package config

import "errors"

// QuorumConfig controls the multi-source quorum feature.
type QuorumConfig struct {
	Enabled  bool `yaml:"enabled"`
	MinVotes int  `yaml:"min_votes"`
}

// defaultQuorumConfig returns a conservative default: quorum is
// disabled and requires two votes when enabled.
func defaultQuorumConfig() QuorumConfig {
	return QuorumConfig{
		Enabled:  false,
		MinVotes: 2,
	}
}

// Validate checks that the QuorumConfig is self-consistent.
func (c QuorumConfig) Validate() error {
	if !c.Enabled {
		return nil
	}
	if c.MinVotes < 1 {
		return errors.New("quorum: min_votes must be at least 1 when enabled")
	}
	return nil
}
