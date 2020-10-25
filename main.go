package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	connectorBridges "github.com/kubemq-hub/builder/connector/bridges"
	"github.com/kubemq-hub/kubemq-bridges/api"
	"github.com/kubemq-hub/kubemq-bridges/binding"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
	"io/ioutil"
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
	log        *logger.Logger
	build      = flag.Bool("build", false, "build sources configuration")
	configFile = flag.String("config", "config.yaml", "set config file name")
)

func loadCfgBindings() []*connectorBridges.Binding {
	file, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return nil
	}
	list := &connectorBridges.Bindings{}
	err = yaml.Unmarshal(file, list)
	if err != nil {
		return nil
	}
	return list.Bindings
}

func buildConfig() error {
	var err error
	var bindingsYaml []byte
	if bindingsYaml, err = connectorBridges.NewBridges("kubemq-bridges").
		SetBindings(loadCfgBindings()).
		Render(); err != nil {
		return err
	}
	return ioutil.WriteFile("config.yaml", bindingsYaml, 0644)
}

func run() error {
	var gracefulShutdown = make(chan os.Signal, 1)
	signal.Notify(gracefulShutdown, syscall.SIGTERM)
	signal.Notify(gracefulShutdown, syscall.SIGINT)
	signal.Notify(gracefulShutdown, syscall.SIGQUIT)
	configCh := make(chan *config.Config)
	cfg, err := config.Load(configCh)
	if err != nil {
		return err
	}
	err = cfg.Validate()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	bindingsService, err := binding.New()
	if err != nil {
		return err
	}
	err = bindingsService.Start(ctx, cfg)
	if err != nil {
		return err
	}
	apiServer, err := api.Start(ctx, cfg.ApiPort, bindingsService)
	if err != nil {
		return err
	}
	for {
		select {
		case newConfig := <-configCh:
			err = newConfig.Validate()
			if err != nil {
				return fmt.Errorf("error on validation new config file: %s", err.Error())

			}
			bindingsService.Stop()
			err = bindingsService.Start(ctx, newConfig)
			if err != nil {
				return fmt.Errorf("error on restarting service with new config file: %s", err.Error())
			}
			if apiServer != nil {
				err = apiServer.Stop()
				if err != nil {
					return fmt.Errorf("error on shutdown api server: %s", err.Error())
				}
			}

			apiServer, err = api.Start(ctx, newConfig.ApiPort, bindingsService)
			if err != nil {
				return fmt.Errorf("error on start api server: %s", err.Error())
			}
		case <-gracefulShutdown:
			_ = apiServer.Stop()
			bindingsService.Stop()
			return nil
		}
	}

}
func main() {
	log = logger.NewLogger("main")
	flag.Parse()
	if *build {
		err := buildConfig()
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
	config.SetConfigFile(*configFile)
	log.Infof("starting kubemq bridges connector version: %s, commit: %s, date %s", version, commit, date)
	if err := run(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
