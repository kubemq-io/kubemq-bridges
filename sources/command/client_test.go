package command

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
		Name: "kubemq-rpc",
		Kind: "",
		Properties: map[string]string{
			"address":                    "localhost:50000",
			"client_id":                  "responseid",
			"auth_token":                 "some-auth token",
			"channel":                    "commands",
			"group":                      "group",
			"concurrency":                "1",
			"auto_reconnect":             "true",
			"reconnect_interval_seconds": "1",
			"max_reconnects":             "0",
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
func sendCommand(t *testing.T, ctx context.Context, req *kubemq.Command, sendChannel string, timeout time.Duration) (*kubemq.CommandResponse, error) {
	client, err := kubemq.NewClient(ctx,
		kubemq.WithAddress("localhost", 50000),
		kubemq.WithClientId(nuid.Next()),
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
			name: "request - command target - executed",
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
			name: "request - command target - not executed",
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
			name: "request - query target - executed",
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
			name: "request - query target - not executed",
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
			name: "request - command target - error",
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
			name: "request - other target type - executed",
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
			c, err := setupClient(ctx, tt.target)
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
					"address":                    "localhost:50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "1",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
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
					"address":                    "localhost:50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "",
					"group":                      "",
					"concurrency":                "1",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad concurrency",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "0",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "1",
					"max_reconnects":             "0",
				},
			},
			wantErr: true,
		},
		{
			name: "init - bad reconnect interval",
			cfg: config.Spec{
				Name: "kubemq-rpc",
				Kind: "",
				Properties: map[string]string{
					"address":                    "localhost:50000",
					"client_id":                  "",
					"auth_token":                 "some-auth token",
					"channel":                    "some-channel",
					"group":                      "",
					"concurrency":                "1",
					"auto_reconnect":             "true",
					"reconnect_interval_seconds": "-1",
					"max_reconnects":             "0",
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
