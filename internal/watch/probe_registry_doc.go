// Package watch provides port monitoring primitives for portwatch.
//
// ProbeRegistry
//
// ProbeRegistry is a concurrency-safe store that tracks the most recent
// probe result for each port. It is intended to complement the Probe type,
// which performs individual TCP reachability checks, by maintaining a
// historical record of outcomes that other components can query without
// re-probing the network.
//
// Usage:
//
//	 reg, err := watch.NewProbeRegistry(2 * time.Second)
//	 if err != nil { ... }
//
//	 // After probing port 8080:
//	 _ = reg.Record(8080, true, nil)
//
//	 // Query the last result:
//	 rec := reg.Get(8080)
//	 if rec != nil && rec.Open {
//	     fmt.Println("port 8080 was open at", rec.LastCheck)
//	 }
//
// Records are replaced on every call to Record; the registry does not
// retain history beyond the most recent check per port.
package watch
