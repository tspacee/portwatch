// Package watch provides primitives for monitoring, classifying, and
// reacting to changes in observed open ports.
//
// Band
//
// Band groups port numbers into named frequency ranges and counts how many
// ports from each scan fall within each range. It is useful for high-level
// reporting such as distinguishing system ports (1–1023), registered ports
// (1024–49151), and dynamic/ephemeral ports (49152–65535).
//
// Example usage:
//
//	b := watch.NewBand()
//	_ = b.Define("system",     1,     1023)
//	_ = b.Define("registered", 1024, 49151)
//	b.Classify(openPorts)
//	fmt.Println(b.Snapshot())
package watch
