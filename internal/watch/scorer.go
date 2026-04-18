package watch

import "errors"

// Scorer assigns a numeric risk score to a port based on configurable weights.
// Higher scores indicate ports that are more likely to be noteworthy.
type Scorer struct {
	weights map[int]float64
	defaultScore float64
}

// NewScorer creates a Scorer with the given default score.
// defaultScore must be >= 0.
func NewScorer(defaultScore float64) (*Scorer, error) {
	if defaultScore < 0 {
		return nil, errors.New("scorer: defaultScore must be >= 0")
	}
	return &Scorer{
		weights:      make(map[int]float64),
		defaultScore: defaultScore,
	}, nil
}

// SetWeight assigns a risk weight to a specific port.
// port must be in [1, 65535] and weight must be >= 0.
func (s *Scorer) SetWeight(port int, weight float64) error {
	if port < 1 || port > 65535 {
		return errors.New("scorer: port out of range")
	}
	if weight < 0 {
		return errors.New("scorer: weight must be >= 0")
	}
	s.weights[port] = weight
	return nil
}

// Score returns the risk score for the given port.
// If no weight is set for the port, the default score is returned.
func (s *Scorer) Score(port int) float64 {
	if w, ok := s.weights[port]; ok {
		return w
	}
	return s.defaultScore
}

// ScoreAll returns a map of port -> score for all provided ports.
func (s *Scorer) ScoreAll(ports []int) map[int]float64 {
	result := make(map[int]float64, len(ports))
	for _, p := range ports {
		result[p] = s.Score(p)
	}
	return result
}
