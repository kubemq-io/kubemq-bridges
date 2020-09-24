package query

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-go"
	"time"
)

type Client struct {
	opts   options
	client *kubemq.Client
}

func New() *Client {
	return &Client{}

}
func (c *Client) Init(ctx context.Context, connection config.Metadata) error {
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
		kubemq.WithCheckConnection(false),
	)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) Do(ctx context.Context, request interface{}) (interface{}, error) {
	var query *kubemq.Query
	switch val := request.(type) {
	case *kubemq.CommandReceive:
		query = c.parseCommand(val)
	case *kubemq.Event:
		query = c.parseEvent(val)
	case *kubemq.EventStoreReceive:
		query = c.parseEventStore(val)
	case *kubemq.QueryReceive:
		query = c.parseQuery(val)
	case *kubemq.QueueMessage:
		query = c.parseQueue(val)
	default:
		return nil, fmt.Errorf("unknown request type")
	}
	if c.opts.defaultChannel != "" {
		query.SetChannel(c.opts.defaultChannel)
	}
	query.SetTimeout(time.Duration(c.opts.timeoutSeconds) * time.Second)
	queryResponse, err := c.client.SetQuery(query).Send(ctx)
	if err != nil {
		return nil, err
	}
	if !queryResponse.Executed {
		return nil, fmt.Errorf(queryResponse.Error)
	}
	return queryResponse, nil
}

func (c *Client) parseEvent(event *kubemq.Event) *kubemq.Query {
	return kubemq.NewQuery().
		SetBody(event.Body).
		SetMetadata(event.Metadata).
		SetId(event.Id).
		SetTags(event.Tags).
		SetChannel(event.Channel)

}
func (c *Client) parseEventStore(eventStore *kubemq.EventStoreReceive) *kubemq.Query {
	return kubemq.NewQuery().
		SetBody(eventStore.Body).
		SetMetadata(eventStore.Metadata).
		SetId(eventStore.Id).
		SetTags(eventStore.Tags).
		SetChannel(eventStore.Channel)
}

func (c *Client) parseQuery(query *kubemq.QueryReceive) *kubemq.Query {
	return kubemq.NewQuery().
		SetBody(query.Body).
		SetMetadata(query.Metadata).
		SetId(query.Id).
		SetTags(query.Tags).
		SetChannel(query.Channel)
}
func (c *Client) parseCommand(command *kubemq.CommandReceive) *kubemq.Query {
	return kubemq.NewQuery().
		SetBody(command.Body).
		SetMetadata(command.Metadata).
		SetId(command.Id).
		SetTags(command.Tags).
		SetChannel(command.Channel)
}
func (c *Client) parseQueue(message *kubemq.QueueMessage) *kubemq.Query {
	return kubemq.NewQuery().
		SetBody(message.Body).
		SetMetadata(message.Metadata).
		SetId(message.MessageID).
		SetTags(message.Tags).
		SetChannel(message.Channel)
}
