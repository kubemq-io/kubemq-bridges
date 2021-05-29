package events

import (
	"context"
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/middleware"
	"github.com/kubemq-hub/kubemq-bridges/pkg/uuid"
	"github.com/kubemq-io/kubemq-go"
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
func setupSource(ctx context.Context, targets []middleware.Middleware) (*Source, error) {
	s := New()
	err := s.Init(ctx, config.Metadata{
		"address":                    "localhost:50000",
		"client_id":                  "",
		"auth_token":                 "",
		"channel":                    "events",
		"group":                      "some-group",
		"auto_reconnect":             "true",
		"reconnect_interval_seconds": "1",
		"max_reconnects":             "0",
		"sources":                    "2",
	}, config.Metadata{}, nil)
	if err != nil {
		return nil, err
	}

	err = s.Start(ctx, targets)
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	return s, nil
}
func sendEvent(t *testing.T, ctx context.Context, req *kubemq.Event, sendChannel string) error {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))

	if err != nil {
		return err
	}
	return client.SetEvent(req).SetChannel(sendChannel).Send(ctx)

}
func TestClient_processEvent(t *testing.T) {
	tests := []struct {
		name   string
		target middleware.Middleware
		req    *kubemq.Event
		sendCh string
		respCh string

		wantErr bool
	}{
		{
			name: "request",
			target: &mockTarget{
				setResponse: nil,
				setError:    nil,
				delay:       0,
			},
			req:     kubemq.NewEvent().SetBody([]byte("some-data")),
			wantErr: false,
			sendCh:  "events",
		},
		{
			name: "request with target error",
			target: &mockTarget{
				setResponse: nil,
				setError:    fmt.Errorf("some-error"),
				delay:       0,
			},
			req:     kubemq.NewEvent().SetBody([]byte("some-data")),
			wantErr: false,
			sendCh:  "events",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c, err := setupSource(ctx, []middleware.Middleware{tt.target})
			require.NoError(t, err)
			defer func() {
				_ = c.Stop()
			}()

			err = sendEvent(t, ctx, tt.req, tt.sendCh)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestClient_Init(t *testing.T) {

	tests := []struct {
		name       string
		connection config.Metadata
		wantErr    bool
	}{
		{
			name: "init",
			connection: config.Metadata{
				"address":                    "localhost:50000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "1",
			},
			wantErr: false,
		},
		{
			name: "init - error",
			connection: config.Metadata{
				"address": "localhost",
			},
			wantErr: true,
		},
		{
			name: "init - bad connection",
			connection: config.Metadata{
				"address":                    "localhost:40000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "1",
			},
			wantErr: true,
		},
		{
			name: "init - bad channel",
			connection: config.Metadata{
				"address":                    "localhost:50000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "1",
			},
			wantErr: true,
		},
		{
			name: "init - bad reconnect interval",
			connection: config.Metadata{
				"address":                    "localhost:50000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "-1",
				"max_reconnects":             "0",
				"sources":                    "1",
			},
			wantErr: true,
		},
		{
			name: "init - bad sources",
			connection: config.Metadata{
				"address":                    "localhost:50000",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "-1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.connection, config.Metadata{}, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
