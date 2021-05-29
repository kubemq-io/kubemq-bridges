package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/roundrobin"
	"time"

	"github.com/kubemq-io/kubemq-go"

	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/logger"
)

const (
	retriesInterval = 1 * time.Second
)

type Source struct {
	opts              options
	clients           []*kubemq.Client
	log               *logger.Logger
	targets           []middleware.Middleware
	isStopped         bool
	properties        config.Metadata
	roundRobin        *roundrobin.RoundRobin
	loadBalancingMode bool
	requeueCache      *requeue
}

func New() *Source {
	return &Source{}

}
func (s *Source) Init(ctx context.Context, connection config.Metadata, properties config.Metadata, log *logger.Logger) error {
	s.log = log
	if s.log == nil {
		s.log = logger.NewLogger("queue")
	}
	var err error
	s.opts, err = parseOptions(connection)
	if err != nil {
		return err
	}
	s.requeueCache = newRequeue(s.opts.maxRequeue)
	s.properties = properties
	for i := 0; i < s.opts.sources; i++ {
		clientId := s.opts.clientId
		if s.opts.sources > 1 {
			clientId = fmt.Sprintf("%s-%d", clientId, i)
		}
		client, err := kubemq.NewClient(ctx,
			kubemq.WithAddress(s.opts.host, s.opts.port),
			kubemq.WithClientId(clientId),
			kubemq.WithTransportType(kubemq.TransportTypeGRPC),
			kubemq.WithCheckConnection(true),
			kubemq.WithAuthToken(s.opts.authToken),
		)
		if err != nil {
			return err
		}
		s.clients = append(s.clients, client)
	}
	return nil
}

func (s *Source) Start(ctx context.Context, target []middleware.Middleware) error {
	s.roundRobin = roundrobin.NewRoundRobin(len(target))
	if s.properties != nil {
		mode, ok := s.properties["load-balancing"]
		if ok && mode == "true" {
			s.loadBalancingMode = true
		}
	}
	s.targets = target
	for i := 0; i < len(s.clients); i++ {
		go s.run(ctx, s.clients[i])
	}
	return nil
}

func (s *Source) run(ctx context.Context, client *kubemq.Client) {
	for {
		if s.isStopped {
			return
		}
		queueMessages, err := s.getQueueMessages(ctx, client)
		if err != nil {
			s.log.Error(err.Error())
			time.Sleep(retriesInterval)
			continue
		}
		if s.loadBalancingMode {
			for _, message := range queueMessages {
				err := s.processQueueMessage(ctx, message, s.targets[s.roundRobin.Next()], client)
				if err != nil {
					s.log.Errorf("error received from target, %w", err)
				}
			}
		} else {
			for _, message := range queueMessages {
				for _, target := range s.targets {
					err := s.processQueueMessage(ctx, message, target, client)
					if err != nil {
						s.log.Errorf("error received from target, %w", err)
					}
				}

			}
		}

		select {
		case <-ctx.Done():
			return
		default:

		}
	}
}
func (s *Source) getQueueMessages(ctx context.Context, client *kubemq.Client) ([]*kubemq.QueueMessage, error) {
	receiveResult, err := client.NewReceiveQueueMessagesRequest().
		SetChannel(s.opts.channel).
		SetMaxNumberOfMessages(s.opts.batchSize).
		SetWaitTimeSeconds(s.opts.waitTimeout).
		Send(ctx)
	if err != nil {
		return nil, err
	}
	return receiveResult.Messages, nil
}

func (s *Source) processQueueMessage(ctx context.Context, msg *kubemq.QueueMessage, target middleware.Middleware, client *kubemq.Client) error {
	_, err := target.Do(ctx, msg)
	if err == nil {
		s.requeueCache.remove(msg.MessageID)
		return nil
	}
	if s.requeueCache.isRequeue(msg.MessageID) {
		_, requeueErr := client.SetQueueMessage(msg).Send(ctx)
		if requeueErr != nil {
			s.requeueCache.remove(msg.MessageID)
			s.log.Errorf("message id %s wasn't requeue due to an error , %s", msg.MessageID, requeueErr.Error())
			return nil
		}
		s.log.Infof("message id %s, requeued back to channel", msg.MessageID)
		return nil
	} else {
		return nil
	}
}

func (s *Source) Stop() error {
	for _, client := range s.clients {
		_ = client.Close()
	}
	return nil
}
