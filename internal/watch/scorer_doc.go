// Package watch provides primitives for monitoring open ports.
//
// # Scorer
//
// Scorer assigns numeric risk scores to ports based on operator-configured
// weights. It is useful for prioritising alerts when many ports change state
// simultaneously.
//
// Usage:
//
//	s, err := watch.NewScorer(1.0)   // default score for unknown ports
//	if err != nil { ... }
//	_ = s.SetWeight(22, 9.0)         // SSH is high-risk
//	_ = s.SetWeight(80, 3.0)         // HTTP is medium-risk
//
//	scores := s.ScoreAll(openPorts)  // map[int]float64
//
// Scores are purely informational; how they are acted upon is left to the
// caller (e.g. filtering, sorting, or annotating alerts).
package watch
