package targets

import (
	"context"
	"fmt"

	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/targets/aws/sqs"
	"github.com/kubemq-hub/kubemq-bridges/targets/cache/memcached"
	"github.com/kubemq-hub/kubemq-bridges/targets/cache/redis"
	"github.com/kubemq-hub/kubemq-bridges/targets/gcp/bigquery"
	"github.com/kubemq-hub/kubemq-bridges/targets/gcp/bigtable"
	"github.com/kubemq-hub/kubemq-bridges/targets/gcp/firestore"
	"github.com/kubemq-hub/kubemq-bridges/targets/gcp/pubsub"
	"github.com/kubemq-hub/kubemq-bridges/targets/gcp/spanner"
	"github.com/kubemq-hub/kubemq-bridges/targets/http"
	"github.com/kubemq-hub/kubemq-bridges/targets/kubemq/command"
	"github.com/kubemq-hub/kubemq-bridges/targets/kubemq/events"
	events_store "github.com/kubemq-hub/kubemq-bridges/targets/kubemq/events-store"
	"github.com/kubemq-hub/kubemq-bridges/targets/kubemq/query"
	"github.com/kubemq-hub/kubemq-bridges/targets/kubemq/queue"
	"github.com/kubemq-hub/kubemq-bridges/targets/logs/elastic"
	"github.com/kubemq-hub/kubemq-bridges/targets/messaging/activemq"
	"github.com/kubemq-hub/kubemq-bridges/targets/messaging/kafka"
	"github.com/kubemq-hub/kubemq-bridges/targets/messaging/mqtt"
	"github.com/kubemq-hub/kubemq-bridges/targets/messaging/rabbitmq"
	"github.com/kubemq-hub/kubemq-bridges/targets/serverless/openfass"
	"github.com/kubemq-hub/kubemq-bridges/targets/storage/minio"
	"github.com/kubemq-hub/kubemq-bridges/targets/stores/cassandra"
	"github.com/kubemq-hub/kubemq-bridges/targets/stores/couchbase"
	"github.com/kubemq-hub/kubemq-bridges/targets/stores/mongodb"
	"github.com/kubemq-hub/kubemq-bridges/targets/stores/mssql"
	"github.com/kubemq-hub/kubemq-bridges/targets/stores/mysql"
	"github.com/kubemq-hub/kubemq-bridges/targets/stores/postgres"
	"github.com/kubemq-hub/kubemq-bridges/types"
)

var (
	errTargetNotImplemented = fmt.Errorf("target not implemented")
)

type Target interface {
	Init(ctx context.Context, cfg config.Metadata) error
	Do(ctx context.Context, request *types.Request) (*types.Response, error)
	Name() string
}

func Init(ctx context.Context, cfg config.Metadata) (Target, error) {

	switch cfg.Kind {
	case "target.aws.sqs":
		target := sqs.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.cache.redis":
		target := redis.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.cache.memcached":
		target := memcached.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.bigquery":
		target := bigquery.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.bigtable":
		target := bigtable.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.firestore":
		target := firestore.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.pubsub":
		target := pubsub.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.gcp.spanner":
		target := spanner.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.http":
		target := http.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.command":
		target := command.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.query":
		target := query.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.events":
		target := events.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.events-store":
		target := events_store.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.kubemq.queue":
		target := queue.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.logs.elastic":
		target := elastic.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.messaging.activemq":
		target := activemq.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.messaging.kafka":
		target := kafka.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.messaging.mqtt":
		target := mqtt.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.messaging.rabbitmq":
		target := rabbitmq.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.cassandra":
		target := cassandra.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.couchbase":
		target := couchbase.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.mongodb":
		target := mongodb.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.mssql":
		target := mssql.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.mysql":
		target := mysql.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.stores.postgres":
		target := postgres.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.serverless.openfaas":
		target := openfass.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	case "target.storage.minio":
		target := minio.New()
		if err := target.Init(ctx, cfg); err != nil {
			return nil, err
		}
		return target, nil
	default:
		return nil, fmt.Errorf("invalid kind %s for target %s", cfg.Kind, cfg.Name)
	}

}
