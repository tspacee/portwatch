// Package watch provides port watching primitives for portwatch.
//
// # Resolver
//
// Resolver maps port numbers to human-readable service names.
// It ships with a built-in table of well-known ports and supports
// runtime registration of custom port-to-service mappings.
//
// Example:
//
//	r := watch.NewResolver()
//	_ = r.Register(9090, "prometheus")
//	fmt.Println(r.Resolve(443))  // "https"
//	fmt.Println(r.Resolve(9090)) // "prometheus"
//	fmt.Println(r.Resolve(9999)) // "unknown"
package watch
