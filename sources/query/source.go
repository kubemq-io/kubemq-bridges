package query

import (
	"context"
	"fmt"

	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/middleware"
	"github.com/kubemq-io/kubemq-bridges/pkg/logger"

	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"

	"github.com/kubemq-io/kubemq-go"
)

type Source struct {
	opts       options
	clients    []*kubemq.Client
	log        *logger.Logger
	targets    []middleware.Middleware
	properties config.Metadata
}

func New() *Source {
	return &Source{}
}

func (s *Source) Init(ctx context.Context, connection config.Metadata, properties config.Metadata, bindingName string, log *logger.Logger) error {
	s.log = log
	if s.log == nil {
		s.log = logger.NewLogger("query")
	}
	var err error
	s.opts, err = parseOptions(connection)
	if err != nil {
		return err
	}
	s.properties = properties
	for i := 0; i < s.opts.sources; i++ {
		clientId := fmt.Sprintf("kubemq-bridges_%s_%s", bindingName, s.opts.clientId)
		if s.opts.sources > 1 {
			clientId = fmt.Sprintf("kubemq-bridges_%s_%s-%d", bindingName, clientId, i)
		}
		client, err := kubemq.NewClient(ctx,
			kubemq.WithAddress(s.opts.host, s.opts.port),
			kubemq.WithClientId(clientId),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC),
			kubemq.WithCheckConnection(true),
			kubemq.WithAuthToken(s.opts.authToken),
			kubemq.WithMaxReconnects(s.opts.maxReconnects),
			kubemq.WithAutoReconnect(s.opts.autoReconnect),
			kubemq.WithReconnectInterval(s.opts.reconnectIntervalSeconds))
		if err != nil {
			return err
		}
		s.clients = append(s.clients, client)
	}
	return nil
}

func (s *Source) Start(ctx context.Context, target []middleware.Middleware) error {
	s.targets = target
	if s.opts.sources > 1 && s.opts.group == "" {
		s.opts.group = uuid.New().String()
	}
	for _, client := range s.clients {
		for _, target := range target {
			err := s.runSubscriber(ctx, s.opts.channel, s.opts.group, target, client)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Source) runSubscriber(ctx context.Context, channel, group string, target middleware.Middleware, client *kubemq.Client) error {
	errCh := make(chan error, 1)
	queriesCh, err := client.SubscribeToQueries(ctx, channel, group, errCh)
	if err != nil {
		return fmt.Errorf("error on subscribing to query channel, %w", err)
	}
	go func(ctx context.Context, commandCh <-chan *kubemq.QueryReceive, errCh chan error, target middleware.Middleware, client *kubemq.Client) {
		s.run(ctx, queriesCh, errCh, target, client)
	}(ctx, queriesCh, errCh, target, client)
	return nil
}

func (s *Source) run(ctx context.Context, queryCh <-chan *kubemq.QueryReceive, errCh chan error, target middleware.Middleware, client *kubemq.Client) {
	for {
		select {
		case query := <-queryCh:

			go func(q *kubemq.QueryReceive) {
				var queryResponse *kubemq.Response
				queryResponse, err := s.processQuery(ctx, query, target, client)
				if err != nil {
					queryResponse = client.NewResponse().
						SetRequestId(query.Id).
						SetResponseTo(query.ResponseTo).
						SetError(err)
				} else {
					queryResponse.
						SetRequestId(query.Id).
						SetResponseTo(query.ResponseTo)
				}
				err = queryResponse.Send(ctx)
				if err != nil {
					s.log.Errorf("error sending query response %s", err.Error())
				}
			}(query)

		case err := <-errCh:
			s.log.Errorf("error received from kuebmq server, %s", err.Error())
			return
		case <-ctx.Done():
			return

		}
	}
}

func (s *Source) processQuery(ctx context.Context, query *kubemq.QueryReceive, target middleware.Middleware, client *kubemq.Client) (*kubemq.Response, error) {
	result, err := target.Do(ctx, query)
	if err != nil {
		return nil, err
	}
	switch val := result.(type) {
	case *kubemq.CommandResponse:
		return s.parseCommandResponse(val, client), nil
	case *kubemq.QueryResponse:
		return s.parseQueryResponse(val, client), nil
	default:
		return client.NewResponse(), nil
	}
}

func (s *Source) Stop() error {
	for _, client := range s.clients {
		_ = client.Close()
	}
	return nil
}

func (s *Source) parseCommandResponse(cmd *kubemq.CommandResponse, client *kubemq.Client) *kubemq.Response {
	resp := client.NewResponse().SetTags(cmd.Tags)
	if cmd.Executed {
		resp.SetExecutedAt(cmd.ExecutedAt)
	} else {
		resp.SetError(fmt.Errorf("%s", cmd.Error))
	}
	return resp
}

func (s *Source) parseQueryResponse(query *kubemq.QueryResponse, client *kubemq.Client) *kubemq.Response {
	resp := client.NewResponse().SetTags(query.Tags).SetMetadata(query.Metadata).SetBody(query.Body)
	if query.Executed {
		resp.SetExecutedAt(query.ExecutedAt)
	} else {
		resp.SetError(fmt.Errorf("%s", query.Error))
	}
	return resp
}
