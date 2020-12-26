package query

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"

	"github.com/nats-io/nuid"

	"github.com/kubemq-io/kubemq-go"
)

type Source struct {
	opts    options
	client  *kubemq.Client
	log     *logger.Logger
	targets []middleware.Middleware
	properties config.Metadata
}

func New() *Source {
	return &Source{}
}
func (s *Source) Init(ctx context.Context, connection config.Metadata,properties config.Metadata) error {
	var err error
	s.opts, err = parseOptions(connection)
	if err != nil {
		return err
	}
	s.properties=properties
	s.client, err = kubemq.NewClient(ctx,
		kubemq.WithAddress(s.opts.host, s.opts.port),
		kubemq.WithClientId(s.opts.clientId),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC),
		kubemq.WithAuthToken(s.opts.authToken),
		kubemq.WithCheckConnection(true),
		kubemq.WithMaxReconnects(s.opts.maxReconnects),
		kubemq.WithAutoReconnect(s.opts.autoReconnect),
		kubemq.WithReconnectInterval(s.opts.reconnectIntervalSeconds))
	if err != nil {
		return err
	}
	return nil
}

func (s *Source) runSubscriber(ctx context.Context, channel, group string, target middleware.Middleware) error {
	errCh := make(chan error, 1)
	queriesCh, err := s.client.SubscribeToQueries(ctx, channel, group, errCh)
	if err != nil {
		return fmt.Errorf("error on subscribing to query channel, %w", err)
	}
	go func(ctx context.Context, commandCh <-chan *kubemq.QueryReceive, errCh chan error, target middleware.Middleware) {
		s.run(ctx, queriesCh, errCh, target)
	}(ctx, queriesCh, errCh, target)
	return nil
}
func (s *Source) Start(ctx context.Context, targets []middleware.Middleware, log *logger.Logger) error {
	s.log = log
	s.targets = targets
	group := nuid.Next()
	if s.opts.group != "" {
		group = s.opts.group
	}
	for _, target := range targets {
		err := s.runSubscriber(ctx, s.opts.channel, group, target)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *Source) run(ctx context.Context, queryCh <-chan *kubemq.QueryReceive, errCh chan error, target middleware.Middleware) {
	for {
		select {
		case query := <-queryCh:

			go func(q *kubemq.QueryReceive) {
				var queryResponse *kubemq.Response
				queryResponse, err := s.processQuery(ctx, query, target)
				if err != nil {
					queryResponse = s.client.NewResponse().
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
func (s *Source) processQuery(ctx context.Context, query *kubemq.QueryReceive, target middleware.Middleware) (*kubemq.Response, error) {
	result, err := target.Do(ctx, query)
	if err != nil {
		return nil, err
	}
	switch val := result.(type) {
	case *kubemq.CommandResponse:
		return s.parseCommandResponse(val), nil
	case *kubemq.QueryResponse:
		return s.parseQueryResponse(val), nil
	default:
		return s.client.NewResponse(), nil
	}
}
func (s *Source) Stop() error {
	return s.client.Close()
}

func (s *Source) parseCommandResponse(cmd *kubemq.CommandResponse) *kubemq.Response {
	resp := s.client.NewResponse().SetTags(cmd.Tags)
	if cmd.Executed {
		resp.SetExecutedAt(cmd.ExecutedAt)
	} else {
		resp.SetError(fmt.Errorf("%s", cmd.Error))
	}
	return resp
}
func (s *Source) parseQueryResponse(query *kubemq.QueryResponse) *kubemq.Response {
	resp := s.client.NewResponse().SetTags(query.Tags).SetMetadata(query.Metadata).SetBody(query.Body)
	if query.Executed {
		resp.SetExecutedAt(query.ExecutedAt)
	} else {
		resp.SetError(fmt.Errorf("%s", query.Error))
	}
	return resp
}
