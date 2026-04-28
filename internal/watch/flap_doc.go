// Package watch provides building blocks for monitoring open ports and
// reacting to changes in network state.
//
// # Flap
//
// Flap detects ports that rapidly alternate between open and closed states,
// a condition commonly called "port flapping". This can indicate an unstable
// service, a misconfigured firewall, or an active network attack.
//
// Usage:
//
//	f, err := watch.NewFlap(30*time.Second, 4)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Record each observed state change:
//	f.Record(port, time.Now())
//
//	// Check whether a port is flapping:
//	if f.IsFlapping(port, time.Now()) {
//		log.Printf("port %d is flapping", port)
//	}
package watch
