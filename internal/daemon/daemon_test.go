package daemon_test

import (
	"bytes"
	"context"
	"testing"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/snapshot"
)

func buildDaemon(t *testing.T, cfg *config.Config) (*daemon.Daemon, *bytes.Buffer) {
	t.Helper()

	engine, err := rules.NewEngine(nil)
	if err != nil {
		t.Fatalf("NewEngine: %v", err)
	}

	buf := &bytes.Buffer{}
	notifier := alert.NewLogNotifier(buf)
	dispatcher := alert.NewDispatcher([]alert.Notifier{notifier})

	dir := t.TempDir()
	manager := snapshot.NewManager(dir + "/snap.json")

	d, err := daemon.New(cfg, engine, dispatcher, manager)
	if err != nil {
		t.Fatalf("daemon.New: %v", err)
	}
	return d, buf
}

func TestNew_ValidConfig(t *testing.T) {
	cfg := config.Default()
	d, _ := buildDaemon(t, cfg)
	if d == nil {
		t.Fatal("expected non-nil daemon")
	}
}

func TestNew_InvalidPortRange(t *testing.T) {
	cfg := config.Default()
	cfg.PortRange = "invalid"

	engine, _ := rules.NewEngine(nil)
	buf := &bytes.Buffer{}
	notifier := alert.NewLogNotifier(buf)
	dispatcher := alert.NewDispatcher([]alert.Notifier{notifier})
	manager := snapshot.NewManager(t.TempDir() + "/snap.json")

	_, err := daemon.New(cfg, engine, dispatcher, manager)
	if err == nil {
		t.Fatal("expected error for invalid port range")
	}
}

func TestRun_StopsOnContextCancel(t *testing.T) {
	cfg := config.Default()
	cfg.Interval = 500 * time.Millisecond

	d, _ := buildDaemon(t, cfg)

	ctx, cancel := context.WithTimeout(context.Background(), 600*time.Millisecond)
	defer cancel()

	done := make(chan struct{})
	go func() {
		d.Run(ctx)
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("daemon did not stop after context cancellation")
	}
}
