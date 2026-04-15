package daemon

import (
	"context"
	"log"
	"time"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/portscanner"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/snapshot"
)

// Daemon orchestrates periodic port scanning, snapshot diffing, rule evaluation, and alerting.
type Daemon struct {
	cfg        *config.Config
	scanner    *portscanner.Scanner
	engine     *rules.Engine
	dispatcher *alert.Dispatcher
	manager    *snapshot.Manager
}

// New creates a Daemon from the provided configuration.
func New(cfg *config.Config, engine *rules.Engine, dispatcher *alert.Dispatcher, manager *snapshot.Manager) (*Daemon, error) {
	scanner, err := portscanner.NewScanner(cfg.PortRange, cfg.Protocol)
	if err != nil {
		return nil, err
	}
	return &Daemon{
		cfg:        cfg,
		scanner:    scanner,
		engine:     engine,
		dispatcher: dispatcher,
		manager:    manager,
	}, nil
}

// Run starts the daemon loop, ticking at the configured interval until ctx is cancelled.
func (d *Daemon) Run(ctx context.Context) {
	ticker := time.NewTicker(d.cfg.Interval)
	defer ticker.Stop()

	log.Printf("portwatch daemon started (interval=%s)", d.cfg.Interval)

	// Run once immediately before waiting for the first tick.
	d.tick()

	for {
		select {
		case <-ticker.C:
			d.tick()
		case <-ctx.Done():
			log.Println("portwatch daemon stopped")
			return
		}
	}
}

// tick performs a single scan-diff-evaluate-alert cycle.
func (d *Daemon) tick() {
	ports, err := d.scanner.OpenPorts()
	if err != nil {
		log.Printf("scan error: %v", err)
		return
	}

	current := snapshot.New(ports)
	prev, err := d.manager.Load()
	if err != nil {
		log.Printf("snapshot load error: %v", err)
	}

	diff := snapshot.Diff(prev, current)

	violations := d.engine.Evaluate(diff)
	for _, v := range violations {
		if err := d.dispatcher.Dispatch(v); err != nil {
			log.Printf("dispatch error: %v", err)
		}
	}

	if err := d.manager.Save(current); err != nil {
		log.Printf("snapshot save error: %v", err)
	}
}
