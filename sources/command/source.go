package command

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
)

type Source struct {
	opts    options
	client  *kubemq.Client
	log     *logger.Logger
	targets []middleware.Middleware
}

func New() *Source {
	return &Source{}

}
func (s *Source) Init(ctx context.Context, connection config.Metadata) error {
	var err error
	s.opts, err = parseOptions(connection)
	if err != nil {
		return err
	}
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
func (s *Source) Start(ctx context.Context, targets []middleware.Middleware, log *logger.Logger) error {
	s.targets = targets
	s.log = log
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

func (s *Source) runSubscriber(ctx context.Context, channel, group string, target middleware.Middleware) error {
	errCh := make(chan error, 1)
	commandsCh, err := s.client.SubscribeToCommands(ctx, channel, group, errCh)
	if err != nil {
		return fmt.Errorf("error on subscribing to command channel, %w", err)
	}
	go func(ctx context.Context, commandCh <-chan *kubemq.CommandReceive, errCh chan error, target middleware.Middleware) {
		s.run(ctx, commandsCh, errCh, target)
	}(ctx, commandsCh, errCh, target)
	return nil
}
func (s *Source) run(ctx context.Context, commandCh <-chan *kubemq.CommandReceive, errCh chan error, target middleware.Middleware) {
	for {
		select {
		case command := <-commandCh:
			go func(q *kubemq.CommandReceive) {
				var cmdResponse *kubemq.Response
				cmdResponse, err := s.processCommand(ctx, command, target)
				if err != nil {
					cmdResponse = s.client.NewResponse().
						SetRequestId(command.Id).
						SetResponseTo(command.ResponseTo).
						SetError(err)
				} else {
					cmdResponse.
						SetRequestId(command.Id).
						SetResponseTo(command.ResponseTo)
				}
				err = cmdResponse.Send(ctx)
				if err != nil {
					s.log.Errorf("error sending command response %s", err.Error())
				}
			}(command)

		case err := <-errCh:
			s.log.Errorf("error received from kuebmq server, %s", err.Error())
			return
		case <-ctx.Done():
			return

		}
	}
}

func (s *Source) processCommand(ctx context.Context, command *kubemq.CommandReceive, target middleware.Middleware) (*kubemq.Response, error) {

	result, err := target.Do(ctx, command)
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
	resp := s.client.NewResponse().SetTags(query.Tags)
	if query.Executed {
		resp.SetExecutedAt(query.ExecutedAt)
	} else {
		resp.SetError(fmt.Errorf("%s", query.Error))
	}
	return resp
}

func (s *Source) Stop() error {
	return s.client.Close()
}
