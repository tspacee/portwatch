package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/daemon"
	"github.com/user/portwatch/internal/rules"
	"github.com/user/portwatch/internal/snapshot"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	cfg, err := config.LoadOrDefault(*cfgPath)
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	engine, err := rules.NewEngine(cfg.Rules)
	if err != nil {
		log.Fatalf("rules engine: %v", err)
	}

	notifier := alert.NewLogNotifier(os.Stdout)
	dispatcher := alert.NewDispatcher([]alert.Notifier{notifier})

	manager := snapshot.NewManager(cfg.SnapshotPath)

	d, err := daemon.New(cfg, engine, dispatcher, manager)
	if err != nil {
		log.Fatalf("daemon init: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	d.Run(ctx)
}
