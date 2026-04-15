// Package metrics provides lightweight, in-process scan statistics for
// portwatch. A Collector accumulates ScanResult values produced by the
// daemon on every scan cycle. A Reporter can print a human-readable
// summary of those statistics to any io.Writer.
//
// Typical usage:
//
//	col := metrics.NewCollector()
//
//	// inside the scan loop:
//	col.Record(metrics.ScanResult{
//		Timestamp:    time.Now(),
//		OpenPorts:    len(ports),
//		Violations:   len(violations),
//		ScanDuration: elapsed,
//	})
//
//	// on SIGUSR1 or shutdown:
//	metrics.NewReporter(col, os.Stdout).Print()
package metrics
