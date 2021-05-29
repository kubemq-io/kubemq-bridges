package queue_stream

import (
	"context"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/roundrobin"

	"github.com/kubemq-io/kubemq-go"

	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
	"strings"
)

type Source struct {
	opts              options
	client            *kubemq.Client
	log               *logger.Logger
	targets           []middleware.Middleware
	isStopped         bool
	properties        config.Metadata
	roundRobin        *roundrobin.RoundRobin
	loadBalancingMode bool
}

func New() *Source {
	return &Source{}

}

func (s *Source) getKubemqClient(ctx context.Context) (*kubemq.Client, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress(s.opts.host, s.opts.port),
		kubemq.WithClientId(s.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(s.opts.authToken))
	if err != nil {
		return nil, err
	}
	return client, nil
}
func (s *Source) Init(ctx context.Context, connection config.Metadata, properties config.Metadata, log *logger.Logger) error {
	s.log = log
	if s.log == nil {
		s.log = logger.NewLogger("queue-stream")
	}
	var err error
	s.opts, err = parseOptions(connection)
	s.properties = properties
	if err != nil {
		return err
	}
	s.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(s.opts.host, s.opts.port),
		kubemq.WithClientId(s.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithCheckConnection(true),
		kubemq.WithAuthToken(s.opts.authToken))
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) Start(ctx context.Context, target []middleware.Middleware) error {
	s.roundRobin = roundrobin.NewRoundRobin(len(target))
	if s.properties != nil {
		mode, ok := s.properties["load-balancing"]
		if ok && mode == "true" {
			s.loadBalancingMode = true
		}
	}
	s.targets = target
	for i := 0; i < s.opts.sources; i++ {
		go s.run(ctx)
	}
	return nil
}

func (s *Source) run(ctx context.Context) {
	for {
		if s.isStopped {
			return
		}
		err := s.processQueueMessage()
		if err != nil {
			if !strings.Contains(err.Error(), "138") {
				s.log.Error(err.Error())
			}
		}
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}
func (s *Source) processQueueMessage() error {
	ctx := context.Background()
	client, err := s.getKubemqClient(ctx)
	if err != nil {
		return err
	}
	stream := client.NewStreamQueueMessage().SetChannel(s.opts.channel)
	defer func() {
		stream.Close()
	}()
	msg, err := stream.Next(ctx, int32(s.opts.visibilityTimeout), int32(s.opts.waitTimeout))
	if err != nil {
		return err
	}
	if s.loadBalancingMode {
		_, err := s.targets[s.roundRobin.Next()].Do(ctx, msg)
		if err != nil {
			if msg.Policy.MaxReceiveCount != msg.Attributes.ReceiveCount {
				return msg.Reject()
			}
			return nil
		}
	} else {
		wasExecuted := false
		for _, target := range s.targets {
			_, err := target.Do(ctx, msg)
			if err == nil {
				wasExecuted = true
			}
		}
		if !wasExecuted {
			if msg.Policy.MaxReceiveCount != msg.Attributes.ReceiveCount {
				return msg.Reject()
			}
			return nil
		}

	}

	return msg.Ack()
}

func (s *Source) Stop() error {
	s.isStopped = true
	_ = s.client.Close()
	return nil
}
