package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
	"github.com/kubemq-hub/kubemq-bridges/pkg/metrics"
	"github.com/kubemq-hub/kubemq-bridges/sources"
	"github.com/kubemq-hub/kubemq-bridges/targets"
)

type Binder struct {
	name    string
	log     *logger.Logger
	sources []sources.Source
	targets []middleware.Middleware
}

func NewBinder() *Binder {
	return &Binder{}
}
func (b *Binder) buildMiddleware(target targets.Target, cfg config.BindingConfig, exporter *metrics.Exporter) (middleware.Middleware, error) {
	log, err := middleware.NewLogMiddleware(cfg.Name, cfg.Properties)
	if err != nil {
		return nil, err
	}
	retry, err := middleware.NewRetryMiddleware(cfg.Properties, b.log)
	if err != nil {
		return nil, err
	}
	rateLimiter, err := middleware.NewRateLimitMiddleware(cfg.Properties)
	if err != nil {
		return nil, err
	}
	met, err := middleware.NewMetricsMiddleware(cfg, exporter)
	if err != nil {
		return nil, err
	}
	md := middleware.Chain(target, middleware.RateLimiter(rateLimiter), middleware.Retry(retry), middleware.Metric(met), middleware.Log(log))
	return md, nil
}
func (b *Binder) Init(ctx context.Context, cfg config.BindingConfig, exporter *metrics.Exporter) error {
	b.name = cfg.Name
	b.log = logger.NewLogger(b.name)
	for _, connection := range cfg.Targets.Connections {
		target, err := targets.Init(ctx, cfg.Targets.Kind, connection)
		if err != nil {
			return fmt.Errorf("error loading targets conntector %s on binding %s, %w", cfg.Targets.Name, b.name, err)
		}
		md, err := b.buildMiddleware(target, cfg, exporter)
		if err != nil {
			return fmt.Errorf("error loading middlewares %s on binding %s, %w", cfg.Targets.Name, b.name, err)
		}
		b.targets = append(b.targets, md)
	}

	for _, connection := range cfg.Sources.Connections {
		source, err := sources.Init(ctx, cfg.Sources.Kind, connection)
		if err != nil {
			return fmt.Errorf("error loading sources conntector %s on binding %s, %w", cfg.Sources.Name, b.name, err)
		}
		b.sources = append(b.sources, source)
	}
	b.log.Infof("binding %s initialized successfully", b.name)
	return nil
}

func (b *Binder) Start(ctx context.Context) error {
	if b.targets == nil {
		return fmt.Errorf("error starting binding connector %s,no valid initialzed targets middleware found", b.name)
	}
	if b.sources == nil {
		return fmt.Errorf("error starting binding connector %s,no valid initialzed sources found", b.name)
	}

	for _, source := range b.sources {
		err := source.Start(ctx, b.targets, b.log)
		if err != nil {
			return err
		}
	}

	b.log.Infof("binding %s started successfully", b.name)
	return nil
}
func (b *Binder) Stop() error {
	for _, source := range b.sources {
		err := source.Stop()
		if err != nil {
			return err
		}
	}
	b.log.Infof("binding %s stopped successfully", b.name)
	return nil
}
