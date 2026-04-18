// Package watch provides adaptive scan scheduling utilities.
//
// # Sampler
//
// Sampler tracks how frequently port scans are recorded within a
// configurable sliding time window. It is intended for adaptive
// interval tuning: callers can query Rate() or Count() to decide
// whether to slow down or speed up scanning.
//
// Example:
//
//	s, _ := watch.NewSampler(30 * time.Second)
//	s.Record()           // called each time a scan completes
//	fmt.Println(s.Rate()) // scans per second over the last 30 s
package watch
