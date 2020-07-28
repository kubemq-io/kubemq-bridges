package queue

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/middleware"

	"github.com/kubemq-io/kubemq-go"
	"github.com/nats-io/nuid"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type mockTarget struct {
	setResponse interface{}
	setError    error
	delay       time.Duration
}

func (m *mockTarget) Do(ctx context.Context, request interface{}) (interface{}, error) {
	time.Sleep(m.delay)
	return m.setResponse, m.setError
}
func setupClient(ctx context.Context, target middleware.Middleware) (*Client, error) {
	c := New()

	err := c.Init(ctx, config.Spec{
		Name: "kubemq-queue",
		Kind: "",
		Properties: map[string]string{
			"address":      "localhost:50000",
			"client_id":    "some-client-id",
			"auth_token":   "",
			"channel":      "queue",
			"batch_size":   "1",
			"wait_timeout": "60",
		},
	})
	if err != nil {
		return nil, err
	}
	err = c.Start(ctx, target)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	return c, nil
}
func sendQueueMessage(t *testing.T, ctx context.Context, req *kubemq.QueueMessage, sendChannel string) error {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))

	if err != nil {
		return err
	}

	_, err = client.SetQueueMessage(req).SetChannel(sendChannel).Send(ctx)

	return err
}

func TestClient_processQueue(t *testing.T) {
	tests := []struct {
		name        string
		target      middleware.Middleware
		respChannel string
		req         *kubemq.QueueMessage
		sendCh      string
		wantErr     bool
	}{
		{
			name: "request",
			target: &mockTarget{
				setResponse: nil,
				setError:    nil,
				delay:       0,
			},
			req:     kubemq.NewQueueMessage().SetBody([]byte("some-data")),
			sendCh:  "queue",
			wantErr: false,
		},
		{
			name: "request with target error",
			target: &mockTarget{
				setResponse: nil,
				setError:    fmt.Errorf("some-error"),
				delay:       0,
			},
			req:     kubemq.NewQueueMessage().SetBody([]byte("some-data")),
			wantErr: false,
			sendCh:  "queue",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c, err := setupClient(ctx, tt.target)
			require.NoError(t, err)
			defer func() {
				_ = c.Stop()
			}()
			err = sendQueueMessage(t, ctx, tt.req, tt.sendCh)
			if tt.wantErr {
				require.Error(t, err)

			} else {
				require.NoError(t, err)
			}
			time.Sleep(time.Second)
		})
	}
}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name    string
		cfg     config.Spec
		wantErr bool
	}{
		{
			name: "init",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":      "localhost:50000",
					"client_id":    "",
					"auth_token":   "some-auth token",
					"channel":      "some-channel",
					"batch_size":   "1",
					"wait_timeout": "60",
				},
			},
			wantErr: false,
		},
		{
			name: "init - error",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address": "localhost",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad channel",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":      "localhost:50000",
					"client_id":    "",
					"auth_token":   "some-auth token",
					"channel":      "",
					"batch_size":   "1",
					"wait_timeout": "60",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad batch size",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":      "localhost:50000",
					"client_id":    "",
					"auth_token":   "some-auth token",
					"channel":      "some-channel",
					"batch_size":   "-1",
					"wait_timeout": "60",
				},
			},
			wantErr: true,
		}, {
			name: "init - bad wait timeout",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":      "localhost:50000",
					"client_id":    "",
					"auth_token":   "some-auth token",
					"channel":      "some-channel",
					"batch_size":   "1",
					"wait_timeout": "-1",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.cfg); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
			require.EqualValues(t, tt.cfg.Name, c.Name())
		})
	}
}
