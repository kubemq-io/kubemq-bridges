package targets

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/targets/command"
	"github.com/kubemq-hub/kubemq-bridges/targets/events"
	events_store "github.com/kubemq-hub/kubemq-bridges/targets/events-store"
	"github.com/kubemq-hub/kubemq-bridges/targets/query"
	"github.com/kubemq-hub/kubemq-bridges/targets/queue"

)

type Target interface {
	Init(ctx context.Context, cfg config.Metadata) error
	Do(ctx context.Context, request interface{}) (interface{}, error)
	Name() string
}

func Init(ctx context.Context, cfg config.Metadata) (Target, error) {

	switch cfg.Kind {
	case "target.command":
		target := command.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.query":
		target := query.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.events":
		target := events.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.events-store":
		target := events_store.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.queue":
		target := queue.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	default:
		return nil, fmt.Errorf("invalid kind %s for target %s", cfg.Kind, cfg.Name)
	}

}
