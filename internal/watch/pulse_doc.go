// Package watch provides primitives for observing and controlling port
// scanning behaviour.
//
// # Pulse
//
// Pulse tracks the heartbeat interval between successive scan cycles.
// Call Beat on each cycle; after two or more beats the Summary method
// returns the count, average, minimum, and maximum intervals observed.
//
// Example:
//
//	pulse := watch.NewPulse()
//
//	for {
//		pulse.Beat()
//		runScan()
//	}
//
//	s, err := pulse.Summary()
//	if err == nil {
//		log.Printf("avg scan interval: %v", s.Avg)
//	}
package watch
