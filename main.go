package main

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/api"
	"github.com/kubemq-hub/kubemq-bridges/binding"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

var (
	log *logger.Logger
)

func run() error {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	err = cfg.Validate()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bindingsService := binding.New()
	err = bindingsService.Start(ctx, cfg)
	if err != nil {
		fmt.Println(err)
		//return err
	}
	apiServer, err := api.Start(ctx, cfg.ApiPort, bindingsService)
	if err != nil {
		return err
	}
	<-gracefulShutdown
	apiServer.Stop()
	bindingsService.Stop()

	return nil
}
func main() {
	log = logger.NewLogger("main")
	log.Infof("starting kubemq bridges connector version: %s, commit: %s, date %s", version, commit, date)
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
