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
	Init(ctx context.Context, connection config.Metadata) error
	Do(ctx context.Context, request interface{}) (interface{}, error)
}

func Init(ctx context.Context, kind string, connection config.Metadata) (Target, error) {

	switch kind {
	case "target.command":
		target := command.New()
		if err := target.Init(ctx, connection); err != nil {
			return nil, err
		}
		return target, nil
	case "target.query":
		target := query.New()
		if err := target.Init(ctx, connection); err != nil {
			return nil, err
		}
		return target, nil
	case "target.events":
		target := events.New()
		if err := target.Init(ctx, connection); err != nil {
			return nil, err
		}
		return target, nil
	case "target.events-store":
		target := events_store.New()
		if err := target.Init(ctx, connection); err != nil {
			return nil, err
		}
		return target, nil
	case "target.queue":
		target := queue.New()
		if err := target.Init(ctx, connection); err != nil {
			return nil, err
		}
		return target, nil
	default:
		return nil, fmt.Errorf("invalid kind %s for target", kind)
	}

}
