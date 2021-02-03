package events

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
	"github.com/kubemq-hub/kubemq-bridges/pkg/roundrobin"

	"github.com/kubemq-hub/kubemq-bridges/pkg/uuid"
	"github.com/kubemq-io/kubemq-go"
)

type Source struct {
	opts              options
	log               *logger.Logger
	clients           []*kubemq.Client
	targets           []middleware.Middleware
	properties        config.Metadata
	roundRobin        *roundrobin.RoundRobin
	loadBalancingMode bool
}

func New() *Source {
	return &Source{}

}
func (s *Source) Init(ctx context.Context, Source config.Metadata, properties config.Metadata) error {
	var err error
	s.opts, err = parseOptions(Source)
	if err != nil {
		return err
	}
	s.properties = properties
	for i := 0; i < s.opts.sources; i++ {
		clientId := s.opts.clientId
		if s.opts.sources > 1 {
			clientId = fmt.Sprintf("%s-%d", clientId, i)
		}
		client, err := kubemq.NewClient(ctx,
			kubemq.WithAddress(s.opts.host, s.opts.port),
			kubemq.WithClientId(clientId),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC),
			kubemq.WithCheckConnection(true),
			kubemq.WithAuthToken(s.opts.authToken),
			kubemq.WithMaxReconnects(s.opts.maxReconnects),
			kubemq.WithAutoReconnect(s.opts.autoReconnect),
			kubemq.WithReconnectInterval(s.opts.reconnectIntervalSeconds))
		if err != nil {
			return err
		}
		s.clients = append(s.clients, client)
	}
	return nil
}

func (s *Source) Start(ctx context.Context, targets []middleware.Middleware, log *logger.Logger) error {
	s.roundRobin = roundrobin.NewRoundRobin(len(targets))
	if s.properties != nil {
		mode, ok := s.properties["load-balancing"]
		if ok && mode == "true" {
			s.loadBalancingMode = true
		}
	}
	s.targets = targets
	s.log = log
	if s.opts.sources > 1 && s.opts.group == "" {
		s.opts.group = uuid.New().String()
	}

	for _, client := range s.clients {
		errCh := make(chan error, 1)
		eventsCh, err := client.SubscribeToEvents(ctx, s.opts.channel, s.opts.group, errCh)
		if err != nil {
			return fmt.Errorf("error on subscribing to events channel, %w", err)
		}
		go func(ctx context.Context, eventsCh <-chan *kubemq.Event, errCh chan error) {
			s.run(ctx, eventsCh, errCh)
		}(ctx, eventsCh, errCh)
	}

	return nil
}

func (s *Source) run(ctx context.Context, eventsCh <-chan *kubemq.Event, errCh chan error) {
	for {
		select {
		case event := <-eventsCh:
			if s.loadBalancingMode {
				go func(event *kubemq.Event, target middleware.Middleware) {
					_, err := target.Do(ctx, event)
					if err != nil {
						s.log.Errorf("error received from target, %s", err.Error())
					}
				}(event, s.targets[s.roundRobin.Next()])
			} else {
				for _, target := range s.targets {
					go func(event *kubemq.Event, target middleware.Middleware) {
						_, err := target.Do(ctx, event)
						if err != nil {
							s.log.Errorf("error received from target, %s", err.Error())
						}
					}(event, target)
				}
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
	for _, client := range s.clients {
		_ = client.Close()
	}
	return nil
}
