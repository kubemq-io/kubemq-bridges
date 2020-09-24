package events_store

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"

	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
)

type Source struct {
	opts    options
	client  *kubemq.Client
	log     *logger.Logger
	targets []middleware.Middleware
}

func New() *Source {
	return &Source{}

}
func (s *Source) Init(ctx context.Context, connection config.Metadata) error {
	var err error
	s.opts, err = parseOptions(connection)
	if err != nil {
		return err
	}
	s.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(s.opts.host, s.opts.port),
		kubemq.WithClientId(s.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(s.opts.authToken),
		kubemq.WithCheckConnection(false),
		kubemq.WithMaxReconnects(s.opts.maxReconnects),
		kubemq.WithAutoReconnect(s.opts.autoReconnect),
		kubemq.WithReconnectInterval(s.opts.reconnectIntervalSeconds))
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) Start(ctx context.Context, targets []middleware.Middleware, log *logger.Logger) error {
	s.log = log
	s.targets = targets
	group := nuid.Next()
	if s.opts.group != "" {
		group = s.opts.group
	}
	errCh := make(chan error, 1)
	eventsCh, err := s.client.SubscribeToEventsStore(ctx, s.opts.channel, group, errCh, kubemq.StartFromNewEvents())
	if err != nil {
		return fmt.Errorf("error on subscribing to events store channel, %w", err)
	}
	go func(ctx context.Context, eventsCh <-chan *kubemq.EventStoreReceive, errCh chan error) {
		s.run(ctx, eventsCh, errCh)
	}(ctx, eventsCh, errCh)

	return nil
}

func (s *Source) run(ctx context.Context, eventsCh <-chan *kubemq.EventStoreReceive, errCh chan error) {
	for {
		select {
		case event := <-eventsCh:
			for _, target := range s.targets {

				go func(event *kubemq.EventStoreReceive, target middleware.Middleware) {
					_, err := target.Do(ctx, event)
					if err != nil {
						s.log.Errorf("error received from target, %w", err)
					}

				}(event, target)
			}
		case err := <-errCh:
			s.log.Errorf("error received from kuebmq server, %s", err.Error())
			return
		case <-ctx.Done():
			return

		}
	}
}

func (s *Source) Stop() error {
	return s.client.Close()
}
