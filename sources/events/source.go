package events

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
	properties config.Metadata
}

func New() *Source {
	return &Source{}

}
func (s *Source) Init(ctx context.Context, Source config.Metadata,properties config.Metadata) error {
	var err error
	s.opts, err = parseOptions(Source)
	if err != nil {
		return err
	}
	s.properties=properties
	s.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(s.opts.host, s.opts.port),
		kubemq.WithClientId(s.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(s.opts.authToken),
		kubemq.WithCheckConnection(true),
		kubemq.WithMaxReconnects(s.opts.maxReconnects),
		kubemq.WithAutoReconnect(s.opts.autoReconnect),
		kubemq.WithReconnectInterval(s.opts.reconnectIntervalSeconds))
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) Start(ctx context.Context, targets []middleware.Middleware, log *logger.Logger) error {
	s.targets = targets
	s.log = log
	group := nuid.Next()
	if s.opts.group != "" {
		group = s.opts.group
	}
	errCh := make(chan error, 1)
	eventsCh, err := s.client.SubscribeToEvents(ctx, s.opts.channel, group, errCh)
	if err != nil {
		return fmt.Errorf("error on subscribing to events channel, %w", err)
	}
	go func(ctx context.Context, eventsCh <-chan *kubemq.Event, errCh chan error) {
		s.run(ctx, eventsCh, errCh)
	}(ctx, eventsCh, errCh)
	return nil
}

func (s *Source) run(ctx context.Context, eventsCh <-chan *kubemq.Event, errCh chan error) {
	for {
		select {
		case event := <-eventsCh:
			for _, target := range s.targets {
				go func(event *kubemq.Event, target middleware.Middleware) {
					_, err := target.Do(ctx, event)
					if err != nil {
						s.log.Errorf("error received from target, %s", err.Error())
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
