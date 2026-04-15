// Package watch provides the core port-watching loop and supporting
// utilities used by the portwatch daemon.
//
// # Port Filtering
//
// PortFilter is the central interface for deciding which ports are
// included in a scan cycle.  Three implementations are provided:
//
//   - RangeFilter  – allows only ports within an inclusive [min, max] band.
//   - ExcludeFilter – blocks a fixed set of port numbers (e.g. well-known
//     system ports that should always be ignored).
//   - ChainFilter  – composes multiple PortFilter values; a port is allowed
//     only when every filter in the chain allows it.
//
// Filters are designed to be constructed once at startup from the loaded
// configuration and then used concurrently without locking.
package watch
