// Package watch provides utilities for monitoring port activity.
//
// # Decay
//
// Decay tracks per-port floating-point scores that diminish over time.
// It is useful for modelling recency-weighted port activity: a port that
// was seen frequently in the recent past has a high score, but if it
// goes quiet the score gradually falls toward zero.
//
// Usage:
//
//	d, err := watch.NewDecay(0.2) // 20 % loss per second
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Observe activity on port 8080.
//	_ = d.Add(8080, 1.0)
//
//	// Call Tick periodically (e.g. from the watcher loop) to apply decay.
//	d.Tick()
//
//	// Query the current score.
//	fmt.Println(d.Score(8080))
package watch
