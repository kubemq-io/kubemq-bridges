package targets

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-bridges/targets/command"
	"github.com/kubemq-io/kubemq-bridges/targets/events"
	events_store "github.com/kubemq-io/kubemq-bridges/targets/events-store"
	"github.com/kubemq-io/kubemq-bridges/targets/query"
	"github.com/kubemq-io/kubemq-bridges/targets/queue"
)

type Target interface {
	Init(ctx context.Context, connection config.Metadata, log *logger.Logger) error
	Do(ctx context.Context, request interface{}) (interface{}, error)
	Stop() error
}

func Init(ctx context.Context, kind string, connection config.Metadata, log *logger.Logger) (Target, error) {

	switch kind {
	case "target.command", "kubemq.command":
		target := command.New()
		if err := target.Init(ctx, connection, log); err != nil {
			return nil, err
		}
		return target, nil
	case "target.query", "kubemq.query":
		target := query.New()
		if err := target.Init(ctx, connection, log); err != nil {
			return nil, err
		}
		return target, nil
	case "target.events", "kubemq.events":
		target := events.New()
		if err := target.Init(ctx, connection, log); err != nil {
			return nil, err
		}
		return target, nil
	case "target.events-store", "kubemq.events-store":
		target := events_store.New()
		if err := target.Init(ctx, connection, log); err != nil {
			return nil, err
		}
		return target, nil
	case "target.queue", "kubemq.queue":
		target := queue.New()
		if err := target.Init(ctx, connection, log); err != nil {
			return nil, err
		}
		return target, nil
	default:
		return nil, fmt.Errorf("invalid kind %s for target", kind)
	}

}
