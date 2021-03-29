package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/kubemq-hub/kubemq-bridges/api"
	"github.com/kubemq-hub/kubemq-bridges/binding"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/browser"
	"github.com/kubemq-hub/kubemq-bridges/pkg/builder"
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
	build      = flag.Bool("build", false, "build bridges configuration")
	buildUrl   = flag.String("get", "", "get config file from url")
	configFile = flag.String("config", "config.yaml", "set config file name")
)

func downloadUrl() error {
	c, err := builder.GetBuildManifest(*buildUrl)
	if err != nil {
		return err
	}
	cfg := &config.Config{}
	err = yaml.Unmarshal([]byte(c.Spec.Config), &cfg)
	if err != nil {
		return err
	}
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("config.yaml", data, 0644)
	if err != nil {
		return err
	}
	return nil
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
		err := browser.OpenURL("https://build.kubemq.io/#/bridges")
		if err != nil {
			log.Error(err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
	if *buildUrl != "" {
		err := downloadUrl()
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
