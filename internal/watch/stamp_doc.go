// Package watch provides port-watching primitives for portwatch.
//
// # Stamp
//
// Stamp records the most-recent time a port was observed as open.
// It is useful for tracking port activity over time and detecting
// ports that have not been seen recently.
//
// Example usage:
//
//	stamp := watch.NewStamp()
//	_ = stamp.Touch(8080)
//	if ts, ok := stamp.Last(8080); ok {
//		fmt.Println("port 8080 last seen:", ts)
//	}
package watch
