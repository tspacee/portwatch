package watch_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/metrics"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/snapshot"
	"github.com/user/portwatch/internal/watch"
)

func buildWatcher(t *testing.T) *watch.Watcher {
	t.Helper()
	scanner, err := portscanner.NewScanner(portscanner.Config{StartPort: 1, EndPort: 1024})
	if err != nil {
		t.Fatalf("scanner: %v", err)
	}
	engine, err := rules.NewEngine(nil)
	if err != nil {
		t.Fatalf("engine: %v", err)
	}
	mgr := snapshot.NewManager("")
	notifier := alert.NewLogNotifier(bytes.NewBuffer(nil))
	disp := alert.NewDispatcher(notifier)
	w, err := watch.New(watch.Config{
		Scanner:    scanner,
		Engine:     engine,
		Manager:    mgr,
		Dispatcher: disp,
		Interval:   50 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("watcher: %v", err)
	}
	return w
}

func buildWatcherWithCollector(t *testing.T, col *metrics.Collector) *watch.Watcher {
	t.Helper()
	scanner, err := portscanner.NewScanner(portscanner.Config{StartPort: 1, EndPort: 100})
	if err != nil {
		t.Fatalf("scanner: %v", err)
	}
	engine, err := rules.NewEngine(nil)
	if err != nil {
		t.Fatalf("engine: %v", err)
	}
	mgr := snapshot.NewManager("")
	notifier := alert.NewLogNotifier(bytes.NewBuffer(nil))
	disp := alert.NewDispatcher(notifier)
	w, err := watch.New(watch.Config{
		Scanner:    scanner,
		Engine:     engine,
		Manager:    mgr,
		Dispatcher: disp,
		Collector:  col,
		Interval:   30 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("watcher with collector: %v", err)
	}
	return w
}

func TestNew_MissingScanner(t *testing.T) {
	_, err := watch.New(watch.Config{})
	if err == nil {
		t.Fatal("expected error for nil scanner")
	}
}

func TestNew_MissingEngine(t *testing.T) {
	scanner, _ := portscanner.NewScanner(portscanner.Config{StartPort: 1, EndPort: 100})
	_, err := watch.New(watch.Config{Scanner: scanner})
	if err == nil {
		t.Fatal("expected error for nil engine")
	}
}

func TestNew_DefaultInterval(t *testing.T) {
	w := buildWatcher(t)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	w := buildWatcher(t)
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- w.Run(ctx) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	select {
	case err := <-done:
		if err != context.Canceled {
			t.Fatalf("expected context.Canceled, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("watcher did not stop after context cancel")
	}
}

func TestRun_RecordsMetrics(t *testing.T) {
	col := metrics.NewCollector()
	w := buildWatcherWithCollector(t, col)
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	w.Run(ctx) //nolint:errcheck
	if col.Summary().TotalScans == 0 {
		t.Error("expected at least one recorded scan")
	}
}
