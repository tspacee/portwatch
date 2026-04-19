// Package watch provides port-watching primitives used by the portwatch daemon.
//
// # Envelope
//
// Envelope wraps a list of scanned ports with metadata:
//
//   - Seq: monotonically increasing scan sequence number
//   - Source: identifier for the scanner or host that produced the result
//   - ScannedAt: wall-clock time when the scan completed
//   - Ports: snapshot of open ports at scan time
//
// EnvelopeBuilder is the primary way to create Envelopes. It maintains an
// internal counter so consumers can detect missed or out-of-order scans.
//
// Example:
//
//	b, _ := watch.NewEnvelopeBuilder("localhost")
//	env := b.Wrap([]int{80, 443})
//	fmt.Println(env.Seq, env.Ports)
package watch
