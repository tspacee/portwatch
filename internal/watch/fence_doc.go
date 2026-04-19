// Package watch provides port watching primitives for portwatch.
//
// # Fence
//
// Fence maintains a set of explicitly allowed ports and reports whether
// a given port falls within that set. It is useful for enforcing strict
// allowlists independent of rule-based evaluation.
//
// Example:
//
//	f, err := watch.NewFence([]int{80, 443, 8080})
//	if err != nil {
//		log.Fatal(err)
//	}
//	if !f.Allow(port) {
//		fmt.Printf("port %d is outside the fence\n", port)
//	}
package watch
