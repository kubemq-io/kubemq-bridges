package queue

import (
	"context"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"time"

	"github.com/kubemq-io/kubemq-go"

	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
)

const (
	retriesInterval = 1 * time.Second
)

type Source struct {
	opts      options
	client    *kubemq.Client
	log       *logger.Logger
	targets   []middleware.Middleware
	isStopped bool
	properties config.Metadata
}

func New() *Source {
	return &Source{}

}
func (s *Source) Init(ctx context.Context, connection config.Metadata,properties config.Metadata) error {
	var err error
	s.opts, err = parseOptions(connection)
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
	)
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) Start(ctx context.Context, targets []middleware.Middleware, log *logger.Logger) error {
	s.log = log
	s.targets = targets
	go s.run(ctx)
	return nil
}

func (s *Source) run(ctx context.Context) {
	for {
		if s.isStopped {
			return
		}
		queueMessages, err := s.getQueueMessages(ctx)
		if err != nil {
			s.log.Error(err.Error())
			time.Sleep(retriesInterval)
			continue
		}
		for _, message := range queueMessages {
			for _, target := range s.targets {
				err := s.processQueueMessage(ctx, message, target)
				if err != nil {
					s.log.Errorf("error received from target, %w", err)
				}
			}

		}
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}
func (s *Source) getQueueMessages(ctx context.Context) ([]*kubemq.QueueMessage, error) {
	receiveResult, err := s.client.NewReceiveQueueMessagesRequest().
		SetChannel(s.opts.channel).
		SetMaxNumberOfMessages(s.opts.batchSize).
		SetWaitTimeSeconds(s.opts.waitTimeout).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	return receiveResult.Messages, nil
}

func (s *Source) processQueueMessage(ctx context.Context, msg *kubemq.QueueMessage, target middleware.Middleware) error {
	_, err := target.Do(ctx, msg)
	if err != nil {
		return err
	}
	return nil

}

func (s *Source) Stop() error {
	s.isStopped = true
	return s.client.Close()
}
