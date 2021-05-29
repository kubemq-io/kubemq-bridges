package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-go"
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
		c.log = logger.NewLogger("queue")
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
	var messages []*kubemq.QueueMessage
	switch val := request.(type) {
	case *kubemq.CommandReceive:
		messages = c.parseCommand(val, c.opts.channels)
	case *kubemq.Event:
		messages = c.parseEvent(val, c.opts.channels)
	case *kubemq.EventStoreReceive:
		messages = c.parseEventStore(val, c.opts.channels)
	case *kubemq.QueryReceive:
		messages = c.parseQuery(val, c.opts.channels)
	case *kubemq.QueueMessage:
		messages = c.parseQueue(val, c.opts.channels)
	default:
		return nil, fmt.Errorf("unknown request type")
	}
	results, err := c.client.SendQueueMessages(ctx, messages)
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		if result.IsError {
			return nil, fmt.Errorf(result.Error)
		}
	}
	return nil, nil
}

func (c *Client) parseEvent(event *kubemq.Event, channels []string) []*kubemq.QueueMessage {
	var messages []*kubemq.QueueMessage
	if len(channels) == 0 {
		channels = append(channels, event.Channel)
	}
	for _, channel := range channels {
		messages = append(messages, c.client.NewQueueMessage().
			SetChannel(channel).
			SetBody(event.Body).
			SetMetadata(event.Metadata).
			SetId(event.Id).
			SetTags(event.Tags).
			SetPolicyDelaySeconds(c.opts.delaySeconds).
			SetPolicyExpirationSeconds(c.opts.expirationSeconds).
			SetPolicyMaxReceiveCount(c.opts.maxReceiveCount).
			SetPolicyMaxReceiveQueue(c.opts.deadLetterQueue))
	}
	return messages

}
func (c *Client) parseEventStore(eventStore *kubemq.EventStoreReceive, channels []string) []*kubemq.QueueMessage {
	var messages []*kubemq.QueueMessage
	if len(channels) == 0 {
		channels = append(channels, eventStore.Channel)
	}
	for _, channel := range channels {
		messages = append(messages, c.client.NewQueueMessage().
			SetChannel(channel).
			SetBody(eventStore.Body).
			SetMetadata(eventStore.Metadata).
			SetId(eventStore.Id).
			SetTags(eventStore.Tags).
			SetPolicyDelaySeconds(c.opts.delaySeconds).
			SetPolicyExpirationSeconds(c.opts.expirationSeconds).
			SetPolicyMaxReceiveCount(c.opts.maxReceiveCount).
			SetPolicyMaxReceiveQueue(c.opts.deadLetterQueue))
	}
	return messages
}

func (c *Client) parseQuery(query *kubemq.QueryReceive, channels []string) []*kubemq.QueueMessage {
	var messages []*kubemq.QueueMessage
	if len(channels) == 0 {
		channels = append(channels, query.Channel)
	}
	for _, channel := range channels {
		messages = append(messages, c.client.NewQueueMessage().
			SetChannel(channel).
			SetBody(query.Body).
			SetMetadata(query.Metadata).
			SetId(query.Id).
			SetTags(query.Tags).
			SetPolicyDelaySeconds(c.opts.delaySeconds).
			SetPolicyExpirationSeconds(c.opts.expirationSeconds).
			SetPolicyMaxReceiveCount(c.opts.maxReceiveCount).
			SetPolicyMaxReceiveQueue(c.opts.deadLetterQueue))
	}
	return messages
}
func (c *Client) parseCommand(command *kubemq.CommandReceive, channels []string) []*kubemq.QueueMessage {
	var messages []*kubemq.QueueMessage
	if len(channels) == 0 {
		channels = append(channels, command.Channel)
	}
	for _, channel := range channels {
		messages = append(messages, c.client.NewQueueMessage().
			SetChannel(channel).
			SetBody(command.Body).
			SetMetadata(command.Metadata).
			SetId(command.Id).
			SetTags(command.Tags).
			SetPolicyDelaySeconds(c.opts.delaySeconds).
			SetPolicyExpirationSeconds(c.opts.expirationSeconds).
			SetPolicyMaxReceiveCount(c.opts.maxReceiveCount).
			SetPolicyMaxReceiveQueue(c.opts.deadLetterQueue))
	}
	return messages
}
func (c *Client) parseQueue(message *kubemq.QueueMessage, channels []string) []*kubemq.QueueMessage {
	var messages []*kubemq.QueueMessage
	if len(channels) == 0 {
		channels = append(channels, message.Channel)
	}
	for _, channel := range channels {
		messages = append(messages, c.client.NewQueueMessage().
			SetChannel(channel).
			SetBody(message.Body).
			SetMetadata(message.Metadata).
			SetId(message.MessageID).
			SetTags(message.Tags).
			SetPolicyDelaySeconds(c.opts.delaySeconds).
			SetPolicyExpirationSeconds(c.opts.expirationSeconds).
			SetPolicyMaxReceiveCount(c.opts.maxReceiveCount).
			SetPolicyMaxReceiveQueue(c.opts.deadLetterQueue))
	}
	return messages
}
