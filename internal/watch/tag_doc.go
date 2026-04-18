// Package watch provides port watching primitives for portwatch.
//
// # Tag Registry
//
// TagRegistry provides a thread-safe mapping from port numbers to
// human-readable labels. Labels are useful for enriching alert output
// and status reports with service names (e.g. port 80 → "http").
//
// Example usage:
//
//	reg := watch.NewTagRegistry()
//	_ = reg.Set(80, "http")
//	_ = reg.Set(443, "https")
//
//	if label, ok := reg.Get(port); ok {
//		fmt.Printf("port %d is %s\n", port, label)
//	}
package watch
