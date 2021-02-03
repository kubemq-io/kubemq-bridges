package sources

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
	"github.com/kubemq-hub/kubemq-bridges/sources/command"
	"github.com/kubemq-hub/kubemq-bridges/sources/events"
	events_store "github.com/kubemq-hub/kubemq-bridges/sources/events-store"
	"github.com/kubemq-hub/kubemq-bridges/sources/query"
	"github.com/kubemq-hub/kubemq-bridges/sources/queue"
	"github.com/kubemq-hub/kubemq-bridges/sources/queue_stream"
)

type Source interface {
	Init(ctx context.Context, connection config.Metadata, properties config.Metadata) error
	Start(ctx context.Context, target []middleware.Middleware, log *logger.Logger) error
	Stop() error
}

func Init(ctx context.Context, kind string, connection config.Metadata, properties config.Metadata) (Source, error) {
	switch kind {
	case "source.command", "kubemq.command":
		source := command.New()
		if err := source.Init(ctx, connection, properties); err != nil {
			return nil, err
		}
		return source, nil
	case "source.query", "kubemq.query":
		source := query.New()
		if err := source.Init(ctx, connection, properties); err != nil {
			return nil, err
		}
		return source, nil
	case "source.events", "kubemq.events":
		source := events.New()
		if err := source.Init(ctx, connection, properties); err != nil {
			return nil, err
		}
		return source, nil
	case "source.events-store", "kubemq.events-store":
		source := events_store.New()
		if err := source.Init(ctx, connection, properties); err != nil {
			return nil, err
		}
		return source, nil
	case "source.queue", "kubemq.queue":
		source := queue.New()
		if err := source.Init(ctx, connection, properties); err != nil {
			return nil, err
		}
		return source, nil
	case "source.queue-stream", "kubemq.queue-stream":
		source := queue_stream.New()
		if err := source.Init(ctx, connection, properties); err != nil {
			return nil, err
		}
		return source, nil
	default:
		return nil, fmt.Errorf("invalid kind %s for source", kind)
	}

}
