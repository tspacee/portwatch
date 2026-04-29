package watch

import (
	"errors"
	"sync"
)

// Anomaly tracks ports that deviate from their expected baseline frequency.
// A port is considered anomalous when its observed rate falls outside the
// configured tolerance band around its historical mean.
type Anomaly struct {
	mu        sync.Mutex
	tolerance float64
	baseline  map[int]float64
	anomalous map[int]bool
}

// NewAnomaly creates an Anomaly detector with the given tolerance factor.
// tolerance must be greater than 0 and represents the maximum allowed
// fractional deviation from the baseline (e.g. 0.5 = 50%).
func NewAnomaly(tolerance float64) (*Anomaly, error) {
	if tolerance <= 0 {
		return nil, errors.New("anomaly: tolerance must be greater than zero")
	}
	return &Anomaly{
		tolerance: tolerance,
		baseline:  make(map[int]float64),
		anomalous: make(map[int]bool),
	}, nil
}

// SetBaseline records the expected frequency for a port.
func (a *Anomaly) SetBaseline(port int, freq float64) error {
	if port < 1 || port > 65535 {
		return errors.New("anomaly: port out of range")
	}
	if freq < 0 {
		return errors.New("anomaly: frequency must be non-negative")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.baseline[port] = freq
	return nil
}

// Observe checks whether the given observed frequency for a port is anomalous
// relative to its baseline. Returns true if the port is considered anomalous.
func (a *Anomaly) Observe(port int, observed float64) (bool, error) {
	if port < 1 || port > 65535 {
		return false, errors.New("anomaly: port out of range")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	base, ok := a.baseline[port]
	if !ok {
		a.anomalous[port] = true
		return true, nil
	}
	var deviation float64
	if base == 0 {
		deviation = observed
	} else {
		deviation = abs64(observed-base) / base
	}
	isAnomalous := deviation > a.tolerance
	a.anomalous[port] = isAnomalous
	return isAnomalous, nil
}

// IsAnomalous reports whether the port was last observed as anomalous.
func (a *Anomaly) IsAnomalous(port int) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.anomalous[port]
}

// Reset clears all anomaly state for all ports.
func (a *Anomaly) Reset() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.anomalous = make(map[int]bool)
}

func abs64(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
