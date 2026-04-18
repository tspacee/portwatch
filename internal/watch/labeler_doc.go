// Package watch provides port watching utilities for portwatch.
//
// # Labeler
//
// Labeler maps port numbers to human-readable labels such as "http" or "ssh".
// It is safe for concurrent use.
//
// Usage:
//
//	l := watch.NewLabeler()
//	_ = l.Register(80, "http")
//	_ = l.Register(443, "https")
//	fmt.Println(l.Label(80))   // "http"
//	fmt.Println(l.Label(9000)) // "port/9000"
package watch
