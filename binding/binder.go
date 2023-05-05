package binding

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/middleware"
	"github.com/kubemq-io/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-bridges/pkg/metrics"
	"github.com/kubemq-io/kubemq-bridges/sources"
	"github.com/kubemq-io/kubemq-bridges/targets"
)

type Binder struct {
	name              string
	log               *logger.Logger
	sources           []sources.Source
	targetsMiddleware []middleware.Middleware
	targets           []targets.Target
}

func NewBinder() *Binder {
	return &Binder{}
}
func (b *Binder) buildMiddleware(target targets.Target, cfg config.BindingConfig, exporter *metrics.Exporter, log *middleware.LogMiddleware) (middleware.Middleware, error) {

	retry, err := middleware.NewRetryMiddleware(cfg.Properties, b.log)
	if err != nil {
		return nil, err
	}
	rateLimiter, err := middleware.NewRateLimitMiddleware(cfg.Properties)
	if err != nil {
		return nil, err
	}
	var md middleware.Middleware
	if exporter != nil {
		met, err := middleware.NewMetricsMiddleware(cfg, exporter)
		if err != nil {
			return nil, err
		}
		md = middleware.Chain(target, middleware.RateLimiter(rateLimiter), middleware.Retry(retry), middleware.Metric(met), middleware.Log(log))
	} else {
		md = middleware.Chain(target, middleware.RateLimiter(rateLimiter), middleware.Retry(retry), middleware.Log(log))
	}

	return md, nil
}
func (b *Binder) Init(ctx context.Context, cfg config.BindingConfig, exporter *metrics.Exporter) error {
	b.name = cfg.Name
	log, err := middleware.NewLogMiddleware(cfg.Name, cfg.Properties)
	if err != nil {
		return err
	}
	b.log = log.Logger
	for _, connection := range cfg.Targets.Connections {
		target, err := targets.Init(ctx, cfg.Targets.Kind, connection, cfg.Name, b.log)
		if err != nil {
			return fmt.Errorf("error loading targets conntector on binding %s, %w", b.name, err)
		}
		md, err := b.buildMiddleware(target, cfg, exporter, log)
		if err != nil {
			return fmt.Errorf("error loading middlewares on binding %s, %w", b.name, err)
		}
		b.targetsMiddleware = append(b.targetsMiddleware, md)
		b.targets = append(b.targets, target)
	}

	for _, connection := range cfg.Sources.Connections {
		source, err := sources.Init(ctx, cfg.Sources.Kind, connection, cfg.Properties, cfg.Name, b.log)
		if err != nil {
			return fmt.Errorf("error loading sources conntector on binding %s, %w", b.name, err)
		}
		b.sources = append(b.sources, source)
	}
	b.log.Infof("binding %s initialized successfully", b.name)
	return nil
}

func (b *Binder) Start(ctx context.Context) error {
	if b.targetsMiddleware == nil {
		return fmt.Errorf("error starting binding connector %s,no valid initialzed targets middleware found", b.name)
	}
	if b.sources == nil {
		return fmt.Errorf("error starting binding connector %s,no valid initialzed sources found", b.name)
	}

	for _, source := range b.sources {
		err := source.Start(ctx, b.targetsMiddleware)
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
	for _, target := range b.targets {
		err := target.Stop()
		if err != nil {
			return err
		}
	}
	b.log.Infof("binding %s stopped successfully", b.name)
	return nil
}
