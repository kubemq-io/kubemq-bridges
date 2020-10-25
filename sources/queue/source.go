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
}

func New() *Source {
	return &Source{}

}
func (c *Source) Init(ctx context.Context, connection config.Metadata) error {
	var err error
	c.opts, err = parseOptions(connection)
	if err != nil {
		return err
	}
	c.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken),
		kubemq.WithCheckConnection(true),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Source) Start(ctx context.Context, targets []middleware.Middleware, log *logger.Logger) error {
	c.log = log
	c.targets = targets
	go c.run(ctx)
	return nil
}

func (c *Source) run(ctx context.Context) {
	for {
		if c.isStopped {
			return
		}
		queueMessages, err := c.getQueueMessages(ctx)
		if err != nil {
			c.log.Error(err.Error())
			time.Sleep(retriesInterval)
			continue
		}
		for _, message := range queueMessages {
			for _, target := range c.targets {
				err := c.processQueueMessage(ctx, message, target)
				if err != nil {
					c.log.Errorf("error received from target, %w", err)
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
func (c *Source) getQueueMessages(ctx context.Context) ([]*kubemq.QueueMessage, error) {
	receiveResult, err := c.client.NewReceiveQueueMessagesRequest().
		SetChannel(c.opts.channel).
		SetMaxNumberOfMessages(c.opts.batchSize).
		SetWaitTimeSeconds(c.opts.waitTimeout).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	return receiveResult.Messages, nil
}

func (c *Source) processQueueMessage(ctx context.Context, msg *kubemq.QueueMessage, target middleware.Middleware) error {
	_, err := target.Do(ctx, msg)
	if err != nil {
		return err
	}
	return nil

}

func (c *Source) Stop() error {
	c.isStopped = true
	return c.client.Close()
}
