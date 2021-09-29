package command

import (
	"context"
	"fmt"
	"github.com/kubemq-io/kubemq-bridges/config"
	"github.com/kubemq-io/kubemq-bridges/middleware"
	"github.com/kubemq-io/kubemq-bridges/pkg/uuid"
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
		"client_id":                  "responseid",
		"auth_token":                 "some-auth token",
		"channel":                    "commands",
		"group":                      "group",
		"auto_reconnect":             "true",
		"reconnect_interval_seconds": "1",
		"max_reconnects":             "0",
		"sources":                    "1",
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
func sendCommand(t *testing.T, ctx context.Context, req *kubemq.Command, sendChannel string, timeout time.Duration) (*kubemq.CommandResponse, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(uuid.New().String()),
		kubemq.WithTransportType(kubemq.TransportTypeGRPC))
	require.NoError(t, err)
	return client.SetCommand(req).SetChannel(sendChannel).SetTimeout(timeout).Send(ctx)

}
func TestClient_processCommand(t *testing.T) {
	tests := []struct {
		name     string
		target   middleware.Middleware
		req      *kubemq.Command
		wantResp *kubemq.CommandResponse
		timeout  time.Duration
		sendCh   string
		wantErr  bool
	}{
		{
			name: "request - command targets - executed",
			target: &mockTarget{
				setResponse: &kubemq.CommandResponse{
					CommandId:        "id",
					ResponseClientId: "responseid",
					Executed:         true,
					ExecutedAt:       time.Unix(1000, 0),
					Error:            "",
					Tags:             nil,
				},
				setError: nil,
				delay:    0,
			},
			req: kubemq.NewCommand().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.CommandResponse{
				CommandId:        "id",
				ResponseClientId: "responseid",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Error:            "",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "commands",
			wantErr: false,
		},
		{
			name: "request - command targets - not executed",
			target: &mockTarget{
				setResponse: &kubemq.CommandResponse{
					CommandId:        "id",
					ResponseClientId: "responseid",
					Executed:         false,
					Error:            "some-error",
					Tags:             nil,
				},
				setError: nil,
				delay:    0,
			},
			req: kubemq.NewCommand().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.CommandResponse{
				CommandId:        "id",
				ResponseClientId: "responseid",
				Executed:         false,
				ExecutedAt:       time.Time{}.In(time.Local),
				Error:            "some-error",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "commands",
			wantErr: false,
		},
		{
			name: "request - query targets - executed",
			target: &mockTarget{
				setResponse: &kubemq.QueryResponse{
					QueryId:          "id",
					Executed:         true,
					ExecutedAt:       time.Unix(1000, 0),
					Metadata:         "some-metadata",
					ResponseClientId: "responseid",
					Body:             nil,
					CacheHit:         false,
					Error:            "",
					Tags:             nil,
				},
				setError: nil,
				delay:    0,
			},
			req: kubemq.NewCommand().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.CommandResponse{
				CommandId:        "id",
				ResponseClientId: "responseid",
				Executed:         true,
				ExecutedAt:       time.Unix(1000, 0),
				Error:            "",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "commands",
			wantErr: false,
		},
		{
			name: "request - query targets - not executed",
			target: &mockTarget{
				setResponse: &kubemq.QueryResponse{
					QueryId:          "id",
					Executed:         false,
					ExecutedAt:       time.Time{},
					Metadata:         "some-metadata",
					ResponseClientId: "responseid",
					Body:             nil,
					CacheHit:         false,
					Error:            "some-error",
					Tags:             nil,
				},
			},
			req: kubemq.NewCommand().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.CommandResponse{
				CommandId:        "id",
				ResponseClientId: "responseid",
				Executed:         false,
				ExecutedAt:       time.Time{}.In(time.Local),
				Error:            "some-error",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "commands",
			wantErr: false,
		},
		{
			name: "request - command targets - error",
			target: &mockTarget{
				setResponse: nil,
				setError:    fmt.Errorf("some-error"),
				delay:       0,
			},
			req: kubemq.NewCommand().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.CommandResponse{
				CommandId:        "id",
				ResponseClientId: "responseid",
				Executed:         false,
				ExecutedAt:       time.Time{}.In(time.Local),
				Error:            "some-error",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "commands",
			wantErr: false,
		},
		{
			name: "request - other targets type - executed",
			target: &mockTarget{
				setResponse: nil,
				setError:    nil,
				delay:       0,
			},
			req: kubemq.NewCommand().SetId("id").SetBody([]byte("some-data")),
			wantResp: &kubemq.CommandResponse{
				CommandId:        "id",
				ResponseClientId: "responseid",
				Executed:         true,
				ExecutedAt:       time.Time{}.In(time.Local),
				Error:            "",
				Tags:             nil,
			},
			timeout: 5 * time.Second,
			sendCh:  "commands",
			wantErr: false,
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
			gotResp, err := sendCommand(t, ctx, tt.req, tt.sendCh, tt.timeout)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.EqualValues(t, tt.wantResp, gotResp)
		})
	}
}

//
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
				"sources":                    "2",
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
				"sources":                    "2",
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
				"sources":                    "2",
			},
			wantErr: true,
		},
		{
			name: "init - bad connection",
			connection: config.Metadata{
				"address":                    "localhost",
				"client_id":                  "",
				"auth_token":                 "some-auth token",
				"channel":                    "some-channel",
				"group":                      "",
				"auto_reconnect":             "true",
				"reconnect_interval_seconds": "1",
				"max_reconnects":             "0",
				"sources":                    "2",
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
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			c := New()
			if err := c.Init(ctx, tt.connection, config.Metadata{}, nil); (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
