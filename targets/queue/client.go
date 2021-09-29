package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-go"
	"github.com/kubemq-io/kubemq-go/queues_stream"
)

type Client struct {
	log          *logger.Logger
	opts         options
	streamClient *queues_stream.QueuesStreamClient
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
	c.streamClient, err = queues_stream.NewQueuesStreamClient(ctx,
		queues_stream.WithAddress(c.opts.host, c.opts.port),
		queues_stream.WithClientId(c.opts.clientId),
		queues_stream.WithCheckConnection(true),
		queues_stream.WithAutoReconnect(true),
		queues_stream.WithAuthToken(c.opts.authToken),
		queues_stream.WithConnectionNotificationFunc(
			func(msg string) {
				c.log.Infof(msg)
			}),
	)
	if err != nil {
		return err
	}

	return nil
}
func (c *Client) Stop() error {
	if c.streamClient != nil {
		return c.streamClient.Close()
	}
	return nil
}
func (c *Client) Do(ctx context.Context, request interface{}) (interface{}, error) {
	var messages []*queues_stream.QueueMessage
	switch val := request.(type) {
	case *kubemq.CommandReceive:
		messages = c.parseCommand(val, c.opts.channels)
	case *kubemq.Event:
		messages = c.parseEvent(val, c.opts.channels)
	case *kubemq.EventStoreReceive:
		messages = c.parseEventStore(val, c.opts.channels)
	case *kubemq.QueryReceive:
		messages = c.parseQuery(val, c.opts.channels)
	case *queues_stream.QueueMessage:
		messages = c.parseQueueStream(val, c.opts.channels)
	case *kubemq.QueueMessage:
		messages = c.parseQueue(val, c.opts.channels)
	default:
		return nil, fmt.Errorf("unknown request type")
	}
	results, err := c.streamClient.Send(ctx, messages...)
	if err != nil {
		return nil, err
	}
	for _, result := range results.Results {
		if result.IsError {
			return nil, fmt.Errorf(result.Error)
		}
	}
	return nil, nil
}

func (c *Client) parseEvent(event *kubemq.Event, channels []string) []*queues_stream.QueueMessage {
	var messages []*queues_stream.QueueMessage
	if len(channels) == 0 {
		channels = append(channels, event.Channel)
	}
	for _, channel := range channels {
		messages = append(messages, queues_stream.NewQueueMessage().
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
func (c *Client) parseEventStore(eventStore *kubemq.EventStoreReceive, channels []string) []*queues_stream.QueueMessage {
	var messages []*queues_stream.QueueMessage

	for _, channel := range channels {
		messages = append(messages, queues_stream.NewQueueMessage().
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

func (c *Client) parseQuery(query *kubemq.QueryReceive, channels []string) []*queues_stream.QueueMessage {
	var messages []*queues_stream.QueueMessage

	for _, channel := range channels {
		messages = append(messages, queues_stream.NewQueueMessage().
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
func (c *Client) parseCommand(command *kubemq.CommandReceive, channels []string) []*queues_stream.QueueMessage {
	var messages []*queues_stream.QueueMessage

	for _, channel := range channels {
		messages = append(messages, queues_stream.NewQueueMessage().
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
func (c *Client) parseQueue(message *kubemq.QueueMessage, channels []string) []*queues_stream.QueueMessage {
	var messages []*queues_stream.QueueMessage

	for _, channel := range channels {
		messages = append(messages, queues_stream.NewQueueMessage().
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
func (c *Client) parseQueueStream(message *queues_stream.QueueMessage, channels []string) []*queues_stream.QueueMessage {
	var messages []*queues_stream.QueueMessage

	for _, channel := range channels {
		messages = append(messages, queues_stream.NewQueueMessage().
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
