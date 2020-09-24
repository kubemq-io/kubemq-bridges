package events

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-go"
	"time"
)

const (
	defaultSendTimeout     = 10 * time.Second
	defaultStreamReconnect = 1 * time.Second
)

type Client struct {
	opts   options
	client *kubemq.Client
	sendCh chan *kubemq.Event
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
	c.sendCh = make(chan *kubemq.Event, 1)
	go c.runStreamProcessing(ctx)
	return nil
}

func (c *Client) Do(ctx context.Context, request interface{}) (interface{}, error) {
	var events []*kubemq.Event
	switch val := request.(type) {
	case *kubemq.CommandReceive:
		events = c.parseCommand(val, c.opts.channels)
	case *kubemq.Event:
		events = c.parseEvent(val, c.opts.channels)
	case *kubemq.EventStoreReceive:
		events = c.parseEventStore(val, c.opts.channels)
	case *kubemq.QueryReceive:
		events = c.parseQuery(val, c.opts.channels)
	case *kubemq.QueueMessage:
		events = c.parseQueue(val, c.opts.channels)
	default:
		return nil, fmt.Errorf("unknown request type")
	}
	for _, event := range events {
		select {
		case c.sendCh <- event:
		case <-time.After(defaultSendTimeout):
			return nil, fmt.Errorf("error timeout on sending event")
		}
	}
	return nil, nil
}

func (c *Client) runStreamProcessing(ctx context.Context) {
	for {
		errCh := make(chan error, 1)
		go func() {
			c.client.StreamEvents(ctx, c.sendCh, errCh)
		}()
		select {
		case <-errCh:
			time.Sleep(defaultStreamReconnect)
		case <-ctx.Done():
			goto done
		}
	}
done:
}

func (c *Client) parseEvent(event *kubemq.Event, channels []string) []*kubemq.Event {
	var events []*kubemq.Event
	if len(channels) == 0 {
		channels = append(channels, event.Channel)
	}

	for _, channel := range channels {
		events = append(events, kubemq.NewEvent().
			SetChannel(channel).
			SetBody(event.Body).
			SetMetadata(event.Metadata).
			SetId(event.Id).
			SetTags(event.Tags))
	}
	return events

}
func (c *Client) parseEventStore(eventStore *kubemq.EventStoreReceive, channels []string) []*kubemq.Event {
	var events []*kubemq.Event
	if len(channels) == 0 {
		channels = append(channels, eventStore.Channel)
	}
	for _, channel := range channels {
		events = append(events, kubemq.NewEvent().
			SetChannel(channel).
			SetBody(eventStore.Body).
			SetMetadata(eventStore.Metadata).
			SetId(eventStore.Id).
			SetTags(eventStore.Tags))
	}
	return events
}

func (c *Client) parseQuery(query *kubemq.QueryReceive, channels []string) []*kubemq.Event {
	var events []*kubemq.Event
	if len(channels) == 0 {
		channels = append(channels, query.Channel)
	}
	for _, channel := range channels {
		events = append(events, kubemq.NewEvent().
			SetChannel(channel).
			SetBody(query.Body).
			SetMetadata(query.Metadata).
			SetId(query.Id).
			SetTags(query.Tags))
	}
	return events
}
func (c *Client) parseCommand(command *kubemq.CommandReceive, channels []string) []*kubemq.Event {
	var events []*kubemq.Event
	if len(channels) == 0 {
		channels = append(channels, command.Channel)
	}

	for _, channel := range channels {
		events = append(events, kubemq.NewEvent().
			SetChannel(channel).
			SetBody(command.Body).
			SetMetadata(command.Metadata).
			SetId(command.Id).
			SetTags(command.Tags))
	}
	return events
}
func (c *Client) parseQueue(message *kubemq.QueueMessage, channels []string) []*kubemq.Event {
	var events []*kubemq.Event
	if len(channels) == 0 {
		channels = append(channels, message.Channel)
	}
	for _, channel := range channels {
		events = append(events, kubemq.NewEvent().
			SetChannel(channel).
			SetBody(message.Body).
			SetMetadata(message.Metadata).
			SetId(message.MessageID).
			SetTags(message.Tags))
	}
	return events
}
