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
	c.log = logger.NewLogger(cfg.Name)
	var err error
	c.opts, err = parseOptions(cfg.Properties)
	if err != nil {
		return err
	}
	c.client, _ = kubemq.NewClient(ctx,
		kubemq.WithAddress(c.opts.host, c.opts.port),
		kubemq.WithClientId(c.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(c.opts.authToken),
		kubemq.WithCheckConnection(true),
		kubemq.WithMaxReconnects(c.opts.maxReconnects),
		kubemq.WithAutoReconnect(c.opts.autoReconnect),
		kubemq.WithReconnectInterval(c.opts.reconnectIntervalSeconds))
	return nil
}

func (c *Client) Start(ctx context.Context, target middleware.Middleware) error {

	c.target = target

	group := nuid.Next()
	if c.opts.group != "" {
		group = c.opts.group
	}

	errCh := make(chan error, 1)
	eventsCh, err := c.client.SubscribeToEventsStore(ctx, c.opts.channel, group, errCh, kubemq.StartFromNewEvents())
	if err != nil {
		return fmt.Errorf("error on subscribing to events channel, %w", err)
	}
	go func(ctx context.Context, eventsCh <-chan *kubemq.EventStoreReceive, errCh chan error) {
		c.run(ctx, eventsCh, errCh)
	}(ctx, eventsCh, errCh)

	return nil
}

func (c *Client) run(ctx context.Context, eventsCh <-chan *kubemq.EventStoreReceive, errCh chan error) {
	for {
		select {
		case event := <-eventsCh:
			go func(event *kubemq.EventStoreReceive) {
				err := c.processEventStore(ctx, event)
				if err != nil {
					c.log.Errorf("error received from target, %w", err)
				}
			}(event)

		case err := <-errCh:
			c.log.Errorf("error received from kuebmq server, %s", err.Error())
			return
		case <-ctx.Done():
			return

		}
	}
}

func (c *Client) processEventStore(ctx context.Context, event *kubemq.EventStoreReceive) error {
	_, err := c.target.Do(ctx, event)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) Stop() error {
	return c.client.Close()
}
