// Package watch implements the core port-watching loop for portwatch.
//
// A Watcher ties together the port scanner, snapshot manager, rules engine
// and alert dispatcher into a single periodic scan cycle:
//
//  1. Scan open ports via portscanner.Scanner.
//  2. Diff the result against the previous snapshot.
//  3. Evaluate the diff against configured rules.
//  4. Dispatch alerts for any violations.
//  5. Persist the current snapshot for the next cycle.
//
// Metrics are recorded on every tick when a metrics.Collector is provided.
//
// # Usage
//
//	w, err := watch.New(watch.Config{
//	    Scanner:    scanner,
//	    Engine:     engine,
//	    Manager:    manager,
//	    Dispatcher: dispatcher,
//	    Collector:  collector,  // optional
//	    Interval:   30 * time.Second,
//	})
//	if err != nil { ... }
//	if err := w.Run(ctx); err != nil && err != context.Canceled { ... }
package watch
