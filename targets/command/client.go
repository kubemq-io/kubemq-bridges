package command

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-go"
	"time"
)

type Client struct {
	log    *logger.Logger
	opts   options
	client *kubemq.Client
}

func New() *Client {
	return &Client{}

}

func (c *Client) Init(ctx context.Context, connection config.Metadata, log *logger.Logger) error {
	c.log = log
	if c.log == nil {
		c.log = logger.NewLogger("commands")
	}
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
func (c *Client) Stop() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
func (c *Client) Do(ctx context.Context, request interface{}) (interface{}, error) {

	var cmd *kubemq.Command
	switch val := request.(type) {
	case *kubemq.CommandReceive:
		cmd = c.parseCommand(val)
	case *kubemq.Event:
		cmd = c.parseEvent(val)
	case *kubemq.EventStoreReceive:
		cmd = c.parseEventStore(val)
	case *kubemq.QueryReceive:
		cmd = c.parseQuery(val)
	case *kubemq.QueueMessage:
		cmd = c.parseQueue(val)
	default:
		return nil, fmt.Errorf("unknown request type")
	}
	if c.opts.defaultChannel != "" {
		cmd.SetChannel(c.opts.defaultChannel)
	}
	cmd.SetTimeout(time.Duration(c.opts.timeoutSeconds) * time.Second)
	cmdResponse, err := c.client.SetCommand(cmd).Send(ctx)
	if err != nil {
		return nil, err
	}
	if !cmdResponse.Executed {
		return nil, fmt.Errorf(cmdResponse.Error)
	}
	return cmdResponse, nil

}

func (c *Client) parseEvent(event *kubemq.Event) *kubemq.Command {
	return kubemq.NewCommand().
		SetBody(event.Body).
		SetMetadata(event.Metadata).
		SetId(event.Id).
		SetTags(event.Tags).
		SetChannel(event.Channel)

}
func (c *Client) parseEventStore(eventStore *kubemq.EventStoreReceive) *kubemq.Command {
	return kubemq.NewCommand().
		SetBody(eventStore.Body).
		SetMetadata(eventStore.Metadata).
		SetId(eventStore.Id).
		SetTags(eventStore.Tags).
		SetChannel(eventStore.Channel)
}

func (c *Client) parseQuery(query *kubemq.QueryReceive) *kubemq.Command {
	return kubemq.NewCommand().
		SetBody(query.Body).
		SetMetadata(query.Metadata).
		SetId(query.Id).
		SetTags(query.Tags).
		SetChannel(query.Channel)
}
func (c *Client) parseCommand(command *kubemq.CommandReceive) *kubemq.Command {
	return kubemq.NewCommand().
		SetBody(command.Body).
		SetMetadata(command.Metadata).
		SetId(command.Id).
		SetTags(command.Tags).
		SetChannel(command.Channel)
}
func (c *Client) parseQueue(message *kubemq.QueueMessage) *kubemq.Command {
	return kubemq.NewCommand().
		SetBody(message.Body).
		SetMetadata(message.Metadata).
		SetId(message.MessageID).
		SetTags(message.Tags).
		SetChannel(message.Channel)
}
