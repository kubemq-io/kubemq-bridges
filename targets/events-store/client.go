package events_store

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
	sendCh chan *kubemq.EventStore
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
		kubemq.WithCheckConnection(true),
	)

	if err != nil {
		return err
	}
	c.sendCh = make(chan *kubemq.EventStore, 1)
	go c.runStreamProcessing(ctx)
	return nil
}
func (c *Client) Stop() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
func (c *Client) Do(ctx context.Context, request interface{}) (interface{}, error) {
	var eventsStore []*kubemq.EventStore
	switch val := request.(type) {
	case *kubemq.CommandReceive:
		eventsStore = c.parseCommand(val, c.opts.channels)
	case *kubemq.Event:
		eventsStore = c.parseEvent(val, c.opts.channels)
	case *kubemq.EventStoreReceive:
		eventsStore = c.parseEventStore(val, c.opts.channels)
	case *kubemq.QueryReceive:
		eventsStore = c.parseQuery(val, c.opts.channels)
	case *kubemq.QueueMessage:
		eventsStore = c.parseQueue(val, c.opts.channels)
	default:
		return nil, fmt.Errorf("unknown request type")
	}
	for _, es := range eventsStore {
		select {
		case c.sendCh <- es:

		case <-time.After(defaultSendTimeout):
			return nil, fmt.Errorf("error timeout on sending event store")
		}
	}
	return nil, nil
}

func (c *Client) runStreamProcessing(ctx context.Context) {
	for {
		errCh := make(chan error, 1)

		go func() {
			resultCh := make(chan *kubemq.EventStoreResult, 1)
			c.client.StreamEventsStore(ctx, c.sendCh, resultCh, errCh)
			for {
				select {
				case <-resultCh:
				case <-ctx.Done():
					return
				}
			}
		}()
		select {
		case <-errCh:
			time.Sleep(defaultStreamReconnect)
			return
		case <-ctx.Done():
			goto done
		}
	}
done:
}

func (c *Client) parseEvent(event *kubemq.Event, channels []string) []*kubemq.EventStore {
	var eventsStores []*kubemq.EventStore
	if len(channels) == 0 {
		channels = append(channels, event.Channel)
	}
	for _, channel := range channels {
		eventsStores = append(eventsStores, kubemq.NewEventStore().
			SetChannel(channel).
			SetBody(event.Body).
			SetMetadata(event.Metadata).
			SetId(event.Id).
			SetTags(event.Tags))
	}
	return eventsStores

}
func (c *Client) parseEventStore(eventStore *kubemq.EventStoreReceive, channels []string) []*kubemq.EventStore {
	var eventsStores []*kubemq.EventStore
	if len(channels) == 0 {
		channels = append(channels, eventStore.Channel)
	}
	for _, channel := range channels {
		eventsStores = append(eventsStores, kubemq.NewEventStore().
			SetChannel(channel).
			SetBody(eventStore.Body).
			SetMetadata(eventStore.Metadata).
			SetId(eventStore.Id).
			SetTags(eventStore.Tags))
	}
	return eventsStores
}

func (c *Client) parseQuery(query *kubemq.QueryReceive, channels []string) []*kubemq.EventStore {
	var eventsStores []*kubemq.EventStore
	if len(channels) == 0 {
		channels = append(channels, query.Channel)
	}
	for _, channel := range channels {
		eventsStores = append(eventsStores, kubemq.NewEventStore().
			SetChannel(channel).
			SetBody(query.Body).
			SetMetadata(query.Metadata).
			SetId(query.Id).
			SetTags(query.Tags))
	}
	return eventsStores
}
func (c *Client) parseCommand(command *kubemq.CommandReceive, channels []string) []*kubemq.EventStore {
	var eventsStores []*kubemq.EventStore
	if len(channels) == 0 {
		channels = append(channels, command.Channel)
	}
	for _, channel := range channels {
		eventsStores = append(eventsStores, kubemq.NewEventStore().
			SetChannel(channel).
			SetBody(command.Body).
			SetMetadata(command.Metadata).
			SetId(command.Id).
			SetTags(command.Tags))
	}
	return eventsStores
}
func (c *Client) parseQueue(message *kubemq.QueueMessage, channels []string) []*kubemq.EventStore {
	var eventsStores []*kubemq.EventStore
	if len(channels) == 0 {
		channels = append(channels, message.Channel)
	}
	for _, channel := range channels {
		eventsStores = append(eventsStores, kubemq.NewEventStore().
			SetChannel(channel).
			SetBody(message.Body).
			SetMetadata(message.Metadata).
			SetId(message.MessageID).
			SetTags(message.Tags))
	}
	return eventsStores
}
