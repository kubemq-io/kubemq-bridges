package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/middleware"
	"github.com/kubemq-io/kubemq-bridges/pkg/roundrobin"
	"time"

	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-go/queues_stream"
)

type Source struct {
	opts options

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

func (s *Source) getQueuesClient(ctx context.Context, id int) (*queues_stream.QueuesStreamClient, error) {
	return queues_stream.NewQueuesStreamClient(ctx,
		queues_stream.WithAddress(s.opts.host, s.opts.port),
		queues_stream.WithClientId(s.opts.clientId),
		queues_stream.WithCheckConnection(true),
		queues_stream.WithAutoReconnect(true),
		queues_stream.WithAuthToken(s.opts.authToken),
		queues_stream.WithConnectionNotificationFunc(
			func(msg string) {
				s.log.Infof(fmt.Sprintf("connection: %d, %s", id, msg))
			}),
	)

}
func (s *Source) onError(err error) {
	s.log.Error(err.Error())
}
func (s *Source) Init(ctx context.Context, connection config.Metadata, properties config.Metadata, bindingName string, log *logger.Logger) error {
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
		client, err := s.getQueuesClient(ctx, i+1)
		if err != nil {
			return err
		}
		go s.run(ctx, client)
	}
	return nil
}

func (s *Source) run(ctx context.Context, client *queues_stream.QueuesStreamClient) {
	defer func() {
		_ = client.Close()
	}()
	for {
		if s.isStopped {
			return
		}
		err := s.processQueueMessage(ctx, client)
		if err != nil {
			s.log.Error(err.Error())
			time.Sleep(time.Second)
		}
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}
func (s *Source) processQueueMessage(ctx context.Context, client *queues_stream.QueuesStreamClient) error {
	pr := queues_stream.NewPollRequest().
		SetChannel(s.opts.channel).
		SetMaxItems(s.opts.batchSize).
		SetWaitTimeout(s.opts.waitTimeout * 1000).
		SetAutoAck(false).
		SetOnErrorFunc(s.onError)
	pollResp, err := client.Poll(ctx, pr)
	if err != nil {
		return err
	}
	if !pollResp.HasMessages() {
		return nil
	}
	for _, message := range pollResp.Messages {
		if s.loadBalancingMode {
			_, err := s.targets[s.roundRobin.Next()].Do(ctx, message)
			if err != nil {
				if message.Policy.MaxReceiveCount < 1024 && message.Policy.MaxReceiveCount != message.Attributes.ReceiveCount {
					return message.NAck()
				}
			}

		} else {
			wasExecuted := false
			for _, target := range s.targets {
				_, err := target.Do(ctx, message)
				if err == nil {
					wasExecuted = true
				}
			}
			if !wasExecuted {
				if message.Policy.MaxReceiveCount < 1024 && message.Policy.MaxReceiveCount != message.Attributes.ReceiveCount {
					return message.NAck()
				}

			}
		}
		err = message.Ack()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Source) Stop() error {
	s.isStopped = true
	return nil
}
