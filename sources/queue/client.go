package queue

import (
	"context"
	"github.com/kubemq-hub/kubemq-bridges/middleware"

	"github.com/kubemq-io/kubemq-go"

	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
)

type Client struct {
	name   string
	opts   options
	client *kubemq.Client
	log    *logger.Logger
	target middleware.Middleware
}

func New() *Client {
	return &Client{}

}
func (c *Client) Name() string {
	return c.name
}
func (c *Client) Init(ctx context.Context, cfg config.Spec) error {
	c.name = cfg.Name
	c.log = logger.NewLogger(fmt.Sprintf("kubemq-queue-source-%s", cfg.Name))
	var err error
	c.opts, err = parseOptions(cfg.Properties)
	if err != nil {
		return err
	}
	c.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {
	c.target = target
	go c.run(ctx)
	return nil
}

func (c *Client) run(ctx context.Context) {
	for {
		queueMessages, err := c.getQueueMessages(ctx)
		if err != nil {
			c.log.Error(err.Error())
			return
		}
		for _, message := range queueMessages {
			err := c.processQueueMessage(ctx, message)
			if err != nil {
				c.log.Errorf("error received from target, %w", err)
			}
		}
		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}
func (c *Client) getQueueMessages(ctx context.Context) ([]*kubemq.QueueMessage, error) {
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

func (c *Client) processQueueMessage(ctx context.Context, msg *kubemq.QueueMessage) error {
	_, err := c.target.Do(ctx, msg)
	if err != nil {
		return err
	}
	return nil

}

func (c *Client) Stop() error {
	return c.client.Close()
}
