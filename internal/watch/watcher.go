// Package watch provides a port change watcher that compares
// current scan results against a baseline snapshot and emits
// violations when unexpected changes are detected.
package watch

import (
	"context"
	"fmt"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/snapshot"
)

// Watcher periodically scans ports, diffs against the last snapshot,
// evaluates rules and dispatches alerts on violations.
type Watcher struct {
	scanner   *portscanner.Scanner
	engine    *rules.Engine
	manager   *snapshot.Manager
	dispatcher *alert.Dispatcher
	collector  *metrics.Collector
	interval  time.Duration
}

// Config holds the parameters needed to build a Watcher.
type Config struct {
	Scanner    *portscanner.Scanner
	Engine     *rules.Engine
	Manager    *snapshot.Manager
	Dispatcher *alert.Dispatcher
	Collector  *metrics.Collector
	Interval   time.Duration
}

// New validates the config and returns a ready Watcher.
func New(cfg Config) (*Watcher, error) {
	if cfg.Scanner == nil {
		return nil, fmt.Errorf("watcher: scanner is required")
	}
	if cfg.Engine == nil {
		return nil, fmt.Errorf("watcher: engine is required")
	}
	if cfg.Manager == nil {
		return nil, fmt.Errorf("watcher: snapshot manager is required")
	}
	if cfg.Dispatcher == nil {
		return nil, fmt.Errorf("watcher: dispatcher is required")
	}
	if cfg.Interval <= 0 {
		cfg.Interval = 30 * time.Second
	}
	return &Watcher{
		scanner:    cfg.Scanner,
		engine:     cfg.Engine,
		manager:    cfg.Manager,
		dispatcher: cfg.Dispatcher,
		collector:  cfg.Collector,
		interval:   cfg.Interval,
	}, nil
}

// Run starts the watch loop and blocks until ctx is cancelled.
func (w *Watcher) Run(ctx context.Context) error {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := w.tick(ctx); err != nil {
				return err
			}
		}
	}
}

func (w *Watcher) tick(ctx context.Context) error {
	start := time.Now()
	ports, err := w.scanner.OpenPorts(ctx)
	duration := time.Since(start)
	if w.collector != nil {
		if err != nil {
			w.collector.RecordError(duration)
		} else {
			w.collector.Record(duration, len(ports))
		}
	}
	if err != nil {
		return fmt.Errorf("watcher: scan failed: %w", err)
	}
	current := snapshot.New(ports)
	prev := w.manager.Latest()
	if prev != nil {
		diff := snapshot.Diff(prev, current)
		violations := w.engine.Evaluate(diff)
		for _, v := range violations {
			w.dispatcher.Dispatch(ctx, v)
		}
	}
	w.manager.Save(current)
	return nil
}
