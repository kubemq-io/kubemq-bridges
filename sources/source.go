package sources

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/middleware"
	"github.com/kubemq-io/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-bridges/sources/command"
	"github.com/kubemq-io/kubemq-bridges/sources/events"
	events_store "github.com/kubemq-io/kubemq-bridges/sources/events-store"
	"github.com/kubemq-io/kubemq-bridges/sources/query"
	"github.com/kubemq-io/kubemq-bridges/sources/queue"
)

type Source interface {
	Init(ctx context.Context, connection config.Metadata, properties config.Metadata, bindingName string, log *logger.Logger) error
	Start(ctx context.Context, target []middleware.Middleware) error
	Stop() error
}

func Init(ctx context.Context, kind string, connection config.Metadata, properties config.Metadata, bindingName string, log *logger.Logger) (Source, error) {
	switch kind {
	case "source.command", "kubemq.command":
		source := command.New()
		if err := source.Init(ctx, connection, properties, bindingName, log); err != nil {
			return nil, err
		}
		return source, nil
	case "source.query", "kubemq.query":
		source := query.New()
		if err := source.Init(ctx, connection, properties, bindingName, log); err != nil {
			return nil, err
		}
		return source, nil
	case "source.events", "kubemq.events":
		source := events.New()
		if err := source.Init(ctx, connection, properties, bindingName, log); err != nil {
			return nil, err
		}
		return source, nil
	case "source.events-store", "kubemq.events-store":
		source := events_store.New()
		if err := source.Init(ctx, connection, properties, bindingName, log); err != nil {
			return nil, err
		}
		return source, nil
	case "source.queue", "kubemq.queue":
		source := queue.New()
		if err := source.Init(ctx, connection, properties, bindingName, log); err != nil {
			return nil, err
		}
		return source, nil
	default:
		return nil, fmt.Errorf("invalid kind %s for source", kind)
	}

}
