// Package watch provides primitives for monitoring port state changes.
//
// Shadow
//
// Shadow maintains a secondary read-only copy of the last known port
// state. It is useful when you need to compare the current scan result
// against a previous baseline without touching the primary snapshot.
//
// Usage:
//
//	s := watch.NewShadow()
//	s.Update([]int{80, 443, 8080})
//
//	if s.Contains(80) {
//		// port was open in the previous scan
//	}
package watch
