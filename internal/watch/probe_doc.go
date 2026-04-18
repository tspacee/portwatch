// Package watch provides utilities for monitoring open ports.
//
// # Probe
//
// Probe performs lightweight TCP connectivity checks against individual ports.
// It is useful for confirming that a port detected as open by the scanner is
// genuinely accepting connections, reducing false positives caused by ephemeral
// socket states.
//
// Usage:
//
//	p, err := watch.NewProbe(500 * time.Millisecond)
//	if err != nil { ... }
//
//	if p.Check(8080) {
//		fmt.Println("port 8080 is reachable")
//	}
//
//	open := p.CheckAll([]int{80, 443, 8080})
package watch
